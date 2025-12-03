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

	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/log/logrus"
	commonbau "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/bau"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/fare"
)

type BAUSearchService struct {
	deps                  deps.Deps
	kafkaPublisherService KafkaPublisherServiceItf
}

func NewBAUSearchService(deps deps.Deps, kafkaPublisherService KafkaPublisherServiceItf) bau.AnalyticsService[fare.IntegratorFareRequest] {
	return &BAUSearchService{
		deps:                  deps,
		kafkaPublisherService: kafkaPublisherService,
	}
}

func (svc *BAUSearchService) CreateToken(mandatoryRequest common.MandatoryRequest, integratorFareRequest fare.IntegratorFareRequest) bau.AnalyticsTokenModifier {
	return &bau.AnalyticSearchToken{
		AnalyticsToken: commonbau.AnalyticsToken{
			Name:                    bau.TransactionTypeSearch,
			StoreId:                 mandatoryRequest.StoreId,
			ChannelId:               mandatoryRequest.ChannelId,
			RequestId:               mandatoryRequest.RequestId,
			Distribution:            integratorFareRequest.DistributionType,
			Airlines:                integratorFareRequest.SupplierRequest.Airlines,
			Supplier:                integratorFareRequest.SupplierRequest.Accounts[0].Code,
			StartMillis:             time.Now().UnixMilli(),
			MutableAnalyticsSegment: make([]*commonbau.AnalyticSegment, 0),
			Counter:                 sync.Map{},
		},
		Origin:      integratorFareRequest.Origin,
		Destination: integratorFareRequest.Destination,
		TripType:    string(integratorFareRequest.TripTypes[0]),
	}
}

func (svc *BAUSearchService) Complete(mr common.MandatoryRequest, token bau.AnalyticsTokenModifier, businessErr error) {
	token.Complete()
	searchToken := token.(*bau.AnalyticSearchToken)

	healthStatus := bau.GetSearchHealthStatus(businessErr)
	req := constructAnalyticSearchRequest(*searchToken, healthStatus)
	ctx := logrus.InjectRequestMetadataToContext(context.Background(), searchToken.RequestId, searchToken.Name, "")

	b, err := json.Marshal(req)
	if err != nil {
		svc.deps.GetLogger(ctx).Errorf("%s %s - %s", "[BAU_Search]", searchToken.RequestId, err.Error())
		return
	}

	go func() {
		err := svc.kafkaPublisherService.Publish(mr, svc.deps.Config.KafkaConfig.Topics.AnalyticBAUSearch, string(b), false)
		if err != nil {
			svc.deps.GetLogger(ctx).Errorf("%s %s - error publish to kafka, with err: %s", "[BAU_Search]", searchToken.RequestId, err.Error())
		}
	}()

	analyticSearchLogRecord := constructAnalyticSearchLogRecord(req, businessErr)
	bAnalyticLog, err := json.Marshal(analyticSearchLogRecord)
	if err != nil {
		svc.deps.GetLogger(ctx).Errorf("%s %s - %s", "[BAU_Log_Search]", searchToken.RequestId, err.Error())
		return
	}

	go func() {
		err := svc.kafkaPublisherService.Publish(mr, svc.deps.Config.KafkaConfig.Topics.AnalyticBAULog, string(bAnalyticLog), false)
		if err != nil {
			svc.deps.GetLogger(ctx).Errorf("%s %s - error publish to kafka, with err: %s", "[BAU_Log_Search]", searchToken.RequestId, err.Error())
		}
	}()
}

func constructAnalyticSearchRequest(token bau.AnalyticSearchToken, healthStatus bau.SearchHealthStatus) bau.AnalyticSearchRequest {
	airlineProcessTimeDetail := bau.ConstructAirlineProcessTimeDetail(token.AnalyticsToken.MutableAnalyticsSegment)
	airlineProcessTime := bau.SumMap(airlineProcessTimeDetail)

	transactionDate := time.UnixMilli(token.StartMillis).Format(util.DateTimeFormatWithSecond)

	return bau.AnalyticSearchRequest{
		RequestID:                token.AnalyticsToken.RequestId,
		Distribution:             token.Distribution,
		Supplier:                 token.Supplier,
		Airlines:                 token.Airlines,
		TransactionDate:          transactionDate,
		AirlineProcessTimeDetail: airlineProcessTimeDetail,
		AirlineProcessTime:       airlineProcessTime,
		ProcessTime:              token.EndMillis - token.StartMillis,
		Origin:                   token.Origin,
		Destination:              token.Destination,
		TripType:                 token.TripType,
		SearchHealthStatus:       healthStatus,
		StoreID:                  token.AnalyticsToken.StoreId,
		ChannelID:                token.AnalyticsToken.ChannelId,
	}
}

func constructAnalyticSearchLogRecord(analyticSearchRequest bau.AnalyticSearchRequest, err error) commonbau.AnalyticLogRecord {
	code, message := determineAnalyticLogCodeAndMessage(err)

	return commonbau.AnalyticLogRecord{
		Code:     code.String(),
		Message:  message,
		Event:    strings.ToUpper(bau.TransactionTypeSearch),
		Supplier: analyticSearchRequest.Supplier,
		Airline:  analyticSearchRequest.Airlines,
	}
}

// code is from analytic code severity
// message is from common-model.ResponseCode.java -- messageId (you can check in the integrator java code)
func determineAnalyticLogCodeAndMessage(err error) (bau.SearchLogHealthStatus, string) {
	switch {
	case err == nil:
		return bau.SearchLogHealthStatusSuccess, "Berhasil"
	case errors.Is(err, shared.ErrFlightNotFound) || errors.Is(err, shared.ErrInvalidRoute):
		return bau.SearchLogHealthStatusFlightNotFoundError, "Penerbangan Tidak Ditemukan"
	default:
		return bau.SearchLogHealthStatusUnknownError, "Error Tidak Diketahui"
	}
}
