package services

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/MDAkramSiddiqui/sf-covid-api/app/constants"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConnectionURI struct {
	PASSWORD string
	DB_NAME  string
}

/* Used to create a singleton object of MongoDB client.
Initialized and exposed through  GetMongoClient().*/
var clientInstance *mongo.Client

//Used during creation of singleton client object in GetMongoClient().
var clientInstanceError error

//Used to execute client creation procedure only once.
var mongoOnce sync.Once

//GetMongoClient - Return mongodb connection to work with
func GetMongoClient() (*mongo.Client, error) {
	clientInstanceError = nil
	//Perform connection creation operation only once.
	mongoOnce.Do(func() {

		connectURI := strings.Replace(os.Getenv(constants.MongoDBUrl), "<password>", os.Getenv(constants.MongoDBPassword), 1)
		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectURI))
		if err != nil {
			clientInstanceError = err
			fmt.Printf("Unable to connect to database : %v", err)
		}
		err = client.Ping(context.Background(), nil)
		if err != nil {
			clientInstanceError = err
		}
		fmt.Println("Connection Made successfully")
		clientInstance = client
	})
	return clientInstance, clientInstanceError
}
