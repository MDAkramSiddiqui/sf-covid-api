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

// func StateNames() map[string]string {
// 	return map[string]string{
// 		"25": "Andhra Pradesh",
// 		"35": "Andaman and Nicobar Islands",
// 		"12": "Arunachal Pradesh",
// 		"18": "Assam",
// 		"10": "Bihar",
// 		"04": "Chandigarh",
// 		"22": "Chhattisgarh",
// 		"26": "Dadra and Nagar Haveli and Daman and Diu",
// 		"07": "Delhi",
// 		"30": "Goa",
// 		"0":  "",
// 		"0":  "",
// 		"0":  "",
// 		"0":  "",
// 		"0":  "",
// 		"0":  "",
// 		"0":  "",
// 		"0":  "",
// 		"0":  "",
// 		"0":  "",
// 		"0":  "",
// 		"0":  "",
// 		"0":  "",
// 		"0":  "",
// 		"0":  "",
// 		"0":  "",
// 		"0":  "",
// 		"0":  "",
// 		"0":  "",
// 		"0":  "",
// 		"0":  "",
// 		"0":  "",
// 		"0":  "",
// 	}
// }
