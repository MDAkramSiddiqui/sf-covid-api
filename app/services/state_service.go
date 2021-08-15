package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/MDAkramSiddiqui/sf-covid-api/app/constants"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/drivers"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/utils"
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
	mongoDriverInstance, _ := drivers.GetMongoDriver()
	coll := mongoDriverInstance.Database(os.Getenv(constants.MongoDBName)).Collection("covid-state")

	stateName = strings.Trim(strings.TrimSpace(stateName), "\"")

	_ = coll.FindOne(
		context.TODO(),
		bson.M{"name": stateName},
	).Decode(&data)

	return data
}

func StateService2() []primitive.M {
	var data []bson.M
	mongoDriverInstance, _ := drivers.GetMongoDriver()
	coll := mongoDriverInstance.Database(os.Getenv(constants.MongoDBName)).Collection("covid-state")
	cursor, _ := coll.Find(context.TODO(), bson.M{})
	_ = cursor.All(context.TODO(), &data)
	return data
}

func FetchCovidStateWiseData() []byte {
	body, _ := utils.GetRequest(constants.CovidDataApi)
	return body
}

func FetchStateName(url string) string {
	body, _ := utils.GetRequest(url)
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
