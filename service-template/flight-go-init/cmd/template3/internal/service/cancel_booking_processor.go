package service

import (
	"context"
	"fmt"

	"{{MODULE_NAME}}/config"
	"{{MODULE_NAME}}/internal/service/rule_supplier"
	"{{MODULE_NAME}}/internal/shared/deps"

	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/credential"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	cancelRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/cancel-book"
	cancelRS "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/cancel-book"
)

type CancelBookingProcessorService interface {
	CancelBookingProcessor(ctx context.Context, mr common.MandatoryRequest, req cancelRQ.IntegratorCancelBookRequest) (*cancelRS.CancelBookResponse, error)
}

type cancelBookingProcessorImpl struct {
	deps                  deps.Deps
	config                *config.Config
	credentialService     credential.Service
	cancelBookingRegistry *rule_supplier.CancelBookingRuleRegistry
}

func NewCancelBookingProcessorService(
	deps deps.Deps,
	cfg *config.Config,
	credentialService credential.Service,
	cancelBookingRegistry *rule_supplier.CancelBookingRuleRegistry,
) (CancelBookingProcessorService, error) {
	return &cancelBookingProcessorImpl{
		deps:                  deps,
		config:                cfg,
		credentialService:     credentialService,
		cancelBookingRegistry: cancelBookingRegistry,
	}, nil
}

func (p *cancelBookingProcessorImpl) CancelBookingProcessor(ctx context.Context, mr common.MandatoryRequest, req cancelRQ.IntegratorCancelBookRequest) (*cancelRS.CancelBookResponse, error) {
	distType := p.config.DistributionType
	if distType == "" {
		return nil, fmt.Errorf("config.DistributionType is required but not set")
	}
	rule, ok := p.cancelBookingRegistry.Get(distType)
	if !ok {
		return nil, fmt.Errorf("cancel booking rule not found for distribution type: %s", distType)
	}
	supplierCode := req.Account.Code
	if supplierCode == "" {
		return nil, fmt.Errorf("supplier code is empty")
	}
	cred, err := p.credentialService.FindBySupplierInCache(ctx, mr, supplierCode)
	if err != nil || cred == nil {
		return nil, fmt.Errorf("credential not found for supplier %s: %w", supplierCode, err)
	}
	return rule.CancelBooking(ctx, mr, req, cred)
}
