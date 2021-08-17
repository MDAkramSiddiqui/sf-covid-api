package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/MDAkramSiddiqui/sf-covid-api/app/constants"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/controllers"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/drivers"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/log"
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

	if os.Getenv(constants.Env) == constants.Production {
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
			crons.StateDataCron.Stop()
		}
	}()

	log.Instance.Info("Starting server at port %v", port)
	crons.StateDataCron.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Instance.Fatal("Shutting down server failed, err: %v", err.Error())
	} else {
		log.Instance.Info("Server shut down successfully")
	}
}
