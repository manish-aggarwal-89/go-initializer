package service

import (
	"context"
	"errors"
	"time"
	"{{MODULE_NAME}}/internal/repository"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/dto"
	"{{MODULE_NAME}}/internal/shared/entity"
	"{{MODULE_NAME}}/internal/shared/helper"
	"{{MODULE_NAME}}/internal/shared/util"

	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
)

var now = time.Now

type (
	CurrencyService interface {
		Upsert(ctx context.Context, currencyKafka dto.CurrencyKafka) error
		GetCurrencyCodeFromCodeToV1(ctx context.Context, mr common.MandatoryRequest, codeFrom string, codeTo string) (entity.CurrencyRepository, error)
	}

	CurrencyImpl struct {
		deps               deps.Deps
		currencyRepository repository.CurrencyRepositoryInterface
	}
)

func NewCurrencyService(deps deps.Deps, currencyRepository repository.CurrencyRepositoryInterface) CurrencyService {
	return &CurrencyImpl{deps: deps, currencyRepository: currencyRepository}
}

func (impl *CurrencyImpl) GetCurrencyCodeFromCodeToV1(ctx context.Context, mr common.MandatoryRequest, codeFrom string, codeTo string) (entity.CurrencyRepository, error) {
	ctxOperation := "currency_get_currency_codefrom_v1"

	cacheKey := helper.BuildCacheKeyCurrencyCodeFromCodeTo(impl.deps.Config.CurrencyV1Config.CacheKeyCurrencyPrefixIntegrator, codeFrom, codeTo)
	dataCurrency, ok := impl.deps.PredefinedCache.Predefined.Get(cacheKey)
	if !ok {
		currencyData, err := impl.currencyRepository.FindByCodeFrom(ctx, codeFrom)
		if err != nil {
			impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err.Error()))
			return entity.CurrencyRepository{}, err
		}
		impl.deps.PredefinedCache.Predefined.Set(cacheKey, currencyData, impl.deps.Config.CurrencyV1Config.CacheCurrencyTTLV1*time.Minute)
		return currencyData, nil
	}

	result, _ := dataCurrency.(entity.CurrencyRepository)
	return result, nil
}

func (impl *CurrencyImpl) Upsert(ctx context.Context, currencyKafka dto.CurrencyKafka) error {
	ctxOperation := "currency_upsert"

	for currencyCode, rate := range currencyKafka {
		mr := common.MandatoryRequest{}.Default()
		currencyData, err := impl.currencyRepository.FindByCodeFrom(ctx, currencyCode)
		if errors.Is(err, shared.ErrDataNotExist) {
			_, err := impl.currencyRepository.Create(ctx, mr, entity.Currency{CodeFrom: currencyCode, Rate: rate})
			if err != nil {
				impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err.Error()))
				continue
			}

			cacheKey := helper.BuildCacheKeyCurrencyCodeFromCodeTo(impl.deps.Config.CurrencyV1Config.CacheKeyCurrencyPrefixIntegrator, currencyCode, shared.CURRENCY_IDR)
			impl.deps.PredefinedCache.Predefined.Set(cacheKey, currencyData, impl.deps.Config.CurrencyV1Config.CacheCurrencyTTLV1*time.Minute)
			continue
		}
		if err != nil {
			impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err.Error()))
			continue
		}

		currencyData.BuyRate = rate
		currencyData.SellRate = rate
		currencyData.UpdateBy = mr.Username
		currencyData.UpdateDate = now()
		currencyData.Version++
		err = impl.currencyRepository.UpdateByID(ctx, currencyData)
		if err != nil {
			impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err.Error()))
			continue
		}
		cacheKey := helper.BuildCacheKeyCurrencyCodeFromCodeTo(impl.deps.Config.CurrencyV1Config.CacheKeyCurrencyPrefixIntegrator, currencyCode, shared.CURRENCY_IDR)
		impl.deps.PredefinedCache.Predefined.Set(cacheKey, currencyData, impl.deps.Config.CurrencyV1Config.CacheCurrencyTTLV1*time.Minute)
	}

	return nil
}
