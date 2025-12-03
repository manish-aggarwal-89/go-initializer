package dto

import (
	"{{MODULE_NAME}}/internal/shared/entity"
)

type CurrencyKafka map[string]float32

func (impl CurrencyKafka) DtoListCurrency() (result entity.ListCurrency) {
	if impl == nil {
		return result
	}

	for currency, rate := range impl {
		result = append(result, entity.Currency{
			CodeFrom: currency,
			Rate:     rate,
		})
	}
	return result
}
