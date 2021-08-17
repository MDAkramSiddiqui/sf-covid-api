package constants

import "time"

// Api Uri constants
const (
	CovidDataApi        = "https://www.mohfw.gov.in/data/datanew.json"
	HereGeoCordinateApi = `https://revgeocode.search.hereapi.com/v1/revgeocode?at={{.LAT}},{{.LONG}}&apiKey={{.API_KEY}}`
)

// Environment Variables Key Constants
const (
	Env             = "ENV"
	MongoDBUrl      = "MONGO_DB_URL"
	MongoDBPassword = "MONGO_DB_PASSWORD"
	MongoDBName     = "MONGO_DB_NAME"
	HereGeoAPIKey   = "HERE_GEO_API_KEY"
	Port            = "PORT"
	RedisPort       = "REDIS_PORT"
	RedisHost       = "REDIS_HOST"
	RedisDB         = "REDIS_DB"
	RedisPassword   = "REDIS_PASSWORD"
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

const (
	RedisTTL = time.Minute * 30
)
