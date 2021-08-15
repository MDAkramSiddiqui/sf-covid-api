package utils

import (
	"io/ioutil"
	"net/http"

	"github.com/MDAkramSiddiqui/sf-covid-api/app/log"
)

func GetRequest(url string) ([]byte, error) {
	log.Instance.Debug("GetRequest is hit")

	response, err := http.Get(url)
	if err != nil {
		log.Instance.Err("URL not found", err)
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Instance.Err("Response not parsed", err)
		return nil, err
	}

	if response.StatusCode != 200 {
		log.Instance.Err("Requested data fetch call failed", err)
		return nil, err
	}

	return body, nil
}
