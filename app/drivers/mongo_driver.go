package drivers

import (
	"context"
	"os"
	"strings"
	"sync"

	"github.com/MDAkramSiddiqui/sf-covid-api/app/constants"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/* Used to create a singleton object of MongoDB client.
Initialized and exposed through  GetMongoDriver().*/
var mongoDriverInstance *mongo.Client

//Used during creation of singleton client object in GetMongoDriver().
var mongoDriverInstanceError error

//Used to execute client creation procedure only once.
var mongoOnce sync.Once

//GetMongoDriver - Return mongodb connection to work with
func GetMongoDriver() (*mongo.Client, error) {
	log.Instance.Debug("GetMongoDriver is hit")

	mongoDriverInstanceError = nil
	//Perform connection creation operation only once.
	mongoOnce.Do(func() {

		connectURI := strings.Replace(os.Getenv(constants.MongoDBUrl), "<password>", os.Getenv(constants.MongoDBPassword), 1)
		driver, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectURI))
		if err != nil {
			mongoDriverInstanceError = err
			log.Instance.Fatal("Unable to connect to database", err)
		}
		err = driver.Ping(context.Background(), nil)
		if err != nil {
			mongoDriverInstanceError = err
		}
		log.Instance.Info("Connection made successfully")
		mongoDriverInstance = driver
	})
	return mongoDriverInstance, mongoDriverInstanceError
}
