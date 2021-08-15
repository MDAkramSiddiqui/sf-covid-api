package controllers

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/MDAkramSiddiqui/sf-covid-api/app/constants"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/logger"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/services"
	"github.com/labstack/echo/v4"
)

func StateController(c echo.Context) error {
	var stateName string
	logger.INFO("State Controller hit")
	stateName = c.QueryParam("name")
	latLang := strings.Split(c.QueryParam("latlng"), ",")
	if c.QueryParam("latlng") != "" && len(latLang) != 2 {
		fmt.Println("Invalid params")
	}

	if len(latLang) == 2 {
		data := map[string]string{
			"API_KEY": os.Getenv(constants.HereGeoAPIKey),
			"LAT":     latLang[0],
			"LONG":    latLang[1],
		}
		buf := bytes.Buffer{}
		t := template.Must(template.New("").Parse(constants.HereGeoCordinateApi))
		t.Execute(&buf, data)
		url := buf.String()
		stateName = services.FetchStateName(url)
	}

	if stateName == "" {
		val := services.StateService2()
		return c.JSON(http.StatusOK, val)
	}
	resp := services.StateService(stateName)
	return c.JSON(http.StatusOK, resp)
}
