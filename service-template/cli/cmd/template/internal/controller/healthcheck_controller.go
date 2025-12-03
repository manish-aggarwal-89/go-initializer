package controller

import (
	"{{MODULE_NAME}}/internal/service"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"

	"github.com/labstack/echo/v4"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
)

type (
	HealthCheckController struct {
		deps    deps.Deps
		service service.HealthCheckServiceItf
	}
)

// @Summary		Endpoint to check health service
// @Description	Endpoint to check health service
//
// @Param			storeId		header	string	true	"Authentication header"	default(TIKETCOM)
// @Param			channelId	header	string	true	"Authentication header"	default(WEB)
// @Param			requestId	header	string	true	"Authentication header"	default(321321321)
// @Param			username	header	string	true	"Authentication header"	default(guest)
// @Param			serviceId	header	string	true	"Authentication header"	default(GATEWAY)
// @Param			resellerId	header	string	false	"Authentication header"
// @Param			identity	header	string	false	"Authentication header"
//
// @Tags			health
// @Accept			json
// @Produce		json
// @Success		200
// @Failure		500
// @Router			/application/health [get]
func (impl *HealthCheckController) HealthCheck(c echo.Context) error {
	ctxOperation := "health_check"
	var ctx = c.Request().Context()
	mr := common.MandatoryRequest{}

	if err := impl.service.HealthCheck(ctx, mr); err != nil {
		impl.deps.GetLogger(ctx).Errorf("%s  %s error :  %s request : %s", shared.REST_IMPL, ctxOperation, err.Error(), mr)
		return c.String(400, err.Error())
	}

	return c.String(200, "ok")
}

func NewHealthCheckController(deps deps.Deps, service service.HealthCheckServiceItf) (*HealthCheckController, error) {
	return &HealthCheckController{deps: deps, service: service}, nil
}
