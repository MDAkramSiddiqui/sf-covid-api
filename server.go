package main

import (
	"fmt"
	"os"

	"github.com/MDAkramSiddiqui/sf-covid-api/app/constants"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/controllers"
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
	err := godotenv.Load()
	if err != nil {
		log.Instance.Fatal("Error while loading environment variables", err)
	}

	log.Init()
	if os.Getenv(constants.Env) == constants.PRODUCTION {
		log.Instance.SetLogLevel(3)
	} else {
		log.Instance.SetLogLevel(0)
	}
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
// @schemes http
func main() {
	e := echo.New()

	// middlewares
	e.Use(middleware.LoggerWithConfig(DefaultLoggerConfig))
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// routes
	e.GET("/api/v1/ping", controllers.HealthController)
	e.GET("/api/v1/covid-data/state", controllers.StateController)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	port := os.Getenv(constants.Port)
	if port == "" {
		port = "5000"
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", port)))
	crons.StateDataCron.Start()
}
