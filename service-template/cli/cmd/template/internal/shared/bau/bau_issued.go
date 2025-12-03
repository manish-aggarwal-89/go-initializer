package bau

import (
	"encoding/json"

	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/bau"
)

const TransactionTypeIssuance = "Issued"

type IssuanceHealthStatus int

const (
	IssuanceHealthStatusSuccess IssuanceHealthStatus = iota
	IssuanceHealthStatusIssuanceFailed
	IssuanceHealthStatusUnmapped
)

func (status IssuanceHealthStatus) MarshalJSON() ([]byte, error) {
	switch status {
	case IssuanceHealthStatusSuccess:
		return json.Marshal("SUCCESS")
	case IssuanceHealthStatusIssuanceFailed:
		return json.Marshal("ISSUED_FAILED")
	default:
		return json.Marshal("UNMAPPED")
	}
}

func (status IssuanceHealthStatus) String() string {
	switch status {
	case IssuanceHealthStatusSuccess:
		return "SUCCESS"
	case IssuanceHealthStatusIssuanceFailed:
		return "ISSUED_FAILED"
	default:
		return "UNMAPPED"
	}
}

type IssuanceLogHealthStatus int

const (
	IssuanceLogHealthStatusSuccess IssuanceLogHealthStatus = iota
	IssuanceLogHealthStatusIssuanceFailed
	IssuanceLogHealthStatusUnknownError
)

func (status IssuanceLogHealthStatus) MarshalJSON() ([]byte, error) {
	switch status {
	case IssuanceLogHealthStatusSuccess:
		return json.Marshal("SUCCESS")
	case IssuanceLogHealthStatusIssuanceFailed:
		return json.Marshal("ISSUED_FAILED")
	default:
		return json.Marshal("UNKNOWN_ERROR")
	}
}

func (status IssuanceLogHealthStatus) String() string {
	switch status {
	case IssuanceLogHealthStatusSuccess:
		return "SUCCESS"
	case IssuanceLogHealthStatusIssuanceFailed:
		return "ISSUED_FAILED"
	default:
		return "UNKNOWN_ERROR"
	}
}

type AnalyticsIssuanceToken struct {
	BookingCode string
	bau.AnalyticsToken
}

type AnalyticsIssuanceRequest struct {
	RequestId                string               `json:"request_id"`
	Distribution             string               `json:"distribution"`
	Supplier                 string               `json:"supplier"`
	BookingCode              string               `json:"booking_code"`
	Airlines                 []string             `json:"airlines"`
	TransactionDate          string               `json:"transaction_date"`
	Status                   IssuanceHealthStatus `json:"health"`
	AirlineProcessTimeDetail map[string]int64     `json:"airline_process_time_detail"`
	AirlineProcessTime       int64                `json:"airline_process_time"`
	StoreId                  string               `json:"store_id"`
	ChannelId                string               `json:"channel_id"`
}
