package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/MDAkramSiddiqui/sf-covid-api/app/services"
	"github.com/labstack/echo/v4"
)

func StateController(c echo.Context) error {
	var stateName string
	var latLang []string

	stateName = c.QueryParam("name")
	latLang = strings.Split(c.QueryParam("latlng"), ",")

	if len(latLang) == 2 {
		stateName = services.GetStateNameUsingLatAndLong(latLang)
	} else {
		fmt.Println("latitude and longitude invalid or not provided")
	}

	if stateName == "" {
		fmt.Println("State name not provided")
		val := services.GetAllStateCovidData()
		return c.JSON(http.StatusOK, val)
	}
	resp := services.StateService(stateName)
	return c.JSON(http.StatusOK, resp)
}
