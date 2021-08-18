package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
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
func GetCovidDataByName(stateName string) (primitive.M, *utils.CustomErr) {
	log.Instance.Debug("GetCovidDataByName is hit")

	var result bson.M
	var isFoundInRedis bool
	var isValidStateName bool

	isFoundInRedis = true
	stateName = strings.Title(strings.ToLower(strings.TrimSpace(strings.Trim(stateName, "\""))))

	if stateName == "" {
		log.Instance.Info("State name is not provided")
		return result, &utils.CustomErr{}
	} else {
		isValidStateName = isValidState(stateName)

		if !isValidStateName {
			return result, &utils.CustomErr{Err: errors.New("state name provided is invalid"), StatusCode: http.StatusBadRequest}
		}
	}

	redisDriverInstance, redisDriverInstanceErr := drivers.GetRedisDriver()
	if redisDriverInstanceErr != nil {
		log.Instance.Err("Redis instance is down using DB to fetch data, err: %v", redisDriverInstanceErr.Error())
	} else {
		redisResult, _ := redisDriverInstance.Get(stateName).Bytes()

		if len(result) == 0 {
			isFoundInRedis = false
			log.Instance.Info("Data for state: %v not found in redis requesting from DB", stateName)
		} else {
			json.Unmarshal(redisResult, &result)
			log.Instance.Info("Data for state: %v found in redis", stateName)
			return result, &utils.CustomErr{}
		}
	}

	mongoDriverInstance, err := drivers.GetMongoDriver()
	if err != nil {
		log.Instance.Err("DB is down, err: %v", stateName, err.Error())
		return result, &utils.CustomErr{Err: err, StatusCode: http.StatusInternalServerError}
	}

	coll := mongoDriverInstance.Database(os.Getenv(constants.MongoDBName)).Collection("covid-state")

	err = coll.FindOne(
		context.TODO(),
		bson.M{"name": stateName},
		options.FindOne().SetProjection(bson.M{"_id": 0}),
	).Decode(&result)

	if err != nil {
		log.Instance.Err("Error while fetching data from DB for state %v, err: %v", stateName, err.Error())
		return result, &utils.CustomErr{Err: err, StatusCode: http.StatusInternalServerError}
	}

	log.Instance.Info("Data for state %v fetch from DB successfully", stateName)

	// Save data into redis if fetched from DB
	if !isFoundInRedis {
		redisData, _ := json.Marshal(result)
		err := redisDriverInstance.Set(stateName, redisData, constants.DefaultRedisTTL).Err()
		if err != nil {
			log.Instance.Err("Error while saving data in redis for state %v, err: %v", stateName, err.Error())
		} else {
			log.Instance.Info("Data saved successfully in redis for state %v", stateName)
		}
	}

	return result, &utils.CustomErr{}
}

// Fetches covid data of a state using provided coordinates
func GetCovidDataByCoordinates(latlngStr string) (primitive.M, *utils.CustomErr) {
	log.Instance.Debug("GetCovidDataByCoordinates is hit")

	var stateName string

	latLang := strings.Split(latlngStr, ",")

	if latlngStr == "" {
		log.Instance.Info("Latitude and longitude are not provided")
		return nil, &utils.CustomErr{}
	}

	if len(latLang) != 2 {
		log.Instance.Info("Latitude and longitude are invalid")
		return nil, &utils.CustomErr{Err: errors.New("invalid coordinates provided"), StatusCode: http.StatusBadRequest}
	}

	latLang[0], latLang[1] = strings.TrimSpace(latLang[0]), strings.TrimSpace(latLang[1])

	if len(latLang[0]) <= 0 || len(latLang[1]) <= 0 {
		log.Instance.Info("Latitude and longitude are invalid")
		return nil, &utils.CustomErr{Err: errors.New("invalid coordinates provided"), StatusCode: http.StatusBadRequest}
	}

	log.Instance.Info("Latitude and longitude provided are %v, %v", latLang[0], latLang[1])
	stateName, latLangErr := GetStateNameUsingLatAndLong(latLang)
	if latLangErr.Err != nil {
		return nil, latLangErr
	}

	log.Instance.Info("Fetching data for state %v", stateName)
	result, covidDataErr := GetCovidDataByName(stateName)
	if covidDataErr.Err != nil {
		return result, covidDataErr
	}

	return result, &utils.CustomErr{}
}

// Fetches all state covid data from DB if state name or coordinates not provided
func GetAllStateCovidData() ([]primitive.M, *utils.CustomErr) {
	log.Instance.Debug("GetAllStateCovidData is hit")

	var result []bson.M
	isFoundInRedis := true

	redisDriverInstance, redisDriverInstanceErr := drivers.GetRedisDriver()
	if redisDriverInstanceErr != nil {
		log.Instance.Err("Redis instance is down using DB to fetch data, err: %v", redisDriverInstanceErr.Error())
	} else {
		redisResult, _ := redisDriverInstance.Get("AllStatesData").Bytes()

		if len(result) == 0 {
			isFoundInRedis = false
			log.Instance.Info("All states data not found in redis requesting from DB")
		} else {
			json.Unmarshal(redisResult, &result)
			log.Instance.Info("All states data found in redis")
			return result, &utils.CustomErr{}
		}
	}

	mongoDriverInstance, mongoDriverInstanceErr := drivers.GetMongoDriver()
	if mongoDriverInstanceErr != nil {
		log.Instance.Err("DB is down, err: %v", mongoDriverInstanceErr.Error())
		return result, &utils.CustomErr{Err: mongoDriverInstanceErr, StatusCode: http.StatusInternalServerError}
	}

	coll := mongoDriverInstance.Database(os.Getenv(constants.MongoDBName)).Collection("covid-state")
	cursor, mongoDataFetchErr := coll.Find(context.TODO(), bson.M{}, options.Find().SetProjection(bson.M{"_id": 0}))
	if mongoDataFetchErr != nil {
		log.Instance.Err("Error while fetching data from DB, err %v", mongoDataFetchErr.Error())
		return result, &utils.CustomErr{Err: mongoDataFetchErr, StatusCode: http.StatusInternalServerError}
	}

	mongoDataReadErr := cursor.All(context.TODO(), &result)
	if mongoDataReadErr != nil {
		log.Instance.Err("Error while reading converting data from DB, err %v", mongoDataReadErr.Error())
		return result, &utils.CustomErr{Err: mongoDataReadErr, StatusCode: http.StatusInternalServerError}
	}

	log.Instance.Info("Data for all states fetch from DB successfully")

	// Save data into redis if fetched from DB
	if !isFoundInRedis {
		redisData, resultCovertErr := json.Marshal(result)
		if resultCovertErr != nil {
			log.Instance.Err("Error while converting data for saving into redis, err %v", resultCovertErr.Error())
		}

		redisResultSetErr := redisDriverInstance.Set("AllStatesData", redisData, constants.DefaultRedisTTL).Err()
		if redisResultSetErr != nil {
			log.Instance.Err("Error while saving all states data in redis, err: %v", redisResultSetErr.Error())
		} else {
			log.Instance.Info("All states data saved in redis")
		}
	}

	return result, &utils.CustomErr{}
}

// Fetches all states covid data from 3rd party API, used by cronjob to sync DB
func GetAllStateCovidDataGovtApi() ([]byte, *utils.CustomErr) {
	log.Instance.Debug("GetAllStateCovidDataGovtApi is hit")

	allStatesCovidData, allStatesCovidDataErr := utils.GetRequest(constants.CovidDataApi)
	if allStatesCovidDataErr.Err != nil {
		log.Instance.Err("Data fetch failed from 3rd party API, err: %v", allStatesCovidDataErr.Message())
		return allStatesCovidData, &utils.CustomErr{Err: allStatesCovidDataErr.Err, StatusCode: allStatesCovidDataErr.StatusCode}
	} else {
		log.Instance.Info("Data fetch successfully from 3rd party API")
	}

	return allStatesCovidData, &utils.CustomErr{}
}

// Determines state using provided coordinates if found else return empty string
func GetStateNameUsingLatAndLong(latLang []string) (string, *utils.CustomErr) {
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

	response, getReqErr := utils.GetRequest(url)
	if getReqErr.Err != nil {
		log.Instance.Err("State name request failed from 3rd party API for coordinates %v, %v, err: %v", latLang[0], latLang[1], getReqErr.Message())
		return stateName, &utils.CustomErr{Err: getReqErr.Err, StatusCode: getReqErr.StatusCode}
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

	return stateName, &utils.CustomErr{}
}

// check state is a valid state or not
func isValidState(stateName string) bool {
	for _, item := range constants.AllValidStates {
		if item == stateName {
			return true
		}
	}
	return false
}
