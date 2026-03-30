package controller

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"{{MODULE_NAME}}/internal/service"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/util"

	"github.com/labstack/echo/v4"
	commonUtil "github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/shared/util"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	issuedRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/issued"
	issuedRS "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/issued"
)

type (
	IssuedController struct {
		deps deps.Deps
		svc  service.IssuedProcessorService
		bau  service.AnalyticsService[issuedRQ.IntegratorIssuedRequest]
	}
)

func NewIssuedController(deps deps.Deps, svc service.IssuedProcessorService, bau service.AnalyticsService[issuedRQ.IntegratorIssuedRequest]) *IssuedController {
	return &IssuedController{deps, svc, bau}
}

// @Summary		Issued endpoint
// @Description	Endpoint for issued process
//
// @Param			storeId		header	string								true	"Authentication header"	default(TIKETCOM)
// @Param			channelId	header	string								true	"Authentication header"	default(WEB)
// @Param			requestId	header	string								true	"Authentication header"	default(321321321)
// @Param			username	header	string								true	"Authentication header"	default(guest)
// @Param			serviceId	header	string								true	"Authentication header"	default(GATEWAY)
// @Param			resellerId	header	string								false	"Authentication header"
// @Param			identity	header	string								false	"Authentication header"
//
// @Param			message		body	issuedRQ.IntegratorIssuedRequest	true	"Request Body"
//
// @Tags			prime-booking
// @Accept			json
// @Produce		json
// @Success		200	{object}	issuedRS.IssuedResponse
// @Failure		400
// @Failure		500
// @Router			/issued [post]
func (controller *IssuedController) Issued(c echo.Context) error {
	ctxOperation := "issued"
	ctx := c.Request().Context()

	mr := util.GetMandatoryRequest(c)

	var request issuedRQ.IntegratorIssuedRequest
	body, _ := io.ReadAll(c.Request().Body)
	if err := json.Unmarshal(body, &request); err != nil {
		return c.JSON(http.StatusBadRequest, common.BaseResponse[issuedRS.IssuedResponse]{
			Message: err.Error(),
		})
	}

	var err error
	analyticToken := controller.bau.CreateToken(mr, request)
	ctx = context.WithValue(ctx, shared.CTX_ANALYTIC_TOKEN_KEY, analyticToken)

	defer func() {
		controller.bau.Complete(ctx, mr, analyticToken, err)
	}()

	logData := commonUtil.GetDataToLogFromIssueRequest(request)
	ctx = commonUtil.ConvertAndMapToSDCContext(ctx, logData)

	controller.deps.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.REQUEST, request))

	response, err := controller.svc.IssuedProcessor(ctx, mr, request)
	if err != nil {
		controller.deps.GetLogger(ctx).Error(util.LogRest(mr, ctxOperation, shared.RESPONSE, err.Error()))
		return c.JSON(http.StatusOK, common.ConstructResponse[*issuedRS.IssuedResponse](GetIssuedResponseCodeMapper(err), err.Error(), []string{err.Error()}, nil))
	}

	controller.deps.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.RESPONSE, response))

	return c.JSON(http.StatusOK, common.ConstructResponse(shared.SUCCESS_RESPONSE_CODE, shared.SUCCESS_RESPONSE_CODE, nil, response))
}
