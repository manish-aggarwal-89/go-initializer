package bau

import (
	"errors"
	"regexp"
	"{{MODULE_NAME}}/internal/shared"

	modelBau "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/bau"
)

func GetSearchHealthStatus(err error) SearchHealthStatus {
	switch {
	case err == nil:
		return SearchingHealthStatusSuccess
	case errors.Is(err, shared.ErrFlightNotFound),
		errors.Is(err, shared.ErrSearchResponseIsEmpty),
		errors.Is(err, shared.ErrCabinClassIsNotEnabled):
		return SearchingHealthStatusFlightNotFound
	case errors.Is(err, shared.ErrInvalidRoute):
		return SearchingHealthStatusInvalidRoute
	case errors.Is(err, shared.ErrRequestAirlineFailed),
		errors.Is(err, shared.ErrSearchFailed),
		errors.Is(err, shared.ErrInvalidCredential):
		return SearchingHealthStatusSearchFailed
	default:
		return SearchingHealthStatusUnmapped
	}
}

func GetBookingHealthStatus(err error) BookingHealthStatus {
	switch {
	case err == nil:
		return BookingHealthStatusSuccess

	case errors.Is(err, shared.ErrItineraryNotFound):
		return BookingHealthStatusNoSeat

	case errors.Is(err, shared.ErrRequestAirlineFailed),
		errors.Is(err, shared.ErrBookingMealsNotAvailable),
		errors.Is(err, shared.ErrBookingSeatNotAvailable):
		return BookingHealthStatusBookingFailed

	default:
		return BookingHealthStatusUnmapped
	}
}

func GetIssuanceHealthStatus(err error) IssuanceHealthStatus {
	switch {
	case err == nil:
		return IssuanceHealthStatusSuccess
	case errors.Is(err, shared.ErrRequestAirlineFailed),
		errors.Is(err, shared.ErrIssuedFailed):
		return IssuanceHealthStatusIssuanceFailed
	default:
		return IssuanceHealthStatusUnmapped
	}
}

func SumMap(mapOfProcessingTime map[string]int64) int64 {
	var result int64
	for _, processingTime := range mapOfProcessingTime {
		result += processingTime
	}

	return result
}

func RemoveNumberSuffix(s string) string {
	reg := regexp.MustCompile("_[0-9]+")
	return reg.ReplaceAllString(s, "")
}

func ConstructAirlineProcessTimeDetail(analyticSegments []*modelBau.AnalyticSegment) map[string]int64 {
	result := make(map[string]int64)

	for _, analyticSegment := range analyticSegments {
		processTime := analyticSegment.EndMillis - analyticSegment.StartMillis

		if !analyticSegment.IsParallel {
			result[analyticSegment.Name] = processTime
			continue
		}

		name := RemoveNumberSuffix(analyticSegment.Name)
		if currentProcessTime, ok := result[name]; ok {
			if processTime > currentProcessTime {
				result[name] = processTime
			}
		} else {
			result[name] = processTime
		}
	}

	return result
}
