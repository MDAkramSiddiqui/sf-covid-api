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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DataCron struct {
	name string
	job  *cron.Cron
}

var StateDataCron *DataCron

func Init() {
	stateCronPeriod := constants.DefaultStateDataCronPeriod
	if os.Getenv(constants.StateDataCronPeriod) != "" {
		log.Instance.Info("State cronjob period changed from default %v to %v", stateCronPeriod, os.Getenv(constants.StateDataCronPeriod))
		stateCronPeriod = os.Getenv(constants.StateDataCronPeriod)
	}

	StateDataCron = &DataCron{"StateDataCron", cron.New()}
	StateDataCron.job.AddFunc(stateCronPeriod, updateCovidData)
}

func (c *DataCron) Start() {
	c.job.Start()
	log.Instance.Info("%v job started successfully", c.name)
}

func (c *DataCron) Stop() {
	c.job.Stop()
	log.Instance.Info("%v job stopped successfully", c.name)
}

// cronjob for periodically fetching and updating covid data of different states
// from 3rd party API
func updateCovidData() {
	log.Instance.Debug("updateCovidData is hit")

	var covidStatesData []schema.TCovidState

	data, err := services.GetAllStateCovidDataGovtApi()
	if err != nil {
		log.Instance.Err("Error while fetching all states data from 3rd party API, err: %v", err.Error())
		return
	}

	json.Unmarshal(data, &covidStatesData)

	mongoDriverInstance, err := drivers.GetMongoDriver()
	if err != nil {
		return
	}

	statesChan := make([]chan bool, len(covidStatesData))

	for i := 0; i < len(covidStatesData); i++ {
		go updateStateData(&covidStatesData[i], mongoDriverInstance, statesChan[i])
	}

	for i := 0; i < len(statesChan); i++ {
		<-statesChan[i]
	}
}

func updateStateData(covidStateData *schema.TCovidState, mongoDriverInstance *mongo.Client, stateChan chan bool) {
	stateName := covidStateData.Name
	if stateName == "" {
		stateName = "India"
	}

	opts := options.FindOneAndReplace().SetUpsert(true)
	filter := bson.M{"name": stateName}
	replacement := bson.D{
		{Key: "name", Value: stateName},
		{Key: "positiveCases", Value: covidStateData.PositiveCases},
		{Key: "activeCases", Value: covidStateData.ActiveCases},
		{Key: "deathCases", Value: covidStateData.DeathCases},
		{Key: "curedCases", Value: covidStateData.CuredCases},
		{Key: "latestPositiveCases", Value: covidStateData.LatestPositiveCases},
		{Key: "latestActiveCases", Value: covidStateData.LatestActiveCases},
		{Key: "latestDeathCases", Value: covidStateData.LatestDeathCases},
		{Key: "latestCuredCases", Value: covidStateData.LatestCuredCases},
		{Key: "updatedAt", Value: time.Now()},
	}

	coll := mongoDriverInstance.Database(os.Getenv(constants.MongoDBName)).Collection("covid-state")
	_ = coll.FindOneAndReplace(
		context.TODO(),
		filter,
		replacement,
		opts,
	)

	log.Instance.Debug("Data for state %v updated successfully", covidStateData.Name)
	stateChan <- true
}
