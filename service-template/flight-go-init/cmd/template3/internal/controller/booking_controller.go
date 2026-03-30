package controller

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"{{MODULE_NAME}}/internal/service"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"

	"github.com/labstack/echo/v4"
	util2 "github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/log/util"
	commonUtil "github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/shared/util"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	bookRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/book"
	bookRS "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/book"
)

type BookingController struct {
	deps              deps.Deps
	bookingProcessor  service.BookingProcessorService
	bauBookingService service.AnalyticsService[bookRQ.IntegratorBookRequest]
}

func NewBookingController(
	deps deps.Deps,
	bookingProcessor service.BookingProcessorService,
	bauBookingService service.AnalyticsService[bookRQ.IntegratorBookRequest],
) *BookingController {
	return &BookingController{
		deps:              deps,
		bookingProcessor:  bookingProcessor,
		bauBookingService: bauBookingService,
	}
}

// @Summary		Booking endpoint
// @Description	endpoint for booking process
//
// @Param			storeId		header	string							true	"Authentication header"	default(TIKETCOM)
// @Param			channelId	header	string							true	"Authentication header"	default(WEB)
// @Param			requestId	header	string							true	"Authentication header"	default(321321321)
// @Param			username	header	string							true	"Authentication header"	default(guest)
// @Param			serviceId	header	string							true	"Authentication header"	default(GATEWAY)
// @Param			resellerId	header	string							false	"Authentication header"
// @Param			identity	header	string							false	"Authentication header"
//
// @Param			message		body	bookRQ.IntegratorBookRequest	true	"Request Body"
//
// @Tags			prime-booking
// @Accept			json
// @Produce		json
// @Success		200	{object}	bookRS.BookResponse
// @Failure		400
// @Router			/booking [post]
func (ctrl *BookingController) Booking(c echo.Context) error {
	ctxOperation := "booking"
	var (
		mr  common.MandatoryRequest
		req bookRQ.IntegratorBookRequest
		ctx = c.Request().Context()
	)
	mr = mr.BindMandatory(c)

	body, _ := io.ReadAll(c.Request().Body)
	if err := json.Unmarshal(body, &req); err != nil {
		return c.JSON(http.StatusBadRequest,
			common.ConstructResponse[*bookRS.BookResponse](shared.FAILURE_RESPONSE_CODE, shared.FAILURE_RESPONSE_CODE, []string{"unable to unmarshal request : " + err.Error()}, nil))
	}

	var err error
	analyticToken := ctrl.bauBookingService.CreateToken(mr, req)
	ctx = context.WithValue(ctx, shared.CTX_ANALYTIC_TOKEN_KEY, analyticToken)
	defer func() {
		ctrl.bauBookingService.Complete(ctx, mr, analyticToken, err)
	}()
	logData := commonUtil.GetDataToLogFromBooking(req)
	ctx = commonUtil.ConvertAndMapToSDCContext(ctx, logData)
	fields, msg := util2.LogRestV2(mr, ctxOperation, shared.REQUEST, req)
	ctrl.deps.GetLogger(ctx).WithFields(fields).Info(msg)
	res, err := ctrl.bookingProcessor.BookingProcessor(ctx, mr, req)
	if err != nil {
		fields, msg = util2.LogRestV2(mr, ctxOperation, shared.RESPONSE, err.Error(), "error in booking controller")
		ctrl.deps.GetLogger(ctx).WithError(err).WithFields(fields).Error(msg)
		return c.JSON(http.StatusOK, common.ConstructResponse[*bookRS.BookResponse](BookingResponseCodeMapper(err), err.Error(), []string{err.Error()}, nil))
	}
	responseFields, responseMsg := util2.LogRestV2(mr, ctxOperation, shared.RESPONSE, res)
	ctrl.deps.GetLogger(ctx).WithFields(responseFields).Info(responseMsg)
	return c.JSON(http.StatusOK, common.ConstructResponse(shared.SUCCESS_RESPONSE_CODE, shared.SUCCESS_RESPONSE_CODE, nil, res))
}
