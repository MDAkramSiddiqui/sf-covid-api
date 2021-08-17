package crons

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/MDAkramSiddiqui/sf-covid-api/app/constants"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/drivers"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/log"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/schema"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/services"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DataCron struct {
	name string
	job  *cron.Cron
}

var StateDataCron *DataCron

func init() {
	StateDataCron = &DataCron{"StateDataCron", cron.New()}
	StateDataCron.job.AddFunc("*/30 * * * *", updateCovidData)
}

func (c *DataCron) Start() {
	c.job.Start()
	log.Instance.Info("%v job started successfully", c.name)
}

func (c *DataCron) Stop() {
	c.job.Stop()
	log.Instance.Info("%v job stopped successfully", c.name)
}

func updateCovidData() {
	log.Instance.Debug("updateCovidData is hit")

	var covidStatesData []schema.TCovidState

	data := services.GetAllStateCovidDataGovtApi()
	json.Unmarshal(data, &covidStatesData)

	mongoDriverInstance, err := drivers.GetMongoDriver()
	if err != nil {
		return
	}

	for i := 0; i < len(covidStatesData); i++ {
		stateName := covidStatesData[i].Name
		if stateName == "" {
			stateName = "India"
		}

		opts := options.FindOneAndReplace().SetUpsert(true)
		filter := bson.M{"name": stateName}
		replacement := bson.D{
			{Key: "name", Value: stateName},
			{Key: "positiveCases", Value: covidStatesData[i].PositiveCases},
			{Key: "activeCases", Value: covidStatesData[i].ActiveCases},
			{Key: "deathCases", Value: covidStatesData[i].DeathCases},
			{Key: "curedCases", Value: covidStatesData[i].CuredCases},
			{Key: "latestPositiveCases", Value: covidStatesData[i].LatestPositiveCases},
			{Key: "latestActiveCases", Value: covidStatesData[i].LatestActiveCases},
			{Key: "latestDeathCases", Value: covidStatesData[i].LatestDeathCases},
			{Key: "latestCuredCases", Value: covidStatesData[i].LatestCuredCases},
			{Key: "updatedAt", Value: time.Now()},
		}

		coll := mongoDriverInstance.Database(os.Getenv(constants.MongoDBName)).Collection("covid-state")
		_ = coll.FindOneAndReplace(
			context.TODO(),
			filter,
			replacement,
			opts,
		)
	}
}
