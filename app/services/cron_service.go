package services

import (
	"context"
	"encoding/json"

	"github.com/MDAkramSiddiqui/sf-covid-api/app/logger"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CovidState struct {
	Name             string `json:"state_name" bson:"state_name"`
	ActiveCases      string `json:"active" bson:"active"`
	PositiveCases    string `json:"positive" bson:"positive"`
	CuredCases       string `json:"cured" bson:"cured"`
	DeathCases       string `json:"death" bson:"death"`
	NewActiveCases   string `json:"new_active" bson:"new_active"`
	NewPositiveCases string `json:"new_positive" bson:"new_positive"`
	NewCuredCases    string `json:"new_cured" bson:"new_cured"`
	NewDeathCases    string `json:"new_death" bson:"new_death"`
	StateCode        string `json:"state_code" bson:"state_code"`
}

func init() {
	c := cron.New()
	c.AddFunc("*/1 * * * *", updateCovidData)
	c.Start()
}

func updateCovidData() {
	data := FetchCovidStateWiseData()
	logger.INFO("PRINTED EVERYTHONF")
	var covidStatesData []CovidState
	json.Unmarshal(data, &covidStatesData)

	client, err := GetMongoClient()
	if err != nil {
		return
	}

	for i := 0; i < len(covidStatesData); i++ {
		opts := options.FindOneAndReplace().SetUpsert(true)
		filter := bson.M{"name": covidStatesData[i].Name}
		replacement := bson.D{
			{Key: "name", Value: covidStatesData[i].Name},
			{Key: "positiveCases", Value: covidStatesData[i].PositiveCases},
			{Key: "activeCases", Value: covidStatesData[i].ActiveCases},
			{Key: "deathCases", Value: covidStatesData[i].DeathCases},
			{Key: "curedCases", Value: covidStatesData[i].CuredCases},
			{Key: "newPositiveCases", Value: covidStatesData[i].NewPositiveCases},
			{Key: "newActiveCases", Value: covidStatesData[i].NewActiveCases},
			{Key: "newDeathCases", Value: covidStatesData[i].NewDeathCases},
			{Key: "newCuredCases", Value: covidStatesData[i].NewCuredCases},
		}

		coll := client.Database("covid").Collection("covid-state")
		_ = coll.FindOneAndReplace(
			context.TODO(),
			filter,
			replacement,
			opts,
		)
	}
}
