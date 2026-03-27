package service

import (
	"context"
	"{{MODULE_NAME}}/internal/outbound"
	"{{MODULE_NAME}}/internal/shared/deps"

	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	cancelRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/cancel-book"
	cancelRS "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/cancel-book"
)

type (
	CancelBookingServiceItf interface {
		CancelBooking(ctx context.Context, mr common.MandatoryRequest, req cancelRQ.IntegratorCancelBookRequest) (cancelRS.CancelBookResponse, error)
	}

	CancelBookingServiceImpl struct {
		deps     deps.Deps
		outbound outbound.OutboundItf
	}
)

func NewCancelBookingServiceImpl(deps deps.Deps,
	outbound outbound.OutboundItf,
) *CancelBookingServiceImpl {
	return &CancelBookingServiceImpl{
		deps:     deps,
		outbound: outbound,
	}
}

func (impl *CancelBookingServiceImpl) CancelBooking(ctx context.Context, mr common.MandatoryRequest, req cancelRQ.IntegratorCancelBookRequest) (cancelRS.CancelBookResponse, error) {
	var (
		//op            = "cancel-booking"
		//supplierCode  = req.Account.Code
		cancelRes = cancelRS.CancelBookResponse{
			MandatoryRequest:            mr,
			IntegratorCancelBookRequest: req,
		}
	)

	//Implement Cancel-Booking Logic

	return cancelRes, nil
}
