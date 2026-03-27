package controller

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"{{MODULE_NAME}}/internal/service"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/util"

	commonUtil "github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/shared/util"

	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	cancelRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/cancel-book"
	cancelRS "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/cancel-book"
)

type (
	CancelBookingController struct {
		deps    deps.Deps
		service service.CancelBookingServiceItf
	}
)

func NewCancelBookingController(
	deps deps.Deps,
	service service.CancelBookingServiceItf,
) *CancelBookingController {
	return &CancelBookingController{
		deps:    deps,
		service: service,
	}
}

// @Summary		Cancel-Booking endpoint
// @Description	Endpoint for Cancel Booking process
//
// @Param			storeId		header	string									true	"Authentication header"	default(TIKETCOM)
// @Param			channelId	header	string									true	"Authentication header"	default(WEB)
// @Param			requestId	header	string									true	"Authentication header"	default(321321321)
// @Param			username	header	string									true	"Authentication header"	default(guest)
// @Param			serviceId	header	string									true	"Authentication header"	default(GATEWAY)
// @Param			resellerId	header	string									false	"Authentication header"
// @Param			identity	header	string									false	"Authentication header"
//
// @Param			message		body	cancelRQ.IntegratorCancelBookRequest	true	"Request Body"
//
// @Tags			prime-booking
// @Accept			json
// @Produce		json
// @Success		200	{object}	cancelRS.CancelBookResponse
// @Failure		400
// @Router			/cancel-booking [post]
func (ctrl *CancelBookingController) CancelBook(c echo.Context) error {
	ctxOperation := "cancel-booking"
	var (
		mr  common.MandatoryRequest
		req cancelRQ.IntegratorCancelBookRequest
		ctx = c.Request().Context()
	)

	mr = mr.BindMandatory(c)

	body, _ := io.ReadAll(c.Request().Body)
	if err := json.Unmarshal(body, &req); err != nil {
		return c.JSON(http.StatusBadRequest, common.BaseResponse[cancelRS.CancelBookResponse]{
			Message: err.Error(),
		})
	}

	logData := commonUtil.GetDataToLogFromCancelRequest(req)
	ctx = commonUtil.ConvertAndMapToSDCContext(ctx, logData)

	ctrl.deps.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.REQUEST, req))

	response, err := ctrl.service.CancelBooking(ctx, mr, req)
	if err != nil {
		ctrl.deps.GetLogger(ctx).Error(util.LogRest(mr, ctxOperation, shared.RESPONSE, err.Error()))
		// our existing java code return 200 with error message
		return c.JSON(http.StatusOK, common.ConstructResponse[*cancelRS.CancelBookResponse](GetCancelBookingResponseCodeMapper(err), err.Error(), []string{err.Error()}, nil))
	}

	ctrl.deps.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.RESPONSE, response))

	return c.JSON(http.StatusOK, common.ConstructResponse(shared.SUCCESS_RESPONSE_CODE, shared.SUCCESS_RESPONSE_CODE, nil, response))
}
