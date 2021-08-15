package models

import (
	"context"
	"log"

	"github.com/MDAkramSiddiqui/sf-covid-api/app/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CovidState struct {
	ID               primitive.ObjectID `json:"_id" bson:"_id"`
	Name             string             `json:"name" bson:"name"`
	ActiveCases      string             `json:"active_cases" bson:"active_cases"`
	PositiveCases    string             `json:"positive_cases" bson:"positive_cases"`
	CuredCases       string             `json:"cured_cases" bson:"cured_cases"`
	DeathCases       string             `json:"death_cases" bson:"death_cases"`
	NewActiveCases   string             `json:"new_active_cases" bson:"new_active_cases"`
	NewPositiveCases string             `json:"new_positive_cases" bson:"new_positive_cases"`
	NewCuredCases    string             `json:"new_cured_cases" bson:"new_cured_cases"`
	NewDeathCases    string             `json:"new_death_cases" bson:"new_death_cases"`
	StateCode        int                `json:"state_code" bson:"state_code"`
}

func CreateState(state CovidState) error {
	client, err := services.GetMongoClient()
	if err != nil {
		return err
	}

	//Create a handle to the respective collection in the database.
	collection := client.Database("covid").Collection("covid-state")
	//Perform InsertOne operation & validate against the error.
	_, err = collection.InsertOne(context.TODO(), state)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}

func GetAllCollectionNames() []string {
	client, err := services.GetMongoClient()
	if err != nil {
		return nil
	}

	result, err := client.Database("covid").ListCollectionNames(context.TODO(), bson.D{{"options.capped", true}})
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func IsStateCollectionAvailable() bool {

	collectionNames := GetAllCollectionNames()
	flag := false

	for _, coll := range collectionNames {
		if coll == "covid-state" {
			flag = true
			break
		}
	}

	return flag
}
