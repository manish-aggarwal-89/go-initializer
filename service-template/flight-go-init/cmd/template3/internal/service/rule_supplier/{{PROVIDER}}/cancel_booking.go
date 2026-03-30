package {{PROVIDER}}

import (
	"context"
	"errors"

	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/credential"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	cancelRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/cancel-book"
	cancelRS "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/cancel-book"
)

var ErrCancelBookingNotImplemented = errors.New("cancel booking not implemented")

type CancelBookingRuleImpl struct {
	Base
}

func NewCancelBookingRule(b Base) *CancelBookingRuleImpl {
	return &CancelBookingRuleImpl{Base: b}
}

func (r *CancelBookingRuleImpl) CancelBooking(ctx context.Context, mr common.MandatoryRequest, request cancelRQ.IntegratorCancelBookRequest, cred *credential.CredentialEntity) (*cancelRS.CancelBookResponse, error) {
	return nil, ErrCancelBookingNotImplemented
}
