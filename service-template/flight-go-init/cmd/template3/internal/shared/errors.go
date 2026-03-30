package shared

import (
	"errors"
	"net"
	"os"
	"strings"
)

var (
	ErrBaggageResponseIsEmpty  = errors.New("baggage response is empty")
	ErrBaggageAncillaryIsEmpty = errors.New("baggage response for ancillary is empty")
	ErrItineraryLengthNotMatch = errors.New("itinerary length is not match")
	ErrItineraryNotFound       = errors.New("itinerary not found")
	ErrTripTypeNotSupported    = errors.New("trip type not supported")

	ErrBookingSSRIsUnavailable    = errors.New("booking ssr is unavailable")
	ErrBookingMealsNotAvailable   = errors.New("booking meals not available")
	ErrBookingSeatNotAvailable    = errors.New("booking seat not available")
	ErrBookingBaggageNotAvailable = errors.New("booking baggage not available")

	ErrAddonsEmptyFlightNumbers = errors.New("addons empty flight numbers")

	ErrSellSsr                                     = errors.New("error sell ssr")
	ErrSSRCodeNotFound                             = errors.New("SSR code not found")
	ErrInvalidCredential                           = errors.New("invalid credential")
	ErrRequestAirlineFailed                        = errors.New("request airline failed")
	ErrFlightNotFound                              = errors.New("flight not found")
	ErrGetFareFailed                               = errors.New("get fare failed")
	ErrInvalidRoute                                = errors.New("invalid route")
	ErrNoRoute                                     = errors.New("no route")
	ErrNoSeat                                      = errors.New("no seat")
	ErrSearchFailed                                = errors.New("search failed")
	ErrBookingFailed                               = errors.New("booking failed")
	ErrGetBaggageFailed                            = errors.New("get baggage failed")
	ErrBaggageNotAvailable                         = errors.New("baggage not available")
	ErrCancelBookingFailed                         = errors.New("cancel booking failed")
	ErrIssuedFailed                                = errors.New("issued Failed")
	ErrDataIsExist                                 = errors.New("data is exist")
	ErrDataNotExist                                = errors.New("data not exist")
	ErrDuplicateData                               = errors.New("duplicate data")
	ErrCountDocumentError                          = errors.New("count document error")
	ErrOriginalCurrencyNotFound                    = errors.New("original currency not found")
	ErrInvalidBaggageUom                           = errors.New("invalid baggage uom")
	ErrPnrExpired                                  = errors.New("pnr already expired")
	ErrSearchResponseIsEmpty                       = errors.New("search RS empty")
	ErrInvalidObjectID                             = errors.New("invalid object id")
	ErrInvalidTypeConversion                       = errors.New("invalid type conversion")
	ErrFareNotFound                                = errors.New("fare not found. unable to find correct fare")
	ErrProductClassNotFound                        = errors.New("product class not found")
	ErrBookingQuoteFailed                          = errors.New("all booking quote failed")
	ErrBookingQuoteNotFound                        = errors.New("booking quote not found")
	ErrCabinClassIsNotEnabled                      = errors.New("cabin class is not enabled")
	ErrMalformedBaggageWeight                      = errors.New("malformed baggage weight")
	ErrMalformedBaggageUom                         = errors.New("malformed baggage unit of measurement")
	ErrGetBaggageNotSupportedMultiPCC              = errors.New("non hold baggage is not supported for multi PCC")
	ErrGetSeatFailed                               = errors.New("get seat failed")
	ErrSeatsNotAvailable                           = errors.New("seats not available")
	ErrIntegratorBookRequestAdditionalDataNotFound = errors.New("IntegratorBookRequest not found in additionalData")
	ErrRevalidateItineraryFailed                   = errors.New("revalidate itinerary failed")

	ErrJourneyNotFound = errors.New("journey not found ")

	ErrFaultyDateFormat = errors.New("faulty date format")

	ErrMealsNotAvailable = errors.New("meals not available")
)

func CheckTimeoutError(err error) bool {
	if err == nil {
		return false
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}

	if os.IsTimeout(err) {
		return true
	}

	// to check if the error is wrapped, and we can't
	if strings.Contains(err.Error(), "context deadline exceeded") {
		return true
	}

	return false
}
