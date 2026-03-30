package common

import (
	"context"

	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/credential"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/constant/enum"
	bookRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/book"
	cancelRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/cancel-book"
	fareRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/fare"
	issuedRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/issued"
	bookRS "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/book"
	cancelRS "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/cancel-book"
	fareRS "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/fare"
	issuedRS "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/issued"
)

type SearchRule interface {
	Search(ctx context.Context, mr common.MandatoryRequest, request fareRQ.IntegratorFareRequest, tripType enum.TripType, cred *credential.CredentialEntity, promoCode string) (*fareRS.FlightIntegratorSearchResponse, error)
}

type BookingRule interface {
	Booking(ctx context.Context, mr common.MandatoryRequest, request bookRQ.IntegratorBookRequest, cred *credential.CredentialEntity, promoCode string) (*bookRS.BookResponse, error)
}

type IssuedRule interface {
	Issued(ctx context.Context, mr common.MandatoryRequest, request issuedRQ.IntegratorIssuedRequest, cred *credential.CredentialEntity) (*issuedRS.IssuedResponse, error)
}

type CancelBookingRule interface {
	CancelBooking(ctx context.Context, mr common.MandatoryRequest, request cancelRQ.IntegratorCancelBookRequest, cred *credential.CredentialEntity) (*cancelRS.CancelBookResponse, error)
}
