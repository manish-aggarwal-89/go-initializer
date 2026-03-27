package helper

import (
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/constant/enum"
	bookRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/book"
	fareRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/fare"
	"strings"
)

func FormatDateToYYMMDD(date string) string {
	return strings.ReplaceAll(date, "-", "")
}

func BuildAirlines(integratorBookRequest bookRQ.IntegratorBookRequest) []string {
	airlinesMap := make(map[string]bool)

	for _, itinerary := range integratorBookRequest.Itineraries {
		airlinesMap[itinerary.Airline] = true
	}

	result := make([]string, 0, len(airlinesMap))

	for airline := range airlinesMap {
		result = append(result, airline)
	}

	return result
}

func DefineTripType(integratorBookRequest bookRQ.IntegratorBookRequest) enum.TripType {
	if len(integratorBookRequest.Itineraries) > 1 {
		return enum.ROUND_TRIP
	}

	return enum.DEPARTURE
}

func buildRoutes(integratorBookRequest bookRQ.IntegratorBookRequest) []fareRQ.Route {
	routes := make([]fareRQ.Route, 0)

	for _, itinerary := range integratorBookRequest.Itineraries {
		routes = append(routes, fareRQ.Route{
			DepartureDate: itinerary.Date,
			Origin:        itinerary.Departure,
			Destination:   itinerary.Arrival,
			Currency:      "IDR",
		})
	}

	return routes
}

func buildReturnDate(integratorBookRequest bookRQ.IntegratorBookRequest) string {
	itinLength := len(integratorBookRequest.Itineraries)

	if itinLength == 1 {
		return ""
	}

	return integratorBookRequest.Itineraries[itinLength-1].Date
}

func collectFlightNumbersFromBooking(integratorBookRequest bookRQ.IntegratorBookRequest) []string {
	itineraryList := integratorBookRequest.Itineraries
	itinerarySize := len(itineraryList)
	journeyFlightSelects := []string{}
	trimmedFlightSelect := strings.ReplaceAll(strings.TrimSpace(integratorBookRequest.FlightSelect), " ", "")

	if itinerarySize == 0 {
		journeyFlightSelects := strings.Split(trimmedFlightSelect, "~")
		return journeyFlightSelects
	}

	for _, itinerary := range itineraryList {
		flightNumbers := []string{}
		for _, schedule := range itinerary.Schedules {
			flightNumber := schedule.OperatingAirline.Code + schedule.OperatingAirline.Number
			flightNumbers = append(flightNumbers, flightNumber)
		}
		journeyFlightSelects = append(journeyFlightSelects, strings.Join(flightNumbers, "|"))
	}
	return journeyFlightSelects
}

func ToIntegratorFareRequest(req bookRQ.IntegratorBookRequest) (fareRQ.IntegratorFareRequest, enum.TripType) {
	tripType := DefineTripType(req)
	isInternationalFlight := isInternationalFlight(req.Itineraries)
	return fareRQ.IntegratorFareRequest{
		DistributionType: req.DistributionType,
		SupplierRequest:  buildSupplierRequest(req),
		SupplierMappings: buildSupplierMappings(req),
		Adult:            req.Adult,
		Child:            req.Child,
		Infant:           req.Infant,
		CabinClass:       req.Itineraries[0].CabinClass,
		Currency: fareRQ.Currency{
			DepartureCurrency: req.Currency,
			ReturnCurrency:    req.Currency,
		},
		DepartureDate: req.Itineraries[0].Date,
		ReturnDate:    buildReturnDate(req),
		Origin:        req.Itineraries[0].Departure,
		Destination:   req.Itineraries[0].Arrival,
		TripTypes:     []enum.TripType{tripType},
		International: isInternationalFlight,
		BlacklistMap:  buildBlacklistMap(req),
		Routes:        buildRoutes(req),
	}, tripType
}

func isInternationalFlight(itineraries []bookRQ.Itinerary) int {
	for _, itinerary := range itineraries {
		if itinerary.International == 1 {
			return 1
		}
	}

	return 0
}

func buildSupplierRequest(integratorBookRequest bookRQ.IntegratorBookRequest) fareRQ.SupplierRequest {
	return fareRQ.SupplierRequest{
		Accounts: []fareRQ.Account{integratorBookRequest.Account},
		Airlines: BuildAirlines(integratorBookRequest),
	}
}

func buildSupplierMappings(integratorBookRequest bookRQ.IntegratorBookRequest) []fareRQ.SupplierMapping {
	return []fareRQ.SupplierMapping{
		{
			Account:  integratorBookRequest.Account,
			Airlines: BuildAirlines(integratorBookRequest),
		},
	}
}

func buildBlacklistMap(integratorBookRequst bookRQ.IntegratorBookRequest) map[string][]string {
	blacklistMap := make(map[string][]string)

	for _, itinerary := range integratorBookRequst.Itineraries {
		for k, v := range itinerary.BlacklistMap {
			_, ok := blacklistMap[k]
			if ok {
				continue
			}

			blacklistMap[k] = v
		}
	}

	return blacklistMap
}
