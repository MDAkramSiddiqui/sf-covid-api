package drivers

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/MDAkramSiddiqui/sf-covid-api/app/constants"
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

	// redisDriverInstanceError = nil
	//Perform connection creation operation only once.
	redisOnce.Do(func() {

		redisDB, _ := strconv.Atoi(os.Getenv(constants.RedisDB))
		driver := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%v:%v", os.Getenv(constants.RedisHost), os.Getenv(constants.RedisPort)),
			Password: os.Getenv(constants.RedisPassword),
			DB:       redisDB,
		})

		_, err := driver.Ping().Result()
		if err != nil {
			redisDriverInstanceError = err
			log.Instance.Err("Unable to connect to redis", err)
		} else {
			log.Instance.Info("Redis connection made successfully")
		}

		redisDriverInstance = driver
	})
	return redisDriverInstance, redisDriverInstanceError
}
