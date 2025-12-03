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
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/fare"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/issued"
)

func TestBauIssuanceServiceImpl_CreateToken(t *testing.T) {

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
					AnalyticBAUIssuance: "topic",
				},
			},
		},
	}

	analyticIssuedService := NewBAUIssuedService(deps, publisherMock)

	req := issued.IntegratorIssuedRequest{
		DistributionType: "distribution",
		Account: fare.Account{
			Code: "supplier",
		},
		FlightSelect: "AB 123|CD 456~EF 890",
		BookingCode:  "bookingCode",
	}

	airlines := []string{"AB", "CD", "EF"}

	token := analyticIssuedService.CreateToken(test_var.MandatoryRequest, req).(*bau.AnalyticsIssuanceToken)

	assert.NotNil(t, &token)
	assert.Equal(t, token.RequestId, test_var.MandatoryRequest.RequestId)
	assert.Equal(t, token.Name, bau.TransactionTypeIssuance)
	assert.Equal(t, token.StoreId, test_var.MandatoryRequest.StoreId)
	assert.Equal(t, token.ChannelId, test_var.MandatoryRequest.ChannelId)
	assert.Equal(t, token.Distribution, "distribution")
	assert.Equal(t, token.Supplier, "supplier")
	assert.NotNil(t, &token.Counter)
	assert.Equal(t, token.Airlines, airlines)
	assert.Equal(t, token.BookingCode, "bookingCode")
	assert.NotEqual(t, token.StartMillis, int64(0))
}

func TestBauIssuanceServiceImpl_Complete(t *testing.T) {
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
					AnalyticBAUIssuance: "topic",
					AnalyticBAULog:      "topic",
				},
			},
		},
	}

	analyticIssuedService := NewBAUIssuedService(deps, publisherMock)

	token := bau.AnalyticsIssuanceToken{
		AnalyticsToken: commonbau.AnalyticsToken{
			RequestId:    test_var.MandatoryRequest.RequestId,
			ChannelId:    test_var.MandatoryRequest.ChannelId,
			StoreId:      test_var.MandatoryRequest.StoreId,
			Distribution: "distribution",
			Supplier:     "supplier",
			Airlines:     []string{"AB", "CD", "EF"},
			Name:         bau.TransactionTypeIssuance,
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
		BookingCode: "bookingCode",
	}

	publisherMock.EXPECT().Publish(test_var.MandatoryRequest, "topic", mock.Anything, false).
		Return(nil).
		Twice()

	analyticIssuedService.Complete(test_var.MandatoryRequest, &token, nil)

	time.Sleep(250 * time.Millisecond)
}

func TestAnalyticsIssuanceToken_CreateSegment(t *testing.T) {
	token := bau.AnalyticsIssuanceToken{
		AnalyticsToken: commonbau.AnalyticsToken{
			Counter:                 sync.Map{},
			MutableAnalyticsSegment: make([]*commonbau.AnalyticSegment, 0),
		},
	}
	token.CreateSegment("abc", false)
	assert.Equal(t, token.MutableAnalyticsSegment[0].Name, "abc")
	assert.Equal(t, token.MutableAnalyticsSegment[0].IsParallel, false)
}
