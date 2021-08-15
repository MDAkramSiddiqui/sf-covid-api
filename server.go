package main

import (
	"os"

	"github.com/MDAkramSiddiqui/sf-covid-api/app/constants"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/controllers"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/log"
	"github.com/MDAkramSiddiqui/sf-covid-api/crons"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

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

func main() {
	e := echo.New()

	// middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())

	// routes
	e.GET("/ping", controllers.HealthController)
	e.GET("/v1/covid-data/state", controllers.StateController)

	e.Logger.Fatal(e.Start(":5000"))
	crons.StateDataCron.Start()
}
