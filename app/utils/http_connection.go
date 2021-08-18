package utils

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/MDAkramSiddiqui/sf-covid-api/app/log"
)

// Wrapper for making get request and returns correspoing results and errors
func GetRequest(url string) ([]byte, error) {
	log.Instance.Debug("GetRequest is hit")

	response, err := http.Get(url)
	if err != nil {
		log.Instance.Err("URL not found, err: %v", err.Error())
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Instance.Err("Response not parsed, err: %v", err.Error())
		return nil, err
	}

	if response.StatusCode != 200 {
		log.Instance.Err("Requested data fetch call failed")
		return nil, errors.New("requested data fetch call failed")
	}

	return body, nil
}
