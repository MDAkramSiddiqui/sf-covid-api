package main

import (
	"github.com/MDAkramSiddiqui/sf-covid-api/app/controllers"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())

	// routes
	e.GET("/", controllers.HealthController)
	e.GET("/v1/state", controllers.StateController)

	logger.Init()
	logger.SetLogLevel(0)
	logger.DEBUG("Server Started")
	e.Logger.Fatal(e.Start(":5000"))
}
