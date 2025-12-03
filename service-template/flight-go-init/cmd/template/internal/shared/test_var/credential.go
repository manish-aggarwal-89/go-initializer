package test_var

import (
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/entity"
)

const (
	CredentialUsernameTiketComFlyjaya = ""
	CredentialPasswordTiketComFlyjaya = ""

	CredentialExtendedDataCurrencyIDR       = "IDR"
	CredentialExtendedDataDomainCode        = "EXT"
	CredentialExtendedDataAgentName         = "API_TIKETCORP"
	CredentialExtendedDataKeyAgentPassword  = "@p1_tiket19"
	CredentialExtendedDataKeyOrganizationId = "organization_id"
	CredentialExtendedDataKeyPaymentMethod  = "method"
)

var (
	CredentialTiketComFlyJaya = entity.CredentialRepository{
		Supplier:         shared.SUPPLIER_CODE,
		DistributionType: shared.DISTRIBUTION_TYPE,
		Username:         CredentialUsernameTiketComFlyjaya,
		Password:         CredentialPasswordTiketComFlyjaya,
		IsStaging:        true,
		ExtendedData: map[string]string{
			shared.ExtendedDataKeyCurrency:       CredentialExtendedDataCurrencyIDR,
			shared.ExtendedDataKeyDomainCode:     CredentialExtendedDataDomainCode,
			shared.ExtendedDataKeyAgentName:      CredentialExtendedDataAgentName,
			shared.ExtendedDataKeyAgentPassword:  CredentialExtendedDataKeyAgentPassword,
			shared.ExtendedDataKeyOrganizationId: CredentialExtendedDataKeyOrganizationId,
			shared.ExtendedDataKeyPaymentMethod:  CredentialExtendedDataKeyPaymentMethod,
		},
	}
)
