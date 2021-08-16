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
	"github.com/MDAkramSiddiqui/sf-covid-api/app/log"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/schema"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetStateCovidData(stateName string) primitive.M {
	log.Instance.Debug("GetStateCovidData is hit")

	var data bson.M
	isFoundInRedis := true

	stateName = strings.Title(strings.ToLower(strings.TrimSpace(strings.Trim(stateName, "\""))))

	redisDriverInstance, redisDriverInstanceErr := drivers.GetRedisDriver()
	if redisDriverInstanceErr != nil {
		log.Instance.Err("Redis instance is down using db to fetch data", redisDriverInstanceErr)
	} else {
		result, _ := redisDriverInstance.Get(stateName).Bytes()
		if len(result) == 0 {
			isFoundInRedis = false
			log.Instance.Info(fmt.Sprintf("State: %v data not found in redis requesting from DB", stateName))
		} else {
			json.Unmarshal(result, &data)
			log.Instance.Info(fmt.Sprintf("State: %v data found in redis", stateName))
			return data
		}
	}

	mongoDriverInstance, _ := drivers.GetMongoDriver()
	coll := mongoDriverInstance.Database(os.Getenv(constants.MongoDBName)).Collection("covid-state")

	_ = coll.FindOne(
		context.TODO(),
		bson.M{"name": stateName},
		options.FindOne().SetProjection(bson.M{"_id": 0}),
	).Decode(&data)

	if !isFoundInRedis {
		redisData, _ := json.Marshal(data)
		err := redisDriverInstance.Set(stateName, redisData, constants.RedisTTL).Err()
		if err != nil {
			log.Instance.Err(fmt.Sprintf("Error while saving %v state data in redis", stateName), err)
		} else {
			log.Instance.Info(fmt.Sprintf("State: %v data saved in redis", stateName))
		}
	}

	return data
}

func GetAllStateCovidData() []primitive.M {
	log.Instance.Debug("GetAllStateCovidData is hit")

	var data []bson.M
	isFoundInRedis := true

	redisDriverInstance, redisDriverInstanceErr := drivers.GetRedisDriver()
	if redisDriverInstanceErr != nil {
		log.Instance.Err("Redis instance is down using db to fetch data", redisDriverInstanceErr)
	} else {
		result, _ := redisDriverInstance.Get("AllStatesData").Bytes()
		if len(result) == 0 {
			isFoundInRedis = false
			log.Instance.Info("All states data not found in redis requesting from DB")
		} else {
			json.Unmarshal(result, &data)
			log.Instance.Info("All states data found in redis")
			return data
		}
	}

	mongoDriverInstance, _ := drivers.GetMongoDriver()
	coll := mongoDriverInstance.Database(os.Getenv(constants.MongoDBName)).Collection("covid-state")
	cursor, _ := coll.Find(context.TODO(), bson.M{}, options.Find().SetProjection(bson.M{"_id": 0}))
	_ = cursor.All(context.TODO(), &data)

	if !isFoundInRedis {
		redisData, _ := json.Marshal(data)
		err := redisDriverInstance.Set("AllStatesData", redisData, constants.RedisTTL).Err()
		if err != nil {
			log.Instance.Err("Error while saving all states data in redis", err)
		} else {
			log.Instance.Info("All states data saved in redis")
		}
	}

	return data
}

func GetAllStateCovidDataGovtApi() []byte {
	log.Instance.Debug("GetAllStateCovidDataGovtApi is hit")

	body, _ := utils.GetRequest(constants.CovidDataApi)
	return body
}

func GetStateNameUsingLatAndLong(latLang []string) string {
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

	response, _ := utils.GetRequest(url)
	json.Unmarshal(response, &stateData)

	if len(stateData.Items) > 0 {
		stateName = stateData.Items[0].Address.StateName
		stateName = strings.ReplaceAll(stateName, "&", "and")

		if stateName == "" {
			log.Instance.Info(fmt.Sprintf("Not state found for coordinates %v, %v", latLang[0], latLang[1]))
		} else {
			log.Instance.Info(fmt.Sprintf("State %v is located for coordinates %v, %v", stateName, latLang[0], latLang[1]))
		}

		return stateName
	}

	return stateName
}
