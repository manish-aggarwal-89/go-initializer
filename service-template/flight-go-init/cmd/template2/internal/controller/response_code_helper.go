package controller

import (
	"errors"
	"{{MODULE_NAME}}/internal/shared"
)

func GetBaggageResponseCodeMapper(err error) string {
	switch {
	case err == nil:
		return shared.SUCCESS_RESPONSE_CODE
	default:
		return shared.BAGGAGE_NOT_AVAILABLE_RESPONSE_CODE
	}
}

// TODO: double check the error mapping in booking
func BookingResponseCodeMapper(err error) string {
	switch {
	case err == nil:
		return shared.SUCCESS_RESPONSE_CODE
	case errors.Is(err, shared.ErrNoSeat):
		return shared.NO_SEAT_RESPONSE_CODE
	case errors.Is(err, shared.ErrJourneyNotFound):
		return shared.NO_SEAT_RESPONSE_CODE
	case errors.Is(err, shared.ErrBookingBaggageNotAvailable):
		return shared.BAGGAGE_NOT_AVAILABLE_RESPONSE_CODE
	case errors.Is(err, shared.ErrBookingSeatNotAvailable):
		return shared.SEAT_NOT_AVAILABLE_RESPONSE_CODE
	case errors.Is(err, shared.ErrBookingMealsNotAvailable):
		return shared.MEALS_NOT_AVAILABLE_RESPONSE_CODE
	case errors.Is(err, shared.ErrBookingSSRIsUnavailable):
		return shared.SSR_IS_UNAVAILABLE_RESPONSE_CODE
	default:
		return shared.BOOKING_FAILED_RESPONSE_CODE
	}
}

// TODO: double check the error mapping in searching
func GetSearchResponseCodeMapper(err error) string {
	switch {
	case err == nil:
		return shared.SUCCESS_RESPONSE_CODE
	case errors.Is(err, shared.ErrInvalidRoute):
		return shared.INVALID_ROUTE_RESPONSE_CODE
	case errors.Is(err, shared.ErrRequestAirlineFailed):
		return shared.SEARCH_FAILED_RESPONSE_CODE
	case errors.Is(err, shared.ErrInvalidCredential):
		return shared.CREDENTIAL_NOT_FOUND_RESPONSE_CODE
	case errors.Is(err, shared.ErrFlightNotFound),
		errors.Is(err, shared.ErrCabinClassIsNotEnabled),
		errors.Is(err, shared.ErrSearchResponseIsEmpty):
		return shared.FLIGHT_NOT_FOUND_RESPONSE_CODE
	default:
		return shared.SEARCH_FAILED_RESPONSE_CODE
	}
}

func GetIssuedResponseCodeMapper(err error) string {
	switch {
	case err == nil:
		return shared.SUCCESS_RESPONSE_CODE
	case errors.Is(err, shared.ErrInvalidCredential):
		return shared.CREDENTIAL_NOT_FOUND_RESPONSE_CODE
	default:
		return shared.ISSUED_FAILED_RESPONSE_CODE
	}
}

func GetCancelBookingResponseCodeMapper(err error) string {
	switch {
	case err == nil:
		return shared.SUCCESS_RESPONSE_CODE
	case errors.Is(err, shared.ErrInvalidCredential):
		return shared.CREDENTIAL_NOT_FOUND_RESPONSE_CODE
	default:
		return shared.CANCEL_BOOKING_FAILED_RESPONSE_CODE
	}
}
