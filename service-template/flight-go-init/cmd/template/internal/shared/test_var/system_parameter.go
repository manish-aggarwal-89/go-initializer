package test_var

import (
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var SYSTEM_PARAMETER_PUBLISH_ON_BOOKING_FLAG_TRUE = entity.SystemParameterRepository{
	Variable: shared.SYS_PARAM_PUBLISH_ON_BOOKING_FLAG_KEY,
	Value:    "true",
}

var SYSTEM_PARAMETER_PUBLISH_ON_BOOKING_FLAG_FALSE = entity.SystemParameterRepository{
	Variable: shared.SYS_PARAM_PUBLISH_ON_BOOKING_FLAG_KEY,
	Value:    "false",
}

var SYSTEM_PARAMETER_CACHE_SSR_FLAG_KEY_TRUE = entity.SystemParameterRepository{
	Variable: shared.SYS_PARAM_CACHE_SSR_FLAG_KEY,
	Value:    "true",
}

var SYSTEM_PARAMETER_CACHE_SSR_FLAG_KEY_FALSE = entity.SystemParameterRepository{
	Variable: shared.SYS_PARAM_CACHE_SSR_FLAG_KEY,
	Value:    "false",
}

var SYSTEM_PARAMETER_BAG_10KG = entity.SystemParameterRepository{
	Variable: shared.SYS_PARAM_ADDON_BAGGAGE + ".10.KG",
	Value:    "BAG1",
}

var SYSTEM_PARAMETER_BAG_20KG = entity.SystemParameterRepository{
	Variable: shared.SYS_PARAM_ADDON_BAGGAGE + ".20.KG",
	Value:    "BAG2",
}

var SYSTEM_PARAMETER_INCLUSIVE_BAGGAGE_CABIN = entity.SystemParameterRepository{
	Variable: shared.SYS_PARAM_INCLUSIVE_BAGGAGE_CABIN,
	Value:    "20|KG",
}

var LIST_SYSTEM_PARAMETER = []entity.SystemParameterRepository{SYSTEM_PARAMETER_PUBLISH_ON_BOOKING_FLAG_FALSE}
var LIST_SYSTEM_PARAMETER_BAGGAGE_CONFIG = []entity.SystemParameterRepository{
	SYSTEM_PARAMETER_BAG_10KG, SYSTEM_PARAMETER_BAG_20KG,
}

var SYSTEM_PARAMETER_ID_STR = "65097c3297a071ee3a5fe554"
var SYSTEM_PARAMETER_ID, _ = primitive.ObjectIDFromHex("65097c3297a071ee3a5fe554")
var SYSTEM_PARAMETER_ID_ERROR_STR = "65097c3297a071ee3a5fe558"
var SYSTEM_PARAMETER_ID_ERROR, _ = primitive.ObjectIDFromHex("65097c3297a071ee3a5fe558")
var VARIABLE = "VARIABLE"

var SYSTEM_PARAMETER_REQUEST = entity.SystemParameter{
	Variable: shared.SYS_PARAM_PUBLISH_ON_BOOKING_FLAG_KEY,
	Value:    "false",
}

var SYSTEM_PARAMETER_JSON = `{"variable":"testing","value":"testing","description":"param testing"}`

var SYSTEM_PARAMETER_CHUNK_FLAG_TRUE = entity.SystemParameterRepository{
	Variable: shared.SYS_PARAM_CHUNK_FLAG_KEY,
	Value:    "true",
}

var SYSTEM_PARAMETER_CHUNK_FLAG_FALSE = entity.SystemParameterRepository{
	Variable: shared.SYS_PARAM_CHUNK_FLAG_KEY,
	Value:    "false",
}

var SYSTEM_PARAMETER_SEARCH_SSR_CODES = entity.SystemParameterRepository{
	Variable: shared.SYS_PARAM_SEARCH_SSR_CODES,
	Value:    "PBAB,PBAC,PBAD,PBAE,PBAF",
}

var SYSTEM_PARAMETER_CHUNK_SIZE_100 = entity.SystemParameterRepository{
	Variable: shared.SYS_PARAM_CHUNK_SIZE_KEY,
	Value:    "100",
}
