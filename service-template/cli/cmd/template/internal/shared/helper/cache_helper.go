package helper

import (
	"fmt"
	"{{MODULE_NAME}}/internal/shared"
)

func BuildCacheSessionKey(supplierCode string) string {
	return fmt.Sprintf("%s-session-%s", shared.CACHE_PREFIX, supplierCode)
}

func BuildCacheKeyCurrencyCodeFromCodeTo(prefix, codeFrom, codeTo string) string {
	return prefix + "-" + "com.tiket.tix.currency-currency" + "-" + codeFrom + "-" + codeTo
}
