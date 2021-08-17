package controllers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/MDAkramSiddiqui/sf-covid-api/app/log"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/response_model"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/services"
	"github.com/labstack/echo/v4"
)

// Statewise Covid Data Doc godoc
// @Summary Serves statewise covid data.
// @Description Get statewise covid data either via using name or latitude and longitude.
// @Tags root
// @Accept */*
// @Produce json
// @Param name query string false "State name for which covid data is required"
// @Param latlng query string false "Latitude and longitude of user"
// @Success 200 {object} map[string]interface{}
// @Router /covid-data/state [get]
func StateController(c echo.Context) error {
	log.Instance.Debug("StateController is hit")

	var stateName string
	var latLang []string

	stateName = c.QueryParam("name")
	if stateName == "" {
		log.Instance.Info("State name is not provided")
	}

	latLangQuery, _ := url.QueryUnescape(c.QueryParam("latlng"))
	latLang = strings.Split(latLangQuery, ",")

	if len(latLang) == 2 {
		log.Instance.Info("Latitude and longitude provided are %v, %v", latLang[0], latLang[1])
		stateName = services.GetStateNameUsingLatAndLong(latLang)
	} else {
		log.Instance.Info("Latitude and longitude are not provided or invalid")
	}

	if stateName == "" {
		log.Instance.Info("Cannot determine requested state, therefore fetching all states covid data")
		val := services.GetAllStateCovidData()
		return c.JSON(http.StatusOK, response_model.DefaultResponse(http.StatusOK, val))
	}

	resp := services.GetStateCovidData(stateName)
	return c.JSON(http.StatusOK, response_model.DefaultResponse(http.StatusOK, resp))
}
