package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/bau"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/util"

	commonbau "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/bau"
	common "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/issued"
)

type BAUIssuedService struct {
	deps deps.Deps

	kafkaPublisherService KafkaPublisherServiceItf
}

func NewBAUIssuedService(deps deps.Deps, kafkaPublisherService KafkaPublisherServiceItf) bau.AnalyticsService[issued.IntegratorIssuedRequest] {
	return &BAUIssuedService{
		deps:                  deps,
		kafkaPublisherService: kafkaPublisherService,
	}
}

func (svc *BAUIssuedService) CreateToken(mandatoryRequest common.MandatoryRequest, integratorIssuedRequest issued.IntegratorIssuedRequest) bau.AnalyticsTokenModifier {
	return &bau.AnalyticsIssuanceToken{
		AnalyticsToken: commonbau.AnalyticsToken{
			Name:                    bau.TransactionTypeIssuance,
			StoreId:                 mandatoryRequest.StoreId,
			ChannelId:               mandatoryRequest.ChannelId,
			RequestId:               mandatoryRequest.RequestId,
			Distribution:            integratorIssuedRequest.DistributionType,
			Airlines:                getAirlines(integratorIssuedRequest),
			Supplier:                integratorIssuedRequest.Account.Code,
			StartMillis:             time.Now().UnixMilli(),
			MutableAnalyticsSegment: make([]*commonbau.AnalyticSegment, 0),
			Counter:                 sync.Map{},
		},
		BookingCode: integratorIssuedRequest.BookingCode,
	}
}

func getAirlines(req issued.IntegratorIssuedRequest) (airlines []string) {
	itineraries := strings.Split(req.FlightSelect, shared.ITINERARY_SEPARATOR)

	for _, itin := range itineraries {
		transits := strings.Split(itin, shared.TRANSIT_SEPARATOR)
		for _, trans := range transits {
			if len(trans) > 1 {
				airlines = append(airlines, trans[0:2])
			}
		}
	}
	return
}

func (service *BAUIssuedService) Complete(mr common.MandatoryRequest, token bau.AnalyticsTokenModifier, businessErr error) {
	token.Complete()
	issuanceToken := token.(*bau.AnalyticsIssuanceToken)

	status := bau.GetIssuanceHealthStatus(businessErr)
	analyticsReq := constructIssuanceAnalyticsReq(issuanceToken, status)
	ctx := context.Background()
	jsonBytes, err := json.Marshal(&analyticsReq)
	if err != nil {
		service.deps.GetLogger(ctx).Errorf("{%s} unable to marshal analytics request", issuanceToken.RequestId)
	}
	jsonStr := string(jsonBytes)

	go func() {
		if err := service.kafkaPublisherService.Publish(mr, service.deps.Config.KafkaConfig.Topics.AnalyticBAUIssuance, jsonStr, false); err != nil {
			service.deps.GetLogger(ctx).Errorf("%s %s - error publish to kafka, with err: %s", "[BAU_Issued]", issuanceToken.RequestId, err.Error())
		}
	}()

	analyticLogRecord := constructAnalyticIssuedLogRecord(analyticsReq, businessErr)
	bAnalyticLog, err := json.Marshal(analyticLogRecord)
	if err != nil {
		service.deps.GetLogger(nil).Errorf("%s %s - %s", "[BAU_Log_Issuance]", issuanceToken.RequestId, err.Error())
		return
	}

	go func() {
		err := service.kafkaPublisherService.Publish(mr, service.deps.Config.KafkaConfig.Topics.AnalyticBAULog, string(bAnalyticLog), false)
		if err != nil {
			service.deps.GetLogger(nil).Errorf("%s %s - error publish to kafka, with err: %s", "[BAU_Log_Issued]", issuanceToken.RequestId, err.Error())
		}
	}()
}

func constructIssuanceAnalyticsReq(token *bau.AnalyticsIssuanceToken, status bau.IssuanceHealthStatus) bau.AnalyticsIssuanceRequest {
	airlineProcessTimes := bau.ConstructAirlineProcessTimeDetail(token.MutableAnalyticsSegment)
	totalProcessTime := bau.SumMap(airlineProcessTimes)

	transactionDate := time.UnixMilli(token.StartMillis).Format(util.DateTimeFormatWithSecond)

	return bau.AnalyticsIssuanceRequest{
		RequestId:                token.RequestId,
		Distribution:             token.Distribution,
		Supplier:                 token.Supplier,
		Airlines:                 token.Airlines,
		TransactionDate:          transactionDate,
		AirlineProcessTime:       totalProcessTime,
		AirlineProcessTimeDetail: airlineProcessTimes,
		Status:                   status,
		BookingCode:              token.BookingCode,
		StoreId:                  token.StoreId,
		ChannelId:                token.ChannelId,
	}
}

func constructAnalyticIssuedLogRecord(analyticIssuedRequest bau.AnalyticsIssuanceRequest, err error) commonbau.AnalyticLogRecord {
	code, message := determineAnalyticIssuedLogCodeAndMessage(err)

	return commonbau.AnalyticLogRecord{
		Code:     code.String(),
		Message:  message,
		Event:    strings.ToUpper(bau.TransactionTypeIssuance),
		Supplier: analyticIssuedRequest.Supplier,
		Airline:  analyticIssuedRequest.Airlines,
	}
}

// code is from analytic code severity
// message is from common-model.ResponseCode.java -- messageId (you can check in the integrator java code)
func determineAnalyticIssuedLogCodeAndMessage(err error) (bau.IssuanceLogHealthStatus, string) {
	switch {
	case err == nil:
		return bau.IssuanceLogHealthStatusSuccess, "Berhasil"
	case errors.Is(err, shared.ErrIssuedFailed):
		return bau.IssuanceLogHealthStatusIssuanceFailed, "Issued failed"
	default:
		return bau.IssuanceLogHealthStatusUnknownError, "Error Tidak Diketahui"
	}
}
