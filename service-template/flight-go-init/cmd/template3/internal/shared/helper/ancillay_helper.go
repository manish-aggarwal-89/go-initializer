package helper

import (
	bookRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/book"
	bookRS "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/book"
)

var SSR_ITINERARY_DEFAULT = bookRS.SsrItinerary{
	SsrPaxBaggages: []bookRS.SsrPaxBaggages{},
	SsrPaxMeals:    []bookRS.SsrPaxMeals{},
	SsrPaxSeats:    []bookRS.SsrPaxSeats{},
	SsrPaxMedicals: []bookRS.SsrPaxMedicals{},
}

var SSR_ITINERARY_LIST_OW_DEFAULT = []bookRS.SsrItinerary{
	SSR_ITINERARY_DEFAULT,
}

var SSR_ITINERARY_LIST_RT_DEFAULT = []bookRS.SsrItinerary{
	SSR_ITINERARY_DEFAULT, SSR_ITINERARY_DEFAULT,
}

// IsAncillaryAdded checks whether there is any ancillary
func IsAncillaryAdded(request bookRQ.IntegratorBookRequest) bool {
	return false // todo - to implement once we have addons , for now this is required to return default ssrFee
}

func constructSsrItineraries(request bookRQ.IntegratorBookRequest) []bookRS.SsrItinerary {
	if !IsAncillaryAdded(request) {
		if len(request.Itineraries) > 1 {
			return SSR_ITINERARY_LIST_RT_DEFAULT
		}
		return SSR_ITINERARY_LIST_OW_DEFAULT
	}
	return make([]bookRS.SsrItinerary, len(request.Itineraries)) // todo - logic to implement once we have addons
}
