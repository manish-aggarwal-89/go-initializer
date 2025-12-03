package helper

import (
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/entity"
)

func GetCurrencyFromExtendedData(credential entity.CredentialRepository) string {
	currency, ok := credential.ExtendedData[shared.ExtendedDataKeyCurrency]
	if !ok {
		return ""
	}

	return currency
}

func GetDomainCodeFromExtendedData(credential entity.CredentialRepository) string {
	domainCode, ok := credential.ExtendedData[shared.ExtendedDataKeyDomainCode]
	if !ok {
		return ""
	}

	return domainCode
}

func GetAgentNameFromExtendedData(credential entity.CredentialRepository) string {
	agentName, ok := credential.ExtendedData[shared.ExtendedDataKeyAgentName]
	if !ok {
		return ""
	}

	return agentName
}

func GetAgenPasswordFromExtendedData(credential entity.CredentialRepository) string {
	agentPwd, ok := credential.ExtendedData[shared.ExtendedDataKeyAgentPassword]
	if !ok {
		return ""
	}

	return agentPwd
}

func GetOrganizationIdFromExtendedData(credential entity.CredentialRepository) string {
	organizationCode, ok := credential.ExtendedData[shared.ExtendedDataKeyOrganizationId]
	if !ok {
		return ""
	}

	return organizationCode
}

func GetPaymentMethodFromExtendedData(credential entity.CredentialRepository) string {
	paymentMethod, ok := credential.ExtendedData[shared.ExtendedDataKeyPaymentMethod]
	if !ok {
		return ""
	}

	return paymentMethod
}
