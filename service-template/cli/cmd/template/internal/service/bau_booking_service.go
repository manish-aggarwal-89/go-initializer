package service

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/bau"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/helper"
	"{{MODULE_NAME}}/internal/shared/util"

	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/log/logrus"
	commonbau "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/bau"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	bookRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/book"
)

type BAUBookingService struct {
	deps deps.Deps

	kafkaPublisherService KafkaPublisherServiceItf
}

func NewBAUBookingService(deps deps.Deps, kafkaPublisherService KafkaPublisherServiceItf) bau.AnalyticsService[bookRQ.IntegratorBookRequest] {
	return &BAUBookingService{
		deps:                  deps,
		kafkaPublisherService: kafkaPublisherService,
	}
}

func (svc *BAUBookingService) CreateToken(mandatoryRequest common.MandatoryRequest, integratorBookRequest bookRQ.IntegratorBookRequest) bau.AnalyticsTokenModifier {

	tripType := helper.DefineTripType(integratorBookRequest)
	return &bau.AnalyticBookingToken{
		AnalyticsToken: commonbau.AnalyticsToken{
			Name:                    bau.TransactionTypeBooking,
			RequestId:               mandatoryRequest.RequestId,
			StoreId:                 mandatoryRequest.StoreId,
			ChannelId:               mandatoryRequest.ChannelId,
			Distribution:            integratorBookRequest.DistributionType,
			Supplier:                integratorBookRequest.Account.Code,
			Airlines:                helper.BuildAirlines(integratorBookRequest),
			StartMillis:             time.Now().UnixMilli(),
			MutableAnalyticsSegment: make([]*commonbau.AnalyticSegment, 0),
			Counter:                 sync.Map{},
		},
		Origin:      integratorBookRequest.Itineraries[0].Departure,
		Destination: integratorBookRequest.Itineraries[0].Arrival,
		TripType:    string(tripType),
	}
}

func (svc *BAUBookingService) Complete(mr common.MandatoryRequest, token bau.AnalyticsTokenModifier, businessErr error) {
	token.Complete()
	bookingToken := token.(*bau.AnalyticBookingToken)

	healthStatus := bau.GetBookingHealthStatus(businessErr)
	analyticBookingRequest := constructAnalyticBookingRequest(*bookingToken, healthStatus)
	ctx := logrus.InjectRequestMetadataToContext(context.Background(), bookingToken.RequestId, bookingToken.Name, "")

	b, err := json.Marshal(analyticBookingRequest)
	if err != nil {
		svc.deps.GetLogger(ctx).Errorf("%s %s - %s", "[BAU_Booking]", bookingToken.RequestId, err.Error())
		return
	}

	go func() {
		err := svc.kafkaPublisherService.Publish(mr, svc.deps.Config.KafkaConfig.Topics.AnalyticBAUBooking, string(b), false)
		if err != nil {
			svc.deps.GetLogger(ctx).Errorf("%s %s - error publish to kafka, with err: %s", "[BAU_Booking]", bookingToken.RequestId, err.Error())
		}
	}()

	analyticBookLogRecord := constructAnalyticBookingLogRecord(analyticBookingRequest, businessErr)
	bAnalyticLog, err := json.Marshal(analyticBookLogRecord)
	if err != nil {
		svc.deps.GetLogger(ctx).Errorf("%s %s - %s", "[BAU_Log_Booking]", bookingToken.RequestId, err.Error())
		return
	}

	go func() {
		err := svc.kafkaPublisherService.Publish(mr, svc.deps.Config.KafkaConfig.Topics.AnalyticBAULog, string(bAnalyticLog), false)
		if err != nil {
			svc.deps.Logger.Errorf("%s %s - error publish to kafka, with err: %s", "[BAU_Log_Booking]", bookingToken.RequestId, err.Error())
		}
	}()
}

func constructAnalyticBookingRequest(token bau.AnalyticBookingToken, healthStatus bau.BookingHealthStatus) bau.AnalyticBookingRequest {

	airlineProcessTimeDetail := bau.ConstructAirlineProcessTimeDetail(token.AnalyticsToken.MutableAnalyticsSegment)
	airlineProcessTime := bau.SumMap(airlineProcessTimeDetail)

	transactionDate := time.UnixMilli(token.StartMillis).Format(util.DateTimeFormatWithSecond)

	return bau.AnalyticBookingRequest{
		RequestID:                token.RequestId,
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
		RetryCount:               token.RetryCount,
		BookingHealthStatus:      healthStatus,
		StoreID:                  token.StoreId,
		ChannelID:                token.ChannelId,
	}
}

func constructAnalyticBookingLogRecord(analyticBookRequest bau.AnalyticBookingRequest, err error) commonbau.AnalyticLogRecord {
	code, message := determineAnalyticBookingLogCodeAndMessage(err)

	return commonbau.AnalyticLogRecord{
		Code:     code.String(),
		Message:  message,
		Event:    strings.ToUpper(bau.TransactionTypeBooking),
		Supplier: analyticBookRequest.Supplier,
		Airline:  analyticBookRequest.Airlines,
	}
}

// code is from analytic code severity
// message is from common-model.ResponseCode.java -- messageId (you can check in the integrator java code)
func determineAnalyticBookingLogCodeAndMessage(err error) (bau.BookingLogHealthStatus, string) {
	switch err {
	case nil:
		return bau.BookingLogHealthStatusSuccess, "Berhasil"
	case shared.ErrBookingFailed:
		return bau.BookingLogHealthStatusBookingFailed, "Gagal Booking"
	case shared.ErrNoSeat, shared.ErrNoRoute:
		return bau.BookingLogHealthStatusNoSeat, "Kursi Tidak Tersedia"
	default:
		return bau.BookingLogHealthStatusUnknownError, "Error Tidak Diketahui"
	}
}
