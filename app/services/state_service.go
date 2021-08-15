package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
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
	var data bson.M

	stateName = strings.Trim(strings.TrimSpace(stateName), "\"")
	mongoDriverInstance, _ := drivers.GetMongoDriver()
	coll := mongoDriverInstance.Database(os.Getenv(constants.MongoDBName)).Collection("covid-state")

	_ = coll.FindOne(
		context.TODO(),
		bson.M{"name": stateName},
	).Decode(&data)

	return data
}

func GetAllStateCovidData() []primitive.M {
	var data []bson.M
	mongoDriverInstance, _ := drivers.GetMongoDriver()
	coll := mongoDriverInstance.Database(os.Getenv(constants.MongoDBName)).Collection("covid-state")
	cursor, _ := coll.Find(context.TODO(), bson.M{})
	_ = cursor.All(context.TODO(), &data)
	return data
}

func GetAllStateCovidDataGovtApi() []byte {
	body, _ := utils.GetRequest(constants.CovidDataApi)
	return body
}

func GetStateNameUsingLatAndLong(latLang []string) string {
	var stateName string
	var stateData StateItems

	hereGeoCordinateApiMapper := map[string]string{
		"API_KEY": os.Getenv(constants.HereGeoAPIKey),
		"LAT":     latLang[0],
		"LONG":    latLang[1],
	}

	buf := bytes.Buffer{}
	t := template.Must(template.New("").Parse(constants.HereGeoCordinateApi))
	t.Execute(&buf, hereGeoCordinateApiMapper)
	url := buf.String()

	response, _ := utils.GetRequest(url)
	json.Unmarshal(response, &stateData)

	if len(stateData.Items) > 0 {
		stateName = stateData.Items[0].Address.StateName
		stateName = strings.ReplaceAll(stateName, "&", "and")
		fmt.Println(stateName)
		return stateName
	}

	return stateName
}
