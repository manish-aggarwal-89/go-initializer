package util

import (
	"encoding/json"
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ObjToJson(data interface{}) string {
	rsp, _ := json.Marshal(data)
	return string(rsp)
}

func GetMandatoryRequest(c echo.Context) common.MandatoryRequest {
	mr := common.MandatoryRequest{}
	return mr.BindMandatory(c)
}

func ObjectIDFromString(s string) (primitive.ObjectID, error) {
	objID, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		return objID, errors.New("INVALID ID")
	}
	return objID, nil
}
