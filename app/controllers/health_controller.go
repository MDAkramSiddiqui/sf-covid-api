package controllers

import (
	"net/http"

	"github.com/MDAkramSiddiqui/sf-covid-api/app/log"
	"github.com/labstack/echo/v4"
)

func HealthController(c echo.Context) error {
	log.Instance.Debug("HealthController is hit")

	return c.String(http.StatusOK, "pong")
}
