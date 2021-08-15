package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/MDAkramSiddiqui/sf-covid-api/app/constants"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StateAddress struct {
	StateCode string `json:"stateCode" bson:"stateCode"`
	StateName string `json:"state" bson:"state"`
}

type StateData struct {
	Address StateAddress `json:"address" bson:"address"`
}

type StateItems struct {
	Items []StateData `json:"items" bson:"items"`
}

func StateService(stateName string) primitive.M {

	if stateName == "" {
		return nil
	}

	var data bson.M
	client, _ := GetMongoClient()
	coll := client.Database("covid").Collection("covid-state")

	stateName = strings.Trim(strings.TrimSpace(stateName), "\"")

	_ = coll.FindOne(
		context.TODO(),
		bson.M{"name": stateName},
	).Decode(&data)

	return data
}

func StateService2() []primitive.M {
	var data []bson.M
	client, _ := GetMongoClient()
	coll := client.Database("covid").Collection("covid-state")
	cursor, _ := coll.Find(context.TODO(), bson.M{})
	_ = cursor.All(context.TODO(), &data)
	return data
}

func FetchCovidStateWiseData() []byte {
	resp, err := http.Get(constants.CovidDataApi)
	if err != nil {
		logger.FATAL(err.Error())
		return nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return body
}

func FetchStateName(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		logger.FATAL(err.Error())
		return ""
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var stateData StateItems
	json.Unmarshal(body, &stateData)
	if len(stateData.Items) > 0 {
		stateName := stateData.Items[0].Address.StateName
		stateName = strings.ReplaceAll(stateName, "&", "and")
		fmt.Println(stateName)
		return stateName
	}
	return ""
}
