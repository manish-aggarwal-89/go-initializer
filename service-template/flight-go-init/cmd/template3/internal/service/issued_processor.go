package service

import (
	"context"
	"fmt"

	"{{MODULE_NAME}}/config"
	"{{MODULE_NAME}}/internal/service/rule_supplier"
	"{{MODULE_NAME}}/internal/shared/deps"

	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/credential"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	issuedRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/issued"
	issuedRS "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/issued"
)

type IssuedProcessorService interface {
	IssuedProcessor(ctx context.Context, mr common.MandatoryRequest, req issuedRQ.IntegratorIssuedRequest) (*issuedRS.IssuedResponse, error)
}

type issuedProcessorImpl struct {
	deps              deps.Deps
	config            *config.Config
	credentialService credential.Service
	issuedRegistry    *rule_supplier.IssuedRuleRegistry
}

func NewIssuedProcessorService(
	deps deps.Deps,
	cfg *config.Config,
	credentialService credential.Service,
	issuedRegistry *rule_supplier.IssuedRuleRegistry,
) (IssuedProcessorService, error) {
	return &issuedProcessorImpl{
		deps:              deps,
		config:            cfg,
		credentialService: credentialService,
		issuedRegistry:    issuedRegistry,
	}, nil
}

func (p *issuedProcessorImpl) IssuedProcessor(ctx context.Context, mr common.MandatoryRequest, req issuedRQ.IntegratorIssuedRequest) (*issuedRS.IssuedResponse, error) {
	distType := p.config.DistributionType
	if distType == "" {
		return nil, fmt.Errorf("config.DistributionType is required but not set")
	}
	rule, ok := p.issuedRegistry.Get(distType)
	if !ok {
		return nil, fmt.Errorf("issued rule not found for distribution type: %s", distType)
	}
	supplierCode := req.Account.Code
	if supplierCode == "" {
		return nil, fmt.Errorf("supplier code is empty")
	}
	cred, err := p.credentialService.FindBySupplierInCache(ctx, mr, supplierCode)
	if err != nil || cred == nil {
		return nil, fmt.Errorf("credential not found for supplier %s: %w", supplierCode, err)
	}
	return rule.Issued(ctx, mr, req, cred)
}
