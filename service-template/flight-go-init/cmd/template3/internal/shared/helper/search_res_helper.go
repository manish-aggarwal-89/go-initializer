package helper

import (
	"errors"
	"{{MODULE_NAME}}/internal/shared"

	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	fareRS "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/fare"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/kafka"
)

func ConstructKafkaMessageResponseSearch(mr common.MandatoryRequest, res fareRS.FlightIntegratorSearchResponse, err error) common.BaseResponse[kafka.KafkaSearchResponse] {
	kafkaRes := kafka.KafkaSearchResponse{
		MandatoryRequest:               mr,
		FlightIntegratorSearchResponse: res,
	}
	if err != nil {
		if errors.Is(err, shared.ErrInvalidRoute) {
			return common.ConstructResponse[kafka.KafkaSearchResponse](shared.INVALID_ROUTE_RESPONSE_CODE, shared.INVALID_ROUTE_RESPONSE_CODE, []string{err.Error()}, kafkaRes)
		}
		if errors.Is(err, shared.ErrFlightNotFound) {
			return common.ConstructResponse[kafka.KafkaSearchResponse](shared.FLIGHT_NOT_FOUND_RESPONSE_CODE, shared.FLIGHT_NOT_FOUND_RESPONSE_CODE, []string{err.Error()}, kafkaRes)
		}
		return common.ConstructResponse[kafka.KafkaSearchResponse](shared.SEARCH_FAILED_RESPONSE_CODE, shared.SEARCH_FAILED_RESPONSE_CODE, []string{err.Error()}, kafkaRes)
	}
	if !containsFlights(res) {
		return common.ConstructResponse[kafka.KafkaSearchResponse](shared.FLIGHT_NOT_FOUND_RESPONSE_CODE, shared.FLIGHT_NOT_FOUND_RESPONSE_CODE, []string{}, kafkaRes)
	}

	baseRes := common.ConstructResponse[kafka.KafkaSearchResponse](shared.SUCCESS_RESPONSE_CODE, shared.SUCCESS_RESPONSE_CODE, nil, kafkaRes)

	return baseRes
}

func containsFlights(res fareRS.FlightIntegratorSearchResponse) bool {
	return len(res.RoundTrips.DepartureFlights) != 0 ||
		len(res.RoundTrips.ReturnFlights) != 0 ||
		len(res.DepartureFlights) != 0 ||
		len(res.ReturnFlights) != 0
}
