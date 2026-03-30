package controller

import (
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
	fareRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/fare"
	fareRS "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/fare"
)

type (
	SearchControllerImpl struct {
		deps               deps.Deps
		searchProcessorSvc service.SearchProcessorService
	}
)

func NewSearchController(deps deps.Deps, searchProcessorSvc service.SearchProcessorService) *SearchControllerImpl {
	return &SearchControllerImpl{deps: deps, searchProcessorSvc: searchProcessorSvc}
}

// @Summary		Search endpoint
// @Description	Endpoint for search process
//
// @Param			storeId		header	string							true	"Authentication header"	default(TIKETCOM)
// @Param			channelId	header	string							true	"Authentication header"	default(WEB)
// @Param			requestId	header	string							true	"Authentication header"	default(321321321)
// @Param			username	header	string							true	"Authentication header"	default(guest)
// @Param			serviceId	header	string							true	"Authentication header"	default(GATEWAY)
// @Param			resellerId	header	string							false	"Authentication header"
// @Param			identity	header	string							false	"Authentication header"
//
// @Param			message		body	fareRQ.IntegratorFareRequest	true	"Request Body"
//
// @Tags			prime-booking
// @Accept			json
// @Produce		json
// @Success		200	{object}	fareRS.FlightIntegratorSearchResponse
// @Failure		400
// @Router			/search [post]
func (ctrl *SearchControllerImpl) Search(c echo.Context) error {
	ctxOperation := "search"
	var ctx = c.Request().Context()

	mr := common.MandatoryRequest{}
	mr = mr.BindMandatory(c)

	var request fareRQ.IntegratorFareRequest
	body, _ := io.ReadAll(c.Request().Body)
	if err := json.Unmarshal(body, &request); err != nil {
		return c.JSON(http.StatusBadRequest,
			common.ConstructResponse[*fareRS.FlightIntegratorSearchResponse](shared.FAILURE_RESPONSE_CODE, shared.FAILURE_RESPONSE_CODE, []string{"unable to unmarshal request : " + err.Error()}, nil))
	}

	logData := commonUtil.GetDataToLogFromSearch(request)
	ctx = commonUtil.ConvertAndMapToSDCContext(ctx, logData)

	ctrl.deps.GetLogger(ctx).WithField("request", request).Info("search request")
	response, err := ctrl.searchProcessorSvc.SearchProcessor(ctx, mr, request, request.TripTypes[0])
	if err != nil {
		ctrl.deps.GetLogger(ctx).Error(util.LogRest(mr, ctxOperation, shared.RESPONSE, err.Error()))
		return c.JSON(http.StatusOK, common.ConstructResponse[*fareRS.FlightIntegratorSearchResponse](GetSearchResponseCodeMapper(err), err.Error(), []string{err.Error()}, nil))
	}

	ctrl.deps.GetLogger(ctx).WithField("response", response).Info("search response")

	return c.JSON(http.StatusOK, common.ConstructResponse(shared.SUCCESS_RESPONSE_CODE, shared.SUCCESS_RESPONSE_CODE, nil, response))
}

// @Summary		Search endpoint
// @Description	Endpoint for search process
//
// @Param			storeId		header	string							true	"Authentication header"	default(TIKETCOM)
// @Param			channelId	header	string							true	"Authentication header"	default(WEB)
// @Param			requestId	header	string							true	"Authentication header"	default(321321321)
// @Param			username	header	string							true	"Authentication header"	default(guest)
// @Param			serviceId	header	string							true	"Authentication header"	default(GATEWAY)
// @Param			resellerId	header	string							false	"Authentication header"
// @Param			identity	header	string							false	"Authentication header"
//
// @Param			message		body	fareRQ.IntegratorFareRequest	true	"Request Body"
//
// @Tags			prime-booking
// @Accept			json
// @Produce		json
// @Success		200	{object}	fareRS.FlightIntegratorSearchResponse
// @Failure		400
// @Router			/v2/search [post]
func (ctrl *SearchControllerImpl) SearchV2(c echo.Context) error {
	ctxOperation := "search"
	var ctx = c.Request().Context()

	mr := common.MandatoryRequest{}
	mr = mr.BindMandatory(c)

	var request fareRQ.IntegratorFareRequest
	body, _ := io.ReadAll(c.Request().Body)
	if err := json.Unmarshal(body, &request); err != nil {
		return c.JSON(http.StatusBadRequest,
			common.ConstructResponse[*fareRS.FlightIntegratorSearchResponse](shared.FAILURE_RESPONSE_CODE, shared.FAILURE_RESPONSE_CODE, []string{"unable to unmarshal request : " + err.Error()}, nil))
	}

	logData := commonUtil.GetDataToLogFromSearch(request)
	ctx = commonUtil.ConvertAndMapToSDCContext(ctx, logData)

	ctrl.deps.GetLogger(ctx).WithField("request", request).Info("search v2 request")
	response, err := ctrl.searchProcessorSvc.SearchAndPublishProcessor(ctx, mr, request, request.TripTypes[0])
	if err != nil {
		ctrl.deps.GetLogger(ctx).Error(util.LogRest(mr, ctxOperation, shared.RESPONSE, err.Error()))
		return c.JSON(http.StatusOK, common.ConstructResponse[*fareRS.FlightIntegratorSearchResponse](GetSearchResponseCodeMapper(err), err.Error(), []string{err.Error()}, nil))
	}

	ctrl.deps.GetLogger(ctx).WithField("response", response).Info("search v2 response")

	return c.JSON(http.StatusOK, common.ConstructResponse(shared.SUCCESS_RESPONSE_CODE, shared.SUCCESS_RESPONSE_CODE, nil, response))
}
