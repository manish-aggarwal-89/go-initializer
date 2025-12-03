package bau

import (
	"encoding/json"

	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/bau"
)

const TransactionTypeSearch = "Search"

const (
	SearchingHealthStatusSuccess SearchHealthStatus = iota
	SearchingHealthStatusSearchFailed
	SearchingHealthStatusFlightNotFound
	SearchingHealthStatusInvalidRoute
	SearchingHealthStatusGetFareFailed
	SearchingHealthStatusUnmapped
)

type SearchHealthStatus int

func (status SearchHealthStatus) MarshalJSON() ([]byte, error) {
	switch status {
	case SearchingHealthStatusSuccess:
		return json.Marshal("SUCCESS")
	case SearchingHealthStatusSearchFailed:
		return json.Marshal("SEARCH_FAILED")
	case SearchingHealthStatusFlightNotFound:
		return json.Marshal("FLIGHT_NOT_FOUND")
	case SearchingHealthStatusInvalidRoute:
		return json.Marshal("INVALID_ROUTE")
	case SearchingHealthStatusGetFareFailed:
		return json.Marshal("GET_FARE_FAILED")
	default:
		return json.Marshal("UNMAPPED")
	}
}

func (status SearchHealthStatus) String() string {
	switch status {
	case SearchingHealthStatusSuccess:
		return "SUCCESS"
	case SearchingHealthStatusSearchFailed:
		return "SEARCH_FAILED"
	case SearchingHealthStatusFlightNotFound:
		return "FLIGHT_NOT_FOUND"
	case SearchingHealthStatusInvalidRoute:
		return "INVALID_ROUTE"
	case SearchingHealthStatusGetFareFailed:
		return "GET_FARE_FAILED"
	default:
		return "UNMAPPED"
	}
}

type SearchLogHealthStatus int

const (
	SearchLogHealthStatusSuccess SearchLogHealthStatus = iota
	SearchLogHealthStatusFlightNotFoundError
	SearchLogHealthStatusUnknownError
)

func (status SearchLogHealthStatus) MarshalJSON() ([]byte, error) {
	switch status {
	case SearchLogHealthStatusSuccess:
		return json.Marshal("SUCCESS")
	case SearchLogHealthStatusFlightNotFoundError:
		return json.Marshal("FLIGHT_NOT_FOUND")
	default:
		return json.Marshal("UNKNOWN_ERROR")
	}
}

func (s SearchLogHealthStatus) String() string {
	switch s {
	case SearchLogHealthStatusSuccess:
		return "SUCCESS"
	case SearchLogHealthStatusFlightNotFoundError:
		return "FLIGHT_NOT_FOUND"
	default:
		return "UNKNOWN_ERROR"
	}
}

type (
	AnalyticSearchToken struct {
		bau.AnalyticsToken

		Origin      string `json:"origin"`
		Destination string `json:"destination"`
		TripType    string `json:"tripType"`
	}

	AnalyticSearchRequest struct {
		RequestID                string             `json:"request_id"`
		Distribution             string             `json:"distribution"`
		Supplier                 string             `json:"supplier"`
		Airlines                 []string           `json:"airlines"`
		TransactionDate          string             `json:"transaction_date"`
		AirlineProcessTimeDetail map[string]int64   `json:"airline_process_time_detail"`
		AirlineProcessTime       int64              `json:"airline_process_time"`
		ProcessTime              int64              `json:"process_time"`
		StoreID                  string             `json:"store_id"`
		ChannelID                string             `json:"channel_id"`
		Origin                   string             `json:"origin"`
		Destination              string             `json:"destination"`
		TripType                 string             `json:"trip_type"`
		SearchHealthStatus       SearchHealthStatus `json:"health"`
	}
)
