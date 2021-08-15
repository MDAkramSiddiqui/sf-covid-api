package controllers

import (
	"net/http"
	"strings"

	"github.com/MDAkramSiddiqui/sf-covid-api/app/log"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/services"
	"github.com/labstack/echo/v4"
)

func StateController(c echo.Context) error {
	log.Instance.Debug("StateController is hit")

	var stateName string
	var latLang []string

	stateName = c.QueryParam("name")
	latLang = strings.Split(c.QueryParam("latlng"), ",")

	if len(latLang) == 2 {
		stateName = services.GetStateNameUsingLatAndLong(latLang)
	} else {
		log.Instance.Warn("Latitude and longitude invalid or not provided")
	}

	if stateName == "" {
		log.Instance.Warn("State name not provided")
		val := services.GetAllStateCovidData()
		return c.JSON(http.StatusOK, val)
	}
	resp := services.GetStateCovidData(stateName)
	return c.JSON(http.StatusOK, resp)
}
