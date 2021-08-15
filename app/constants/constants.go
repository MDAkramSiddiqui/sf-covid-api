package constants

// Api Uri constants
const (
	CovidDataApi        = "https://www.mohfw.gov.in/data/datanew.json"
	HereGeoCordinateApi = `https://revgeocode.search.hereapi.com/v1/revgeocode?at={{.LAT}}%2C{{.LONG}}&apiKey={{.API_KEY}}`
)

// Environment Variables Key Constants
const (
	Env             = "ENV"
	MongoDBUrl      = "MONGO_DB_URL"
	MongoDBPassword = "MONGO_DB_PASSWORD"
	MongoDBName     = "MONGO_DB_NAME"
	HereGeoAPIKey   = "HERE_GEO_API_KEY"
)

// Severity levels.
const (
	Debug int = iota
	Info
	Warn
	Err
	Fatal
)

// Possible App Environments
const (
	DEVELOPMENT = "DEVELOPMENT"
	PRODUCTION  = "PRODUCTION"
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
