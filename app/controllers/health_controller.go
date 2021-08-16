package controllers

import (
	"net/http"

	"github.com/MDAkramSiddiqui/sf-covid-api/app/log"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/response_model"
	"github.com/labstack/echo/v4"
)

// HealthCheck godoc
// @Summary Show the status of server.
// @Description Get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /ping [get]
func HealthController(c echo.Context) error {
	log.Instance.Debug("HealthController is hit")

	return c.JSON(http.StatusOK, response_model.DefaultResponse(http.StatusOK, "pong"))
}
