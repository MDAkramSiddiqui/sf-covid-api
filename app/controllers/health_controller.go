package controllers

import (
	"net/http"

	"github.com/MDAkramSiddiqui/sf-covid-api/app/log"
	"github.com/labstack/echo/v4"
)

// HealthCheck godoc
// @Summary Show the status of server.
// @Description Get the status of server.
// @Tags root
// @Accept */*
// @Produce plain
// @Success 200
// @Router /ping [get]
func HealthController(c echo.Context) error {
	log.Instance.Debug("HealthController is hit")

	return c.String(http.StatusOK, "pong")
}
