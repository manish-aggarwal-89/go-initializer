package {{PROVIDER}}

import (
	"context"
	"errors"

	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/credential"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/constant/enum"
	fareRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/fare"
	fareRS "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/fare"
)

var ErrSearchNotImplemented = errors.New("search not implemented")

type SearchRuleImpl struct {
	Base
}

func NewSearchRule(b Base) *SearchRuleImpl {
	return &SearchRuleImpl{Base: b}
}

func (r *SearchRuleImpl) Search(ctx context.Context, mr common.MandatoryRequest, request fareRQ.IntegratorFareRequest, tripType enum.TripType, cred *credential.CredentialEntity, promoCode string) (*fareRS.FlightIntegratorSearchResponse, error) {
	return nil, ErrSearchNotImplemented
}
