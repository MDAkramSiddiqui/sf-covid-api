package constants

import "time"

// Api Uri constants
const (
	CovidDataApi        = "https://www.mohfw.gov.in/data/datanew.json"
	HereGeoCordinateApi = `https://revgeocode.search.hereapi.com/v1/revgeocode?at={{.LAT}},{{.LONG}}&apiKey={{.API_KEY}}`
)

// Environment Variables Key Constants
const (
	Env                 = "ENV"
	MongoDBUrl          = "MONGO_DB_URL"
	MongoDBPassword     = "MONGO_DB_PASSWORD"
	MongoDBName         = "MONGO_DB_NAME"
	HereGeoAPIKey       = "HERE_GEO_API_KEY"
	Port                = "PORT"
	RedisPort           = "REDIS_PORT"
	RedisHost           = "REDIS_HOST"
	RedisDB             = "REDIS_DB"
	RedisPassword       = "REDIS_PASSWORD"
	LogLevel            = "LOG_LEVEL"
	StateDataCronPeriod = "STATE_DATA_CRON"
)

// Severity levels.
const (
	DebugLevel int = iota
	InfoLevel
	WarnLevel
	ErrLevel
	FatalLevel
)

// Possible App Environments
const (
	Development = "DEVELOPMENT"
	Production  = "PRODUCTION"
)

// Some default parameters
const (
	DefaultRedisTTL            = time.Minute * 30
	DefaultStateDataCronPeriod = "*/30 * * * *"
)

// all states name
var (
	AllValidStates = []string{"Andaman and Nicobar Islands", "Andhra Pradesh", "Arunachal Pradesh", "Assam", "Bihar", "Chandigarh", "Chhattisgarh", "Dadra and Nagar Haveli and Daman and Diu", "Delhi", "Goa", "Gujarat", "Haryana", "Himachal Pradesh", "Jammu and Kashmir", "Jharkhand", "Karnataka", "Kerala", "Ladakh", "Lakshadweep", "Madhya Pradesh", "Maharashtra", "Manipur", "Meghalaya", "Mizoram", "Nagaland", "Odisha", "Puducherry", "Punjab", "Rajasthan", "Sikkim", "Tamil Nadu", "Telangana", "Tripura", "Uttarakhand", "Uttar Pradesh", "West Bengal", "India"}
)
