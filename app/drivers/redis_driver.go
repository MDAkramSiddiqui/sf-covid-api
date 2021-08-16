package drivers

import (
	"sync"

	"github.com/MDAkramSiddiqui/sf-covid-api/app/log"
	"github.com/go-redis/redis"
)

/* Used to create a singleton object of MongoDB client.
Initialized and exposed through  GetRedisDriver().*/
var redisDriverInstance *redis.Client

//Used during creation of singleton client object in GetRedisDriver().
var redisDriverInstanceError error

//Used to execute client creation procedure only once.
var redisOnce sync.Once

//GetRedisDriver - Return mongodb connection to work with
func GetRedisDriver() (*redis.Client, error) {
	log.Instance.Debug("GetRedisDriver is hit")

	redisDriverInstanceError = nil
	//Perform connection creation operation only once.
	redisOnce.Do(func() {

		driver := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})

		_, err := driver.Ping().Result()
		if err != nil {
			redisDriverInstanceError = err
			log.Instance.Fatal("Unable to connect to redis", err)
		}

		log.Instance.Info("Redis connection made successfully")
		redisDriverInstance = driver
	})
	return redisDriverInstance, redisDriverInstanceError
}
