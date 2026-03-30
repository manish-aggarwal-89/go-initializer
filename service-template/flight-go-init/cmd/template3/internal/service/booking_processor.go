package service

import (
	"context"
	"fmt"
	"strconv"

	"{{MODULE_NAME}}/config"
	"{{MODULE_NAME}}/internal/service/rule_supplier"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"

	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/credential"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/system_param"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	bookRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/book"
	bookRS "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/book"
)

type BookingProcessorService interface {
	BookingProcessor(ctx context.Context, mr common.MandatoryRequest, req bookRQ.IntegratorBookRequest) (*bookRS.BookResponse, error)
}

type bookingProcessorImpl struct {
	deps               deps.Deps
	config             *config.Config
	credentialService  credential.Service
	systemParameterSvc system_param.Service
	bookingRegistry    *rule_supplier.BookingRuleRegistry
}

func NewBookingProcessorService(
	deps deps.Deps,
	cfg *config.Config,
	credentialService credential.Service,
	systemParameterSvc system_param.Service,
	bookingRegistry *rule_supplier.BookingRuleRegistry,
) (BookingProcessorService, error) {
	return &bookingProcessorImpl{
		deps:               deps,
		config:             cfg,
		credentialService:  credentialService,
		systemParameterSvc: systemParameterSvc,
		bookingRegistry:    bookingRegistry,
	}, nil
}

func (p *bookingProcessorImpl) getPublishOnBookingFlag(ctx context.Context, mr common.MandatoryRequest) bool {
	systemParameter, err := p.systemParameterSvc.FindByVariableInCache(ctx, mr, shared.SYS_PARAM_PUBLISH_ON_BOOKING_FLAG_KEY)
	if err != nil {
		return true
	}

	result, err := strconv.ParseBool(systemParameter.Value)
	if err != nil {
		return true
	}

	return result
}

func (p *bookingProcessorImpl) BookingProcessor(ctx context.Context, mr common.MandatoryRequest, req bookRQ.IntegratorBookRequest) (*bookRS.BookResponse, error) {
	distType := p.config.DistributionType
	if distType == "" {
		return nil, fmt.Errorf("config.DistributionType is required but not set")
	}
	rule, ok := p.bookingRegistry.Get(distType)
	if !ok {
		return nil, fmt.Errorf("booking rule not found for distribution type: %s", distType)
	}
	supplierCode := req.Account.Code
	if supplierCode == "" {
		return nil, fmt.Errorf("supplier code is empty")
	}
	cred, err := p.credentialService.FindBySupplierInCache(ctx, mr, supplierCode)
	if err != nil || cred == nil {
		return nil, fmt.Errorf("credential not found for supplier %s: %w", supplierCode, err)
	}
	return rule.Booking(ctx, mr, req, cred, "")
}
