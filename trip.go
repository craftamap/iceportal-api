package iceportalapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type StopInfo struct {
	ScheduledNext     string `json:"scheduledNext"`
	ActualNext        string `json:"actualNext"`
	ActualLast        string `json:"actualLast"`
	ActualLastStarted string `json:"actualLastStarted"`
	FinalStationName  string `json:"finalStationName"`
	FinalStationEvaNr string `json:"finalStationEvaNr"`
}

type Station struct {
	Code           interface{} `json:"code"`
	EvaNr          string      `json:"evaNr"`
	Name           string      `json:"name"`
	Geocoordinates struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"geocoordinates"`
}

type Timetable struct {
	ArrivalDelay            string `json:"arrivalDelay"`
	DepartureDelay          string `json:"departureDelay"`
	ScheduledDepartureTime  int64  `json:"scheduledDepartureTime"`
	ActualDepartureTime     int64  `json:"actualDepartureTime"`
	ShowActualDepartureTime bool   `json:"showActualDepartureTime"`
	ShowActualArrivalTime   bool   `json:"showActualArrivalTime"`
	ScheduledArrivalTime    int64  `json:"scheduledArrivalTime"`
	ActualArrivalTime       int64  `json:"actualArrivalTime"`
}

type Track struct {
	Scheduled string `json:"scheduled"`
	Actual    string `json:"actual"`
}
type StopsInfo struct {
	PositionStatus    string `json:"positionStatus"`
	Status            int64  `json:"status"`
	Distance          int64  `json:"distance"`
	DistanceFromStart int64  `json:"distanceFromStart"`
	Passed            bool   `json:"passed"`
}

type Stop struct {
	Station   Station   `json:"station"`
	Timetable Timetable `json:"timetable"`
	Track     Track     `json:"track"`
	Info      StopsInfo `json:"info"`
}

type Trip struct {
	TripDate             string   `json:"tripDate"`
	TrainType            string   `json:"trainType"`
	VZN                  string   `json:"vzn"`
	StopInfo             StopInfo `json:"stopInfo"`
	Stops                []Stop   `json:"stops"`
	ActualPosition       int64    `json:"actualPosition"`
	DistanceFromLastStop int64    `json:"distanceFromLastStop"`
	TotalDistance        int64    `json:"totalDistance"`
}

type TripInfo struct {
	Connection    interface{} `json:"connection"`
	SelectedRoute struct {
		Mobility     interface{} `json:"mobility"`
		ConflictInfo struct {
			Text   interface{} `json:"text"`
			Status string      `json:"status"`
		} `json:"conflictInfo"`
	} `json:"selectedRoute"`
	Trip Trip `json:"trip"`
}

func FetchTrip() (TripInfo, error) {
	resp, err := http.Get("https://iceportal.de/api1/rs/tripInfo/trip")
	if err != nil {
		return TripInfo{}, err
	}

	// FIXME: better error handling if status is not 200!
	if resp.StatusCode >= 300 {
		return TripInfo{}, fmt.Errorf("%d", resp.StatusCode)
	}
	// FIXME: better errors here
	if !strings.HasPrefix(resp.Header.Get("content-type"), "application/json") {
		return TripInfo{}, fmt.Errorf("content-type is not application/json")
	}

	trip := TripInfo{}
	err = json.NewDecoder(resp.Body).Decode(&trip)

	return trip, err
}
