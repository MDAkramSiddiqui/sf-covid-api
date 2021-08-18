package utils

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/MDAkramSiddiqui/sf-covid-api/app/log"
)

// Wrapper for making get request and returns correspoing results and errors
func GetRequest(url string) ([]byte, *CustomErr) {
	log.Instance.Debug("GetRequest is hit")

	response, err := http.Get(url)
	if err != nil {
		log.Instance.Err("URL not found, err: %v", err.Error())
		return nil, &CustomErr{Err: err, StatusCode: http.StatusInternalServerError}
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Instance.Err("Response not parsed, err: %v", err.Error())
		return nil, &CustomErr{Err: err, StatusCode: http.StatusInternalServerError}
	}

	if response.StatusCode != 200 {
		log.Instance.Err("Requested data fetch call failed")
		return nil, &CustomErr{Err: errors.New("requested data fetch call failed"), StatusCode: response.StatusCode}
	}

	return body, &CustomErr{}
}
