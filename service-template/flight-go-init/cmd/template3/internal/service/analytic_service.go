package service

import (
	"context"
	"log/slog"

	"{{MODULE_NAME}}/internal/shared/deps"

	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/bau"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/errormapper"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	bookRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/book"
	fareRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/fare"
	issuedRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/issued"
)

// AnalyticsService is the BAU analytics interface (search, booking, issuance).
type AnalyticsService[T any] interface {
	CreateToken(mr common.MandatoryRequest, integratorRequest T) bau.AnalyticsTokenModifier
	Complete(ctx context.Context, mr common.MandatoryRequest, analyticToken bau.AnalyticsTokenModifier, err error)
}

func NewBauSearchService(deps deps.Deps) AnalyticsService[fareRQ.IntegratorFareRequest] {
	log := slog.Default()
	return bau.NewBAUSearchingService(log, deps.Kafka.KafkaFlight, deps.Config.KafkaConfig.Topics.AnalyticBAUSearch, deps.Config.KafkaConfig.Topics.AnalyticBAULog)
}

func NewBauBookingService(deps deps.Deps, integratorErrorMappingService errormapper.IntegratorErrorMappingService) AnalyticsService[bookRQ.IntegratorBookRequest] {
	log := slog.Default()
	return bau.NewBAUBookingService(log, deps.Kafka.KafkaFlight, deps.Config.KafkaConfig.Topics.AnalyticBAUBooking, deps.Config.KafkaConfig.Topics.AnalyticBAULog, integratorErrorMappingService)
}

func NewBauIssuedService(deps deps.Deps, integratorErrorMappingService errormapper.IntegratorErrorMappingService) AnalyticsService[issuedRQ.IntegratorIssuedRequest] {
	log := slog.Default()
	return bau.NewBAUIssuanceService(log, deps.Kafka.KafkaFlight, deps.Config.KafkaConfig.Topics.AnalyticBAUIssuance, deps.Config.KafkaConfig.Topics.AnalyticBAULog, integratorErrorMappingService)
}
