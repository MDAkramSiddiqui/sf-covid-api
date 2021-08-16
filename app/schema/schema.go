package schema

type TCovidState struct {
	Name             string `json:"state_name" bson:"state_name"`
	ActiveCases      string `json:"active" bson:"active"`
	PositiveCases    string `json:"positive" bson:"positive"`
	CuredCases       string `json:"cured" bson:"cured"`
	DeathCases       string `json:"death" bson:"death"`
	NewActiveCases   string `json:"new_active" bson:"new_active"`
	NewPositiveCases string `json:"new_positive" bson:"new_positive"`
	NewCuredCases    string `json:"new_cured" bson:"new_cured"`
	NewDeathCases    string `json:"new_death" bson:"new_death"`
	StateCode        string `json:"state_code" bson:"state_code"`
}

type TCovidStateAddress struct {
	StateCode string `json:"stateCode" bson:"stateCode"`
	StateName string `json:"state" bson:"state"`
}

type TCovidStateData struct {
	Address TCovidStateAddress `json:"address" bson:"address"`
}

type TCovidStateItems struct {
	Items []TCovidStateData `json:"items" bson:"items"`
}
