package {{PROVIDER}}

import (
	"context"
	"errors"

	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/credential"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	issuedRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/issued"
	issuedRS "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/issued"
)

var ErrIssuedNotImplemented = errors.New("issued not implemented")

type IssuedRuleImpl struct {
	Base
}

func NewIssuedRule(b Base) *IssuedRuleImpl {
	return &IssuedRuleImpl{Base: b}
}

func (r *IssuedRuleImpl) Issued(ctx context.Context, mr common.MandatoryRequest, request issuedRQ.IntegratorIssuedRequest, cred *credential.CredentialEntity) (*issuedRS.IssuedResponse, error) {
	return nil, ErrIssuedNotImplemented
}
