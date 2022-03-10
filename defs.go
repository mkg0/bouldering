package main

type BoulderStudio struct {
	name         string
	shiftModelId int
	clientId     int
	tariffId     int
}

var Bouldergarten = BoulderStudio{
	name:         "Bouldergarten",
	shiftModelId: 67359814,
	clientId:     66563713,
	tariffId:     68083827,
}

var Boulderklub = BoulderStudio{
	name:         "Boulderklub",
	shiftModelId: 67361411,
	clientId:     67359314,
	tariffId:     67783230,
}

var boulderStudios = []BoulderStudio{
	Bouldergarten,
	Boulderklub,
}

type Slot struct {
	Selector                      []interface{} `json:"selector"`
	BookableFrom                  int           `json:"bookableFrom"`
	State                         string        `json:"state"`
	BookableUntilDuration         int           `json:"bookableUntilDuration"`
	MinCourseParticipantCount     int           `json:"minCourseParticipantCount"`
	MaxCourseParticipantCount     int           `json:"maxCourseParticipantCount"`
	CurrentCourseParticipantCount int           `json:"currentCourseParticipantCount"`
	DateList                      []struct {
		Start int64 `json:"start"`
		End   int64 `json:"end"`
	} `json:"dateList"`
}

type Profile struct {
	DateOfBirth string
	Address     string
	PostCode    string
	City        string
	Phone       string
	Name        string
	Surname     string
	Email       string
	USCid       string
}

type persistData struct {
	Profiles             []Profile
	ApiEndpoint          string
	RemoteBookingEnabled bool
}

var global = persistData{}
var dayCount = 7
var defaultRetryIntervalSeconds int = 10
