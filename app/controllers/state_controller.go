package controllers

import (
	"net/http"
	"net/url"

	"github.com/MDAkramSiddiqui/sf-covid-api/app/log"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/response_model"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/services"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type customStateChannel struct {
	Result primitive.M
	Err    *utils.CustomErr
}

// Statewise Covid Data Doc godoc
// @Summary Serves statewise covid data.
// @Description Get statewise covid data either via using name or latitude and longitude.
// @Tags root
// @Accept */*
// @Produce json
// @Param name query string false "State name for which covid data is required"
// @Param latlng query string false "Latitude and longitude of user, eg. latlng=23.223,23.222"
// @Success 200 {object} map[string]interface{}
// @Router /covid-data/state [get]
func StateController(c echo.Context) error {
	log.Instance.Debug("StateController is hit")

	var stateName string
	var latLangQuery string
	var responseData []bson.M

	// checking provided state name
	stateName = c.QueryParam("name")
	latLangQuery, _ = url.QueryUnescape(c.QueryParam("latlng"))

	log.Instance.Debug("Raw state name provided is %v", stateName)
	log.Instance.Debug("Raw coordinates provided are %v", latLangQuery)

	dataByStateName := &customStateChannel{}
	dataByCoordinate := &customStateChannel{}

	dataByStateNamechan := make(chan bool)
	dataByCoordinateChan := make(chan bool)

	// get data via state name
	go func(dataByStateName *customStateChannel, dataByStateNamechan chan bool) {
		dataByStateName.Result, dataByStateName.Err = services.GetCovidDataByName(stateName)
		dataByStateNamechan <- true
	}(dataByStateName, dataByStateNamechan)

	// get data via coordinates
	go func(dataByCoordinate *customStateChannel, dataByCoordinateChan chan bool) {
		dataByCoordinate.Result, dataByCoordinate.Err = services.GetCovidDataByCoordinates(latLangQuery)
		dataByCoordinateChan <- true
	}(dataByCoordinate, dataByCoordinateChan)

	<-dataByStateNamechan
	<-dataByCoordinateChan

	if dataByStateName.Err.Err != nil {
		return c.JSON(response_model.DefaultResponse(dataByStateName.Err.StatusCode, dataByStateName.Err.Message(), true))
	} else if dataByStateName.Result != nil {
		log.Instance.Info("Covid data by state name found, appending to response")
		responseData = append(responseData, dataByStateName.Result)
	}

	if dataByCoordinate.Err.Err != nil {
		return c.JSON(response_model.DefaultResponse(dataByCoordinate.Err.StatusCode, dataByCoordinate.Err.Message(), true))
	} else if dataByCoordinate.Result != nil {
		log.Instance.Info("Covid data by coordinates found, appending to response")
		responseData = append(responseData, dataByCoordinate.Result)
	}

	// fetch all states data if none query params are provided
	if len(responseData) == 0 {
		log.Instance.Info("Fetching all states data as coordinates and state name not provided")

		allStatesData, allStatesDataErr := services.GetAllStateCovidData()
		if allStatesDataErr.Err != nil {
			return c.JSON(response_model.DefaultResponse(allStatesDataErr.StatusCode, allStatesDataErr.Message(), true))
		} else {
			responseData = allStatesData
		}
	}

	return c.JSON(response_model.DefaultResponse(http.StatusOK, responseData, false))
}
