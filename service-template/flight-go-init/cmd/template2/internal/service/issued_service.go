package service

import (
	"context"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/credential"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	issuedReq "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/issued"
	issuedRes "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/issued"
	"{{MODULE_NAME}}/internal/outbound"
	"{{MODULE_NAME}}/internal/shared/deps"
)

type IssuedServiceImpl struct {
	deps          deps.Deps
	outbound      outbound.OutboundItf
	credentialSvc credential.Service
}

type IssuedServiceItf interface {
	Issued(ctx context.Context, mr common.MandatoryRequest, req issuedReq.IntegratorIssuedRequest) (issuedRes.IssuedResponse, error)
}

func NewIssuedServiceImpl(deps deps.Deps, outbound outbound.OutboundItf, credentialSvc credential.Service) *IssuedServiceImpl {
	return &IssuedServiceImpl{
		deps:          deps,
		outbound:      outbound,
		credentialSvc: credentialSvc,
	}
}

func (impl *IssuedServiceImpl) Issued(ctx context.Context, mr common.MandatoryRequest, req issuedReq.IntegratorIssuedRequest) (issuedRes.IssuedResponse, error) {
	return issuedRes.IssuedResponse{}, nil
}
