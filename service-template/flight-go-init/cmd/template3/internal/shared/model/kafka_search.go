package model

import (
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
)

// KafkaSearchRequest is the payload for the integrator search Kafka topic.
// Replace IntegratorFareRequests with concrete types when implementing search.
type KafkaSearchRequest struct {
	MandatoryRequest       common.MandatoryRequest `json:"mandatoryRequest"`
	IntegratorFareRequests []interface{}           `json:"integratorFareRequests"`
}
