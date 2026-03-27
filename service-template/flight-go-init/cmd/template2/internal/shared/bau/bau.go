package bau

import (
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/bau"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
)

type AnalyticsTokenModifier interface {
	Complete()
	CreateSegment(string, bool) *bau.AnalyticSegment
}

type AnalyticsService[T any] interface {
	CreateToken(mr common.MandatoryRequest, integratorRequest T) AnalyticsTokenModifier
	Complete(mr common.MandatoryRequest, analyticToken AnalyticsTokenModifier, err error)
}
