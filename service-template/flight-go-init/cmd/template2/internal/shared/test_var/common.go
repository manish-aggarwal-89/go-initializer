package test_var

import (
	"encoding/json"
	"{{MODULE_NAME}}/internal/shared"

	commonModel "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	fareRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/fare"
)

var (
	MandatoryRequest = commonModel.MandatoryRequest{
		StoreId:   "TIKETCOM",
		ChannelId: "iOS",
		RequestId: "REQUEST_ID",
		ServiceId: "LOGIN",
		Username:  "testuser",
	}

	DISTRIBUTION_TYPE = shared.DISTRIBUTION_TYPE

	ACCOUNT = fareRQ.Account{
		Name: shared.SUPPLIER_CODE,
		Code: shared.SUPPLIER_CODE,
	}

	SUPPLIER_MAPPING = fareRQ.SupplierMapping{
		Account:  ACCOUNT,
		Airlines: []string{"KS"},
	}

	SUPPLIER_REQUEST = fareRQ.SupplierRequest{
		Accounts: []fareRQ.Account{ACCOUNT},
		Airlines: []string{"KS"},
	}

	CURRENCY = "IDR"

	SUPPLIER_CODE = "{{PROVIDER}}"
)

func Marshal(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
