package controller

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"{{MODULE_NAME}}/internal/service"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/bau"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/util"

	commonUtil "github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/shared/util"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	bookRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/book"
	bookRS "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/book"
)

type BookingController struct {
	deps              deps.Deps
	bookingService    service.BookingServiceItf
	bauBookingService bau.AnalyticsService[bookRQ.IntegratorBookRequest]
}

func NewBookingController(
	deps deps.Deps,
	bookingService service.BookingServiceItf,
	bauBookingService bau.AnalyticsService[bookRQ.IntegratorBookRequest],
) *BookingController {
	return &BookingController{
		deps:              deps,
		bookingService:    bookingService,
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
// @Router			/{{BASE_PATH}}/booking [post]
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
		ctrl.bauBookingService.Complete(mr, analyticToken, err)
	}()

	logData := commonUtil.GetDataToLogFromBooking(req)
	ctx = commonUtil.ConvertAndMapToSDCContext(ctx, logData)

	ctrl.deps.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.REQUEST, req))

	res, err := ctrl.bookingService.Booking(ctx, mr, req)
	if err != nil {
		ctrl.deps.GetLogger(ctx).Error(util.LogRest(mr, ctxOperation, shared.RESPONSE, err.Error()))
		// our existing java code return 200 with error message
		return c.JSON(http.StatusOK, common.ConstructResponse[*bookRS.BookResponse](BookingResponseCodeMapper(err), err.Error(), []string{err.Error()}, nil))
	}

	ctrl.deps.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.RESPONSE, res))
	return c.JSON(http.StatusOK, common.ConstructResponse(shared.SUCCESS_RESPONSE_CODE, shared.SUCCESS_RESPONSE_CODE, nil, res))
}
