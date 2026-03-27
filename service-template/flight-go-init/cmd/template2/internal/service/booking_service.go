package service

import (
	"context"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/credential"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/system_param"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	bookRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/book"
	bookRS "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/book"
	"strconv"
	"{{MODULE_NAME}}/internal/outbound"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"
)

type BookingServiceItf interface {
	Booking(ctx context.Context, mr common.MandatoryRequest, integratorBookRequest bookRQ.IntegratorBookRequest) (bookRS.BookResponse, error)
}

type BookingServiceImpl struct {
	deps               deps.Deps
	outbound           outbound.OutboundItf
	systemParameterSvc system_param.Service
	credentialService  credential.Service
	searchSvc          SearchServiceItf
}

func NewBookingServiceImpl(deps deps.Deps,
	outbound outbound.OutboundItf,
	systemParameterSvc system_param.Service,
	searchSvc SearchServiceItf,
	credentialService credential.Service,
) *BookingServiceImpl {
	return &BookingServiceImpl{
		deps:               deps,
		outbound:           outbound,
		systemParameterSvc: systemParameterSvc,
		credentialService:  credentialService,
		searchSvc:          searchSvc,
	}
}

func (impl *BookingServiceImpl) getPublishOnBookingFlag(ctx context.Context, mr common.MandatoryRequest) bool {
	systemParameter, err := impl.systemParameterSvc.FindByVariableInCache(ctx, mr, shared.SYS_PARAM_PUBLISH_ON_BOOKING_FLAG_KEY)
	if err != nil {
		return true
	}

	result, err := strconv.ParseBool(systemParameter.Value)
	if err != nil {
		return true
	}

	return result
}

func (impl *BookingServiceImpl) Booking(ctx context.Context, mr common.MandatoryRequest, integratorBookRequest bookRQ.IntegratorBookRequest) (bookRS.BookResponse, error) {
	return bookRS.BookResponse{}, nil
}
