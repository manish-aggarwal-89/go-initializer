package service

import (
	"context"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	bookRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/book"
	bookRS "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/book"
	"strconv"
	"{{MODULE_NAME}}/internal/outbound"
	"{{MODULE_NAME}}/internal/predefine"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"
)

type BookingServiceItf interface {
	Booking(ctx context.Context, mr common.MandatoryRequest, integratorBookRequest bookRQ.IntegratorBookRequest) (bookRS.BookResponse, error)
}

type BookingServiceImpl struct {
	deps               deps.Deps
	outbound           outbound.OutboundItf
	sessionSvc         SessionService
	systemParameterSvc SystemParameterServiceItf
	credentialService  CredentialServiceItf
	predefined         predefine.Predefined
	searchSvc          SearchServiceItf
}

func NewBookingServiceImpl(deps deps.Deps,
	outbound outbound.OutboundItf,
	sessionSvc SessionService,
	systemParameterSvc SystemParameterServiceItf,
	predefined predefine.Predefined,
	searchSvc SearchServiceItf,
	credentialService CredentialServiceItf,
) *BookingServiceImpl {
	return &BookingServiceImpl{
		deps:               deps,
		sessionSvc:         sessionSvc,
		outbound:           outbound,
		systemParameterSvc: systemParameterSvc,
		credentialService:  credentialService,
		predefined:         predefined,
		searchSvc:          searchSvc,
	}
}

func (impl *BookingServiceImpl) getPublishOnBookingFlag(ctx context.Context, mr common.MandatoryRequest) bool {
	systemParameter, err := impl.systemParameterSvc.FindByVariableByCache(ctx, mr, shared.SYS_PARAM_PUBLISH_ON_BOOKING_FLAG_KEY)
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
