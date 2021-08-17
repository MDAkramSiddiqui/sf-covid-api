package services

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"os"
	"strings"

	"github.com/MDAkramSiddiqui/sf-covid-api/app/constants"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/drivers"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/log"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/schema"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Fetches a single state covid data from DB using provided state name or coordinates
func GetStateCovidData(stateName string) ([]primitive.M, error) {
	log.Instance.Debug("GetStateCovidData is hit")

	var data = make([]bson.M, 1)
	isFoundInRedis := true

	stateName = strings.Title(strings.ToLower(strings.TrimSpace(strings.Trim(stateName, "\""))))

	redisDriverInstance, redisDriverInstanceErr := drivers.GetRedisDriver()
	if redisDriverInstanceErr != nil {
		log.Instance.Err("Redis instance is down using DB to fetch data, err: %v", redisDriverInstanceErr.Error())
	} else {
		result, _ := redisDriverInstance.Get(stateName).Bytes()

		if len(result) == 0 {
			isFoundInRedis = false
			log.Instance.Info("Data for state: %v not found in redis requesting from DB", stateName)
		} else {
			json.Unmarshal(result, &data[0])
			log.Instance.Info("Data for state: %v found in redis", stateName)
			return data, nil
		}
	}

	mongoDriverInstance, err := drivers.GetMongoDriver()
	if err != nil {
		log.Instance.Err("DB is down, err: %v", stateName, err.Error())
		return data, err
	}

	coll := mongoDriverInstance.Database(os.Getenv(constants.MongoDBName)).Collection("covid-state")

	err = coll.FindOne(
		context.TODO(),
		bson.M{"name": stateName},
		options.FindOne().SetProjection(bson.M{"_id": 0}),
	).Decode(&data[0])

	if err != nil {
		log.Instance.Err("Error while fetching data from DB for state %v, err: %v", stateName, err.Error())
		return data, err
	}

	log.Instance.Info("Data for state %v fetch from DB successfully", stateName)

	// Save data into redis if fetched from DB
	if !isFoundInRedis {
		redisData, _ := json.Marshal(data[0])
		err := redisDriverInstance.Set(stateName, redisData, constants.RedisTTL).Err()
		if err != nil {
			log.Instance.Err("Error while saving data in redis for state %v, err: %v", stateName, err.Error())
		} else {
			log.Instance.Info("Data saved successfully in redis for state %v", stateName)
		}
	}

	return data, nil
}

// Fetches all state covid data from DB if state name or coordinates not provided
func GetAllStateCovidData() ([]primitive.M, error) {
	log.Instance.Debug("GetAllStateCovidData is hit")

	var data []bson.M
	isFoundInRedis := true

	redisDriverInstance, redisDriverInstanceErr := drivers.GetRedisDriver()
	if redisDriverInstanceErr != nil {
		log.Instance.Err("Redis instance is down using DB to fetch data, err: %v", redisDriverInstanceErr.Error())
	} else {
		result, _ := redisDriverInstance.Get("AllStatesData").Bytes()

		if len(result) == 0 {
			isFoundInRedis = false
			log.Instance.Info("All states data not found in redis requesting from DB")
		} else {
			json.Unmarshal(result, &data)
			log.Instance.Info("All states data found in redis")
			return data, nil
		}
	}

	mongoDriverInstance, err := drivers.GetMongoDriver()
	if err != nil {
		log.Instance.Err("DB is down, err: %v", err.Error())
		return data, err
	}

	coll := mongoDriverInstance.Database(os.Getenv(constants.MongoDBName)).Collection("covid-state")
	cursor, err := coll.Find(context.TODO(), bson.M{}, options.Find().SetProjection(bson.M{"_id": 0}))
	if err != nil {
		log.Instance.Err("Error while fetching data from DB, err %v", err.Error())
		return data, err
	}

	err = cursor.All(context.TODO(), &data)
	if err != nil {
		log.Instance.Err("Error while reading converting data from DB, err %v", err.Error())
		return data, err
	}

	log.Instance.Info("Data for all states fetch from DB successfully")

	// Save data into redis if fetched from DB
	if !isFoundInRedis {
		redisData, err := json.Marshal(data)
		if err != nil {
			log.Instance.Err("Error while converting data for saving into redis, err %v", err.Error())
		}

		err = redisDriverInstance.Set("AllStatesData", redisData, constants.RedisTTL).Err()
		if err != nil {
			log.Instance.Err("Error while saving all states data in redis, err: %v", err.Error())
		} else {
			log.Instance.Info("All states data saved in redis")
		}
	}

	return data, nil
}

// Fetches all states covid data from 3rd party API, used by cronjob to sync DB
func GetAllStateCovidDataGovtApi() ([]byte, error) {
	log.Instance.Debug("GetAllStateCovidDataGovtApi is hit")

	body, err := utils.GetRequest(constants.CovidDataApi)
	if err != nil {
		log.Instance.Err("Data fetch failed from 3rd party API, err: %v", err.Error())
		return body, err
	} else {
		log.Instance.Info("Data fetch successfully from 3rd party API")
	}

	return body, nil
}

// Determines state using provided coordinates if found else return empty string
func GetStateNameUsingLatAndLong(latLang []string) (string, error) {
	log.Instance.Debug("GetStateNameUsingLatAndLong is hit")

	var stateName string
	var stateData schema.TCovidStateItems

	hereGeoCordinateApiMapper := map[string]string{
		"API_KEY": strings.TrimSpace(os.Getenv(constants.HereGeoAPIKey)),
		"LAT":     latLang[0],
		"LONG":    latLang[1],
	}

	buf := bytes.Buffer{}
	t := template.Must(template.New("").Parse(constants.HereGeoCordinateApi))
	t.Execute(&buf, hereGeoCordinateApiMapper)
	url := buf.String()

	response, err := utils.GetRequest(url)
	if err != nil {
		log.Instance.Err("State name request failed from 3rd party API for coordinates %v, %v, err: %v", latLang[0], latLang[1], err.Error())
		return stateName, err
	}

	json.Unmarshal(response, &stateData)

	if len(stateData.Items) > 0 {
		stateName = stateData.Items[0].Address.StateName
		stateName = strings.ReplaceAll(stateName, "&", "and")

		if stateName == "" {
			log.Instance.Info("No state found for coordinates %v, %v", latLang[0], latLang[1])
		} else {
			log.Instance.Info("State %v is located for coordinates %v, %v", stateName, latLang[0], latLang[1])
		}
	}

	return stateName, nil
}
