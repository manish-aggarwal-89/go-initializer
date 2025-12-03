package util

import (
	"fmt"
	"{{MODULE_NAME}}/internal/shared"

	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
)

func LogRest(mr common.MandatoryRequest, ctx, rqrs string, data interface{}) string {
	return fmt.Sprintf("[%s][%s] [%s] - %s - %s", shared.REST_IMPL, ctx, rqrs, ObjToJson(mr), ObjToJson(data))
}

func LogService(mr common.MandatoryRequest, ctx, rqrs string, data interface{}) string {
	return fmt.Sprintf("[%s][%s] [%s] - [%s] - %s", shared.SERVICE_IMPL, ctx, rqrs, mr.RequestId, ObjToJson(data))
}

func LogOutbound(mr common.MandatoryRequest, ctx, rqrs string, data string) string {
	return fmt.Sprintf("[%s][%s] [%s] - [%s] - %s", shared.OUTBOUND_IMPL, ctx, rqrs, mr.RequestId, data)
}

func LogInbound(mr common.MandatoryRequest, ctx, rqrs string, data interface{}) string {
	return fmt.Sprintf("[%s][%s] [%s] - %s - %s", shared.INBOUND_IMPL, ctx, rqrs, ObjToJson(mr), ObjToJson(data))
}

func LogScheduler(prs string, message interface{}) string {
	return fmt.Sprintf("[scheduler][%s] - %s", prs, message)
}
