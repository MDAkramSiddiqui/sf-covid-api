package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetRequest(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ErrorMessageURLNotFound", "", err.Error())
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("ErrorMessageResponseNotParsed", "", url)
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("ErrorMessageCouldNotGetURL", "", url)
	}
	return body, nil
}
