package service

import (
	"sync"
	"testing"
	"time"
	"{{MODULE_NAME}}/config"
	"{{MODULE_NAME}}/internal/shared/bau"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/test_var"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/log/logrus"
	commonbau "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/bau"
	book "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/book"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/fare"
)

func TestBauBookingServiceImpl_CreateToken(t *testing.T) {

	publisherMock := NewMockKafkaPublisherServiceItf(t)

	logger, err := logrus.DefaultLog()
	if err != nil {
		panic(err)
	}

	deps := deps.Deps{
		Logger: logger,
		Config: &config.Config{
			KafkaConfig: config.KafkaConfig{
				Topics: config.KafkaTopics{
					AnalyticBAUBooking: "topic",
				},
			},
		},
	}

	analyticBookingService := NewBAUBookingService(deps, publisherMock)

	req := book.IntegratorBookRequest{
		DistributionType: "distribution",
		Account: fare.Account{
			Code: "supplier",
		},
		FlightSelect: "AB 123|CD 456~EF 890",
		Itineraries: []book.Itinerary{
			{
				Departure: "CGK",
				Arrival:   "DPS",
				Airline:   "AB",
			},
			{
				Departure: "CGK",
				Arrival:   "DPS",
				Airline:   "CD",
			},
			{
				Departure: "CGK",
				Arrival:   "DPS",
				Airline:   "EF",
			},
		},
		MaxBookRetry: 0,
	}

	airlines := []string{"AB", "CD", "EF"}

	token := analyticBookingService.CreateToken(test_var.MandatoryRequest, req).(*bau.AnalyticBookingToken)

	assert.NotNil(t, &token)
	assert.Equal(t, test_var.MandatoryRequest.RequestId, token.RequestId)
	assert.Equal(t, bau.TransactionTypeBooking, token.Name)
	assert.Equal(t, test_var.MandatoryRequest.StoreId, token.StoreId)
	assert.Equal(t, test_var.MandatoryRequest.ChannelId, token.ChannelId)
	assert.Equal(t, "distribution", token.Distribution)
	assert.Equal(t, "supplier", token.Supplier)
	assert.NotNil(t, &token.Counter)
	for _, val := range airlines {
		assert.True(t, isArrayContain(token.Airlines, val))
	}
	assert.Equal(t, "CGK", token.Origin)
	assert.Equal(t, "DPS", token.Destination)
	assert.Equal(t, "ROUND_TRIP", token.TripType)
	assert.Equal(t, 0, token.RetryCount)
	assert.NotEqual(t, token.StartMillis, int64(0))
}

func TestBauBookingServiceImpl_Complete(t *testing.T) {
	logger, err := logrus.DefaultLog()
	if err != nil {
		panic(err)
	}

	publisherMock := NewMockKafkaPublisherServiceItf(t)

	deps := deps.Deps{
		Logger: logger,
		Config: &config.Config{
			KafkaConfig: config.KafkaConfig{
				Topics: config.KafkaTopics{
					AnalyticBAUBooking: "topic",
					AnalyticBAULog:     "topic",
				},
			},
		},
	}

	analyticBookingService := NewBAUBookingService(deps, publisherMock)

	token := bau.AnalyticBookingToken{
		AnalyticsToken: commonbau.AnalyticsToken{
			RequestId:    test_var.MandatoryRequest.RequestId,
			ChannelId:    test_var.MandatoryRequest.ChannelId,
			StoreId:      test_var.MandatoryRequest.StoreId,
			Distribution: "distribution",
			Supplier:     "supplier",
			Airlines:     []string{"AB", "CD", "EF"},
			Name:         bau.TransactionTypeBooking,
			StartMillis:  time.Now().UnixMilli(),
			Counter:      sync.Map{},
			MutableAnalyticsSegment: []*commonbau.AnalyticSegment{
				{
					Name:        "abc",
					StartMillis: time.Now().UnixMilli(),
					EndMillis:   time.Now().Add(time.Second).UnixMilli(),
					IsParallel:  true,
				},
				{
					Name:        "abc_1",
					StartMillis: time.Now().UnixMilli(),
					EndMillis:   time.Now().Add(time.Second).UnixMilli(),
					IsParallel:  true,
				}, {
					Name:        "def",
					StartMillis: time.Now().UnixMilli(),
					EndMillis:   time.Now().Add(time.Second).UnixMilli(),
					IsParallel:  false,
				},
			},
		},
		TripType:    "DEPARTURE",
		Origin:      "CGK",
		Destination: "DPS",
		RetryCount:  0,
	}

	publisherMock.EXPECT().Publish(test_var.MandatoryRequest, "topic", mock.Anything, false).
		Return(nil).
		Twice()

	analyticBookingService.Complete(test_var.MandatoryRequest, &token, nil)

	time.Sleep(250 * time.Millisecond)
}

func TestAnalyticsBookingToken_CreateSegment(t *testing.T) {
	token := bau.AnalyticBookingToken{
		AnalyticsToken: commonbau.AnalyticsToken{
			Counter:                 sync.Map{},
			MutableAnalyticsSegment: make([]*commonbau.AnalyticSegment, 0),
		},
	}
	token.CreateSegment("abc", false)
	assert.Equal(t, token.MutableAnalyticsSegment[0].Name, "abc")
	assert.Equal(t, token.MutableAnalyticsSegment[0].IsParallel, false)
}

func isArrayContain(arr []string, check string) bool {
	for _, val := range arr {
		if val == check {
			return true
		}
	}
	return false
}
