package {{PROVIDER}}

import (
	"context"
	"errors"

	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/credential"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	bookRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/book"
	bookRS "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/book"
)

var ErrBookingNotImplemented = errors.New("booking not implemented")

type BookingRuleImpl struct {
	Base
}

func NewBookingRule(b Base) *BookingRuleImpl {
	return &BookingRuleImpl{Base: b}
}

func (r *BookingRuleImpl) Booking(ctx context.Context, mr common.MandatoryRequest, request bookRQ.IntegratorBookRequest, cred *credential.CredentialEntity, promoCode string) (*bookRS.BookResponse, error) {
	return nil, ErrBookingNotImplemented
}
