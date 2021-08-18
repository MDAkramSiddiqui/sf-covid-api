package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/MDAkramSiddiqui/sf-covid-api/app/constants"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/controllers"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/drivers"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/log"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/response_model"
	"github.com/MDAkramSiddiqui/sf-covid-api/crons"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/MDAkramSiddiqui/sf-covid-api/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

var DefaultLoggerConfig = middleware.LoggerConfig{
	Skipper:          middleware.DefaultSkipper,
	Format:           `[MIDDL]: ${time_custom} Req.Id: ${id}, Method: ${method}, URI: ${uri}, IP: ${remote_ip}, Latency: ${latency}, UserAgent: ${user_agent}` + "\n",
	CustomTimeFormat: "2006/01/02 15:04:05",
}

func init() {
	log.Init()

	err := godotenv.Load()
	if err != nil {
		log.Instance.Err("Error while loading environment variables, err: %v", err.Error())
	}

	if os.Getenv(constants.LogLevel) != "" {
		logLevel, err := strconv.Atoi(os.Getenv(constants.LogLevel))
		if err != nil || logLevel < constants.DebugLevel || logLevel > constants.FatalLevel {
			log.Instance.Fatal("Invalid log level provided")
		}
		log.Instance.SetLogLevel(logLevel)

	} else if os.Getenv(constants.Env) == constants.Production {
		log.Instance.SetLogLevel(constants.ErrLevel)
	} else {
		log.Instance.SetLogLevel(constants.DebugLevel)
	}

	mongoConnectionChan := make(chan bool)
	redisConnectionChan := make(chan bool)
	go func(mongoConnectionChan chan bool) {
		drivers.GetMongoDriver()
		mongoConnectionChan <- true
	}(mongoConnectionChan)

	go func(redisConnectionChan chan bool) {
		drivers.GetRedisDriver()
		redisConnectionChan <- true
	}(redisConnectionChan)

	<-mongoConnectionChan
	<-redisConnectionChan

	crons.Init()
}

// custom http error handler
func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	log.Instance.Err("Some internal error happened, err: %v", err.Error())
	c.JSON(response_model.DefaultResponse(code, nil, false))
}

// custom 404 error handler
func custom404ErrorHandler(c echo.Context) error {
	log.Instance.Err("Requested url not found on server, url: %v", c.Request().URL)
	return c.JSON(response_model.DefaultResponse(http.StatusNotFound, "requested url not found", true))
}

// @title SF-Covid-State Api
// @version 1.0
// @description This is a simple server for requesting covid data of state
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api/v1/
// @schemes https http
func main() {
	e := echo.New()

	// Set custom http error handler
	e.HTTPErrorHandler = customHTTPErrorHandler

	// Set custom 404 error handler
	echo.NotFoundHandler = custom404ErrorHandler

	// middlewares
	e.Use(middleware.LoggerWithConfig(DefaultLoggerConfig))
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// routes
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "/swagger/index.html")
	})
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/api/v1/ping", controllers.HealthController)
	e.GET("/api/v1/covid-data/state", controllers.StateController)

	port := os.Getenv(constants.Port)
	if port == "" {
		port = "5000"
		log.Instance.Info("Server port not provided, using default port")
	}

	// Start server
	go func() {
		if err := e.Start(fmt.Sprintf(":%v", port)); err != nil {
			if err != http.ErrServerClosed {
				log.Instance.Fatal("Server start failed, shutting down server, err: %v", err.Error())
			}
		}
	}()

	log.Instance.Info("Starting server at port %v", port)
	crons.StateDataCron.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// Stop Cron job
	crons.StateDataCron.Stop()

	// close connections with DB and redis
	mongoConnectionChan := make(chan bool)
	redisConnectionChan := make(chan bool)

	mongoDriverInstance, mongoDriverInstanceErr := drivers.GetMongoDriver()
	if mongoDriverInstance != nil && mongoDriverInstanceErr == nil {
		go func(mongoConnectionChan chan bool) {
			err := mongoDriverInstance.Disconnect(context.TODO())
			if err != nil {
				log.Instance.Err("Error while closing connection with DB, err: %v", err.Error())
			} else {
				log.Instance.Info("Connection with DB closed successfully")
			}
			mongoConnectionChan <- true
		}(mongoConnectionChan)
	}

	redisDriverInstance, redisDriverInstanceErr := drivers.GetRedisDriver()
	if redisDriverInstance != nil && redisDriverInstanceErr == nil {
		go func(redisConnectionChan chan bool) {
			err := redisDriverInstance.Close()
			if err != nil {
				log.Instance.Err("Error while closing connection with redis, err: %v", err.Error())
			} else {
				log.Instance.Info("Connection with redis closed successfully")
			}
			redisConnectionChan <- true
		}(redisConnectionChan)
	}

	<-mongoConnectionChan
	<-redisConnectionChan

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Instance.Fatal("Shutting down server failed, err: %v", err.Error())
	} else {
		log.Instance.Info("Server shut down successfully")
	}
}
