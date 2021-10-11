package iceportalapi_test

import (
	"net/http"
	"reflect"
	"testing"

	iceportalapi "github.com/craftamap/iceportal-api"
	"github.com/jarcoal/httpmock"
)

const JSON_RESPONSE = `{
    "connection": true,
    "servicelevel": "AVAILABLE_SERVICE",
    "internet": "HIGH",
    "gpsStatus": "VALID",
    "tileY": 303,
    "tileX": 216,
    "series": "011",
    "latitude": 52.766562,
    "longitude": 10.251847,
    "serverTime": 1603913237508,
    "speed": 169,
    "trainType": "ICE",
    "tzn": "Tz1191",
    "wagonClass": "SECOND",
    "connectivity": {
        "currentState": "HIGH",
        "nextState": "WEAK",
        "remainingTimeSeconds": 58
    }
}`

func TestFetchStatus(t *testing.T) {
	t.Run("correct behaviour", func(t *testing.T) {
		httpmock.Activate()
		httpmock.RegisterResponder("GET", "https://iceportal.de/api1/rs/status", func(*http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, JSON_RESPONSE)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

		status, err := iceportalapi.FetchStatus()
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(
			status,
			iceportalapi.Status{
				ServiceLevel: "AVAILABLE_SERVICE",
				GpsStatus:    "VALID",
				Internet:     "HIGH",
				Latitude:     52.766562,
				Longitude:    10.251847,
				TileY:        303,
				TileX:        216,
				Series:       "011",
				ServerTime:   1603913237508,
				Speed:        169,
				TrainType:    "ICE",
				TZN:          "Tz1191",
				WagonClass:   "SECOND",
				Connectivity: struct {
					NextState            string `json:"nextState"`
					CurrentState         string `json:"currentState"`
					RemainingTimeSeconds int64  `json:"remainingTimeSeconds"`
				}{RemainingTimeSeconds: 58, NextState: "WEAK", CurrentState: "HIGH"},
				Connection: true,
				BAP:        false,
			}) {
			t.Error()
		}
		httpmock.DeactivateAndReset()
	})
}
