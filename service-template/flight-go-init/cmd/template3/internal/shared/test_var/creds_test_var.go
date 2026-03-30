package test_var

//import (
//	"time"
//	"{{MODULE_NAME}}/internal/shared"
//	"{{MODULE_NAME}}/internal/shared/dto"
//	"{{MODULE_NAME}}/internal/shared/entity"
//	"{{MODULE_NAME}}/internal/shared/util"
//
//	"go.mongodb.org/mongo-driver/bson/primitive"
//)
//
//const (
//	CredentialIdString       = "65097c3297a071ee3a5fe554"
//	CredentialSupplierString = "string"
//	CredentialJson           = `{"distributionType":"string","expiredDate":"string","extendedData":{"additionalProp1":"string","additionalProp2":"string","additionalProp3":"string"},"password":"string","supplier":"string","username":"string"}`
//)
//
//var (
//	TimeNow                = time.Now()
//	CredentialPredefineKey = shared.PREDEFINED_CREDENTIAL_CACHE_KEY_PREFIX + CredentialSupplierString
//	CredentialId, _        = primitive.ObjectIDFromHex(CredentialIdString)
//
//	CREDENTIAL = entity.Credential{
//		DistributionType: "string",
//		Supplier:         CredentialSupplierString,
//		Username:         "string",
//		Password:         "string",
//		ExpiredDate:      "2023-09-25",
//		ExtendedData:     ExtendedData(),
//	}
//
//	CredentialRepo = entity.CredentialRepository{
//		DistributionType: "string",
//		Supplier:         CredentialSupplierString,
//		Username:         "string",
//		Password:         "string",
//		ExpiredDate:      &TimeNow,
//		ExtendedData:     ExtendedData(),
//		IsDeleted:        0,
//		Version:          1,
//		CreateDate:       TimeNow,
//		UpdateDate:       TimeNow,
//		CreateBy:         MandatoryRequest.Username,
//		UpdateBy:         MandatoryRequest.Username,
//	}
//
//	CredentialRepos = []entity.CredentialRepository{CredentialRepo}
//
//	CredentialLog = entity.CredentialLog{
//		Value: &CredentialRepo,
//	}
//
//	CredentialLogRepo = entity.CredentialLogRepository{
//		ID:         CredentialId,
//		Value:      &CredentialRepo,
//		IsDeleted:  0,
//		Version:    1,
//		CreateDate: TimeNow,
//		UpdateDate: TimeNow,
//		CreateBy:   MandatoryRequest.Username,
//		UpdateBy:   MandatoryRequest.Username,
//	}
//
//	CredentialLogRepos = []entity.CredentialLogRepository{CredentialLogRepo}
//
//	CredentialLogResp = dto.CredentialLogResponse{
//		LastUpdatedDate: util.FormatTime(util.DateTimeFormatWithSecond, CredentialLogRepo.CreateDate),
//		LastUpdatedBy:   CredentialLogRepo.CreateBy,
//		LogData:         util.ObjToJson(CredentialLogRepo.Value),
//	}
//
//	CredentialLogResps = []dto.CredentialLogResponse{CredentialLogResp}
//
//	CredentialFilter = dto.CredentialFilter{
//		DistributionType: "string",
//		Supplier:         CredentialSupplierString,
//		Username:         "string",
//	}
//
//	CREDENTIAL_ENTITY_DATA = entity.Credential{
//		DistributionType: "{{PROVIDER}}",
//		Supplier:         CredentialSupplierString,
//		Username:         "string",
//		Password:         "string",
//		ExpiredDate:      "2023-09-25",
//		ExtendedData:     ExtendedDataAgent(),
//	}
//)
//
//func ExtendedData() map[string]string {
//	m := make(map[string]string)
//	m["string"] = "string"
//	m[shared.ExtendedDataKeyDomainCode] = shared.ExtendedDataKeyDomainCode
//	return m
//}
//
//func ExtendedDataAgent() map[string]string {
//	m := make(map[string]string)
//	m[shared.ExtendedDataKeyCurrency] = CredentialExtendedDataCurrencyIDR
//	m[shared.ExtendedDataKeyDomainCode] = CredentialExtendedDataDomainCode
//	m[shared.ExtendedDataKeyAgentName] = CredentialExtendedDataAgentName
//	m[shared.ExtendedDataKeyAgentPassword] = CredentialExtendedDataKeyAgentPassword
//	return m
//}
