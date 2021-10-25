package iceportalapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Status struct {
	ServiceLevel string  `json:"serviceLevel"`
	GpsStatus    string  `json:"gpsStatus"`
	Internet     string  `json:"internet"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	TileY        int64   `json:"tileY"`
	TileX        int64   `json:"tileX"`
	Series       string  `json:"series"`
	ServerTime   int64   `json:"serverTime"`
	Speed        float64 `json:"speed"`
	TrainType    string  `json:"trainType"`
	TZN          string  `json:"tzn"`
	WagonClass   string  `json:"wagonClass"`
	Connectivity struct {
		NextState            string `json:"nextState"`
		CurrentState         string `json:"currentState"`
		RemainingTimeSeconds int64  `json:"remainingTimeSeconds"`
	}
	Connection bool `json:"connection"`
	BAP        bool `json:"bap"`
}

func FetchStatus() (Status, error) {
	resp, err := http.Get("https://iceportal.de/api1/rs/status")
	if err != nil {
		return Status{}, err
	}

	// FIXME: better error handling if status is not 200!
	if resp.StatusCode >= 300 {
		return Status{}, fmt.Errorf("%d", resp.StatusCode)
	}

	// FIXME: better errors here
	if !strings.HasPrefix(resp.Header.Get("content-type"), "application/json") {
		return Status{}, fmt.Errorf("content-type is not application/json")
	}

	status := Status{}
	err = json.NewDecoder(resp.Body).Decode(&status)

	return status, err
}
