package bau

import (
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/bau"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	bookRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/book"
	issuedRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/issued"
)

type noopToken struct{}

func (noopToken) Complete() {}
func (noopToken) CreateSegment(string, bool) *bau.AnalyticSegment {
	return nil
}

type NoopAnalyticsService[T any] struct{}

func (NoopAnalyticsService[T]) CreateToken(mr common.MandatoryRequest, req T) AnalyticsTokenModifier {
	return noopToken{}
}

func (NoopAnalyticsService[T]) Complete(mr common.MandatoryRequest, token AnalyticsTokenModifier, err error) {}

func NewNoopAnalyticsBooking() AnalyticsService[bookRQ.IntegratorBookRequest] {
	return NoopAnalyticsService[bookRQ.IntegratorBookRequest]{}
}

func NewNoopAnalyticsIssued() AnalyticsService[issuedRQ.IntegratorIssuedRequest] {
	return NoopAnalyticsService[issuedRQ.IntegratorIssuedRequest]{}
}
