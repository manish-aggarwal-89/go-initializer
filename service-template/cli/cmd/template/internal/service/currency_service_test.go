package service

import (
	"context"
	"testing"
	"time"
	"{{MODULE_NAME}}/config"
	"{{MODULE_NAME}}/internal/repository"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/dto"
	"{{MODULE_NAME}}/internal/shared/entity"
	"{{MODULE_NAME}}/internal/shared/helper"

	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/log/logrus"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
)

func TestGetCurrencyCodeFromCodeToV1(t *testing.T) {
	logger, err := logrus.DefaultLog()
	if err != nil {
		panic(err)
	}

	c := cache.New(0, 0)

	deps := deps.Deps{
		Logger: logger,
		PredefinedCache: &deps.PredefinedCache{
			Predefined: *c,
		},
		Config: &config.Config{
			CurrencyV1Config: config.CurrencyV1Config{
				CacheCurrencyTTLV1: 1,
			},
		},
	}
	currencyRepositoryMock := repository.NewMockCurrencyRepositoryInterface(t)

	currencyImpl := NewCurrencyService(deps, currencyRepositoryMock)

	tests := []struct {
		name        string
		stub        func()
		codeFrom    string
		codeTo      string
		expected    entity.CurrencyRepository
		expectedErr error
	}{
		{
			name: "success_get_from_database",
			stub: func() {
				currencyRepositoryMock.EXPECT().FindByCodeFrom(mock.Anything, "MYR").
					Return(entity.CurrencyRepository{
						CodeFrom: "MYR",
						CodeTo:   "IDR",
						BuyRate:  0.5,
						SellRate: 1.5,
					}, nil).
					Once()
			},
			codeFrom: "MYR",
			codeTo:   "IDR",
			expected: entity.CurrencyRepository{
				CodeFrom: "MYR",
				CodeTo:   "IDR",
				BuyRate:  0.5,
				SellRate: 1.5,
			},
			expectedErr: nil,
		},
		{
			name: "success_get_from_predefine",
			stub: func() {
				cacheKey := helper.BuildCacheKeyCurrencyCodeFromCodeTo(deps.Config.CurrencyV1Config.CacheKeyCurrencyPrefixIntegrator, "MYR", "IDR")
				c.Set(cacheKey, entity.CurrencyRepository{
					CodeFrom: "MYR",
					CodeTo:   "IDR",
					BuyRate:  0.5,
					SellRate: 1.5,
				}, 0)
			},
			codeFrom: "MYR",
			codeTo:   "IDR",
			expected: entity.CurrencyRepository{
				CodeFrom: "MYR",
				CodeTo:   "IDR",
				BuyRate:  0.5,
				SellRate: 1.5,
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stub()

			actual, err := currencyImpl.GetCurrencyCodeFromCodeToV1(context.Background(), common.MandatoryRequest{}, tt.codeFrom, tt.codeTo)

			require.Equal(t, tt.expected, actual)
			require.Equal(t, tt.expectedErr, err)

			c.Flush()
		})
	}
}

func TestUpsert(t *testing.T) {
	logger, err := logrus.DefaultLog()
	if err != nil {
		panic(err)
	}

	now = func() time.Time {
		return time.Date(2024, 10, 10, 0, 0, 0, 0, time.UTC)
	}

	c := cache.New(0, 0)

	deps := deps.Deps{
		Logger: logger,
		PredefinedCache: &deps.PredefinedCache{
			Predefined: *c,
		},
		Config: &config.Config{
			CurrencyV1Config: config.CurrencyV1Config{
				CacheCurrencyTTLV1: 1,
			},
		},
	}
	currencyRepositoryMock := repository.NewMockCurrencyRepositoryInterface(t)

	currencyImpl := NewCurrencyService(deps, currencyRepositoryMock)

	tests := []struct {
		name          string
		stub          func()
		currencyKafka dto.CurrencyKafka
		expectedErr   error
	}{
		{
			name: "success_insert",
			stub: func() {
				currencyRepositoryMock.EXPECT().FindByCodeFrom(mock.Anything, "MYR").
					Return(entity.CurrencyRepository{}, shared.ErrDataNotExist).
					Once()

				currencyRepositoryMock.EXPECT().Create(mock.Anything, mock.Anything, entity.Currency{CodeFrom: "MYR", Rate: 0.5}).
					Return(entity.CurrencyRepository{}, nil).
					Once()
			},
			currencyKafka: dto.CurrencyKafka{
				"MYR": 0.5,
			},
			expectedErr: nil,
		},
		{
			name: "success_update",
			stub: func() {
				currencyRepositoryMock.EXPECT().FindByCodeFrom(mock.Anything, "MYR").
					Return(entity.CurrencyRepository{
						CodeFrom: "MYR",
						CodeTo:   "IDR",
						BuyRate:  0.3,
						SellRate: 1.3,
						Version:  0,
					}, nil).
					Once()

				currencyRepository := entity.CurrencyRepository{
					CodeFrom:   "MYR",
					CodeTo:     "IDR",
					BuyRate:    0.5,
					SellRate:   0.5,
					Version:    1,
					UpdateDate: now(),
					UpdateBy:   "System",
				}

				currencyRepositoryMock.EXPECT().UpdateByID(mock.Anything, currencyRepository).
					Return(nil).
					Once()
			},
			currencyKafka: dto.CurrencyKafka{
				"MYR": 0.5,
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stub()

			err := currencyImpl.Upsert(context.Background(), tt.currencyKafka)

			require.Equal(t, tt.expectedErr, err)

			c.Flush()
		})
	}
}
