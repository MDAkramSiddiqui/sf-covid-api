package schema

// schema for 3rd party covid data api
type TCovidState struct {
	Name                string `json:"state_name" bson:"state_name"`
	ActiveCases         string `json:"active" bson:"active"`
	PositiveCases       string `json:"positive" bson:"positive"`
	CuredCases          string `json:"cured" bson:"cured"`
	DeathCases          string `json:"death" bson:"death"`
	LatestActiveCases   string `json:"new_active" bson:"new_active"`
	LatestPositiveCases string `json:"new_positive" bson:"new_positive"`
	LatestCuredCases    string `json:"new_cured" bson:"new_cured"`
	LatestDeathCases    string `json:"new_death" bson:"new_death"`
	StateCode           string `json:"state_code" bson:"state_code"`
}

// schema for 3rd party geo-coordinates Api that returns state related data
type TCovidStateAddress struct {
	StateCode string `json:"stateCode" bson:"stateCode"`
	StateName string `json:"state" bson:"state"`
}

// schema for 3rd party geo-coordinates Api that returns state related data
type TCovidStateData struct {
	Address TCovidStateAddress `json:"address" bson:"address"`
}

// schema for 3rd party geo-coordinates Api that returns state related data
type TCovidStateItems struct {
	Items []TCovidStateData `json:"items" bson:"items"`
}

// schema for default response model
type TDefaultResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
	Message interface{} `json:"message"`
}
