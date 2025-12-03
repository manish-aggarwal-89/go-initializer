package bau

import (
	"encoding/json"

	modelBau "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/bau"
)

const TransactionTypeBooking = "Booking"

const (
	BookingHealthStatusSuccess BookingHealthStatus = iota
	BookingHealthStatusBookingFailed
	BookingHealthStatusNoSeat
	BookingHealthStatusUnmapped
)

type BookingHealthStatus int

func (status BookingHealthStatus) MarshalJSON() ([]byte, error) {
	switch status {
	case BookingHealthStatusSuccess:
		return json.Marshal("SUCCESS")
	case BookingHealthStatusBookingFailed:
		return json.Marshal("BOOKING_FAILED")
	case BookingHealthStatusNoSeat:
		return json.Marshal("NO_SEAT")
	default:
		return json.Marshal("UNMAPPED")
	}
}

func (s BookingHealthStatus) String() string {
	switch s {
	case BookingHealthStatusSuccess:
		return "SUCCESS"
	case BookingHealthStatusBookingFailed:
		return "BOOKING_FAILED"
	case BookingHealthStatusNoSeat:
		return "NO_SEAT"
	default:
		return "UNMAPPED"
	}
}

type BookingLogHealthStatus int

const (
	BookingLogHealthStatusSuccess BookingLogHealthStatus = iota
	BookingLogHealthStatusBookingFailed
	BookingLogHealthStatusNoSeat
	BookingLogHealthStatusUnknownError
)

func (status BookingLogHealthStatus) MarshalJSON() ([]byte, error) {
	switch status {
	case BookingLogHealthStatusSuccess:
		return json.Marshal("SUCCESS")
	case BookingLogHealthStatusBookingFailed:
		return json.Marshal("BOOKING_FAILED")
	case BookingLogHealthStatusNoSeat:
		return json.Marshal("NO_SEAT")
	default:
		return json.Marshal("UNKNOWN_ERROR")
	}
}

func (s BookingLogHealthStatus) String() string {
	switch s {
	case BookingLogHealthStatusSuccess:
		return "SUCCESS"
	case BookingLogHealthStatusBookingFailed:
		return "BOOKING_FAILED"
	case BookingLogHealthStatusNoSeat:
		return "NO_SEAT"
	default:
		return "UNKNOWN_ERROR"
	}
}

type AnalyticBookingToken struct {
	modelBau.AnalyticsToken

	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	TripType    string `json:"tripType"`
	RetryCount  int    `json:"retryCount"`
}

type AnalyticBookingRequest struct {
	RequestID                string              `json:"request_id"`
	Distribution             string              `json:"distribution"`
	Supplier                 string              `json:"supplier"`
	Airlines                 []string            `json:"airlines"`
	TransactionDate          string              `json:"transaction_date"`
	AirlineProcessTimeDetail map[string]int64    `json:"airline_process_time_detail"`
	AirlineProcessTime       int64               `json:"airline_process_time"`
	ProcessTime              int64               `json:"process_time"`
	StoreID                  string              `json:"store_id"`
	ChannelID                string              `json:"channel_id"`
	Origin                   string              `json:"origin"`
	Destination              string              `json:"destination"`
	TripType                 string              `json:"trip_type"`
	RetryCount               int                 `json:"retry_count"`
	BookingHealthStatus      BookingHealthStatus `json:"health"`
}
