package controller

import (
	"encoding/json"
	"io"
	"strconv"
	"{{MODULE_NAME}}/internal/service"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/entity"
	"github.com/labstack/echo/v4"
	"{{MODULE_NAME}}/internal/shared/util"

	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
)

type (
	SystemParameterController struct {
		deps    deps.Deps
		service service.SystemParameterServiceItf
	}
)

func NewSystemParameterController(deps deps.Deps, service service.SystemParameterServiceItf) (*SystemParameterController, error) {
	return &SystemParameterController{deps: deps, service: service}, nil
}

// FindAll @Summary      Endpoint to find all system parameter pagination
//
//	@Description	Find all system parameter
//
//	@Param			storeId		header	string	true	"Authentication header"	default(TIKETCOM)
//	@Param			channelId	header	string	true	"Authentication header"	default(WEB)
//	@Param			requestId	header	string	true	"Authentication header"	default(321321321)
//	@Param			username	header	string	true	"Authentication header"	default(guest)
//	@Param			serviceId	header	string	true	"Authentication header"	default(GATEWAY)
//	@Param			resellerId	header	string	false	"Authentication header"
//	@Param			identity	header	string	false	"Authentication header"
//
//	@Param			limit		query	int		true	"Limit"	default(10)
//	@Param			page		query	int		true	"Page"	default(0)
//
//	@Tags			system_parameter
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		400
//	@Router			/{{BASE_PATH}}/system-parameter [get]
func (impl *SystemParameterController) FindAll(c echo.Context) error {
	ctxOperation := "system_parameter_find_all"

	mr := util.GetMandatoryRequest(c)
	ctx := c.Request().Context()
	var limit int64 = 10
	var page int64 = 0
	if c.QueryParams().Has("limit") {
		limit, _ = strconv.ParseInt(c.QueryParams().Get("limit"), 10, 64)
	}
	if c.QueryParams().Has("page") {
		page, _ = strconv.ParseInt(c.QueryParams().Get("page"), 10, 64)
	}
	result, err := impl.service.FindAllPaginated(ctx, mr, limit, page)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogRest(mr, ctxOperation, shared.RESPONSE, result))
		return c.String(400, err.Error())
	}

	response := common.ConstructResponse(shared.SUCCESS_RESPONSE_CODE, shared.SUCCESS_RESPONSE_CODE, nil, result)
	impl.deps.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.RESPONSE, response))

	return c.JSON(200, response)
}

// FindByID @Summary      Endpoint to find system parameter by id
//
//	@Description	Find system parameter by id
//
//	@Param			id			path	string	true	"ID system parameter"
//	@Param			storeId		header	string	true	"Authentication header"	default(TIKETCOM)
//	@Param			channelId	header	string	true	"Authentication header"	default(WEB)
//	@Param			requestId	header	string	true	"Authentication header"	default(321321321)
//	@Param			username	header	string	true	"Authentication header"	default(guest)
//	@Param			serviceId	header	string	true	"Authentication header"	default(GATEWAY)
//	@Param			resellerId	header	string	false	"Authentication header"
//	@Param			identity	header	string	false	"Authentication header"
//
//	@Tags			system_parameter
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.SystemParameterRepository
//	@Failure		400
//	@Router			/{{BASE_PATH}}/system-parameter/{id} [get]
func (impl *SystemParameterController) FindByID(c echo.Context) error {
	ctxOperation := "system_parameter_find_by_id"

	mr := util.GetMandatoryRequest(c)
	ctx := c.Request().Context()

	id := c.Param("id")
	impl.deps.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.REQUEST, id))
	result, err := impl.service.FindByID(ctx, mr, id)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogRest(mr, ctxOperation, shared.RESPONSE, result))
		return c.JSON(400, result)
	}

	impl.deps.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.RESPONSE, result))
	return c.JSON(200, result)
}

// FindByVariable @Summary      Endpoint to find system parameter by id
//
//	@Description	Find system parameter by id
//
//	@Param			variable	path	string	true	"ID system parameter"
//	@Param			storeId		header	string	true	"Authentication header"	default(TIKETCOM)
//	@Param			channelId	header	string	true	"Authentication header"	default(WEB)
//	@Param			requestId	header	string	true	"Authentication header"	default(321321321)
//	@Param			username	header	string	true	"Authentication header"	default(guest)
//	@Param			serviceId	header	string	true	"Authentication header"	default(GATEWAY)
//	@Param			resellerId	header	string	false	"Authentication header"
//	@Param			identity	header	string	false	"Authentication header"
//
//	@Tags			system_parameter
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.SystemParameterRepository
//	@Failure		400
//	@Router			/{{BASE_PATH}}/system-parameter/variable/{variable} [get]
func (impl *SystemParameterController) FindByVariable(c echo.Context) error {
	ctxOperation := "system_parameter_find_variable"

	mr := util.GetMandatoryRequest(c)

	variable := c.Param("variable")
	ctx := c.Request().Context()

	impl.deps.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.REQUEST, variable))
	result, err := impl.service.FindByVariableByCache(ctx, mr, variable)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogRest(mr, ctxOperation, shared.RESPONSE, result))
		return c.JSON(400, result)
	}

	response := common.ConstructResponse(shared.SUCCESS_RESPONSE_CODE, shared.SUCCESS_RESPONSE_CODE, nil, result)
	impl.deps.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.RESPONSE, response))
	return c.JSON(200, response)
}

// Create @Summary      Endpoint to create system parameter
//
//	@Description	Create system parameter
//
//	@Param			storeId		header	string					true	"Authentication header"	default(TIKETCOM)
//	@Param			channelId	header	string					true	"Authentication header"	default(WEB)
//	@Param			requestId	header	string					true	"Authentication header"	default(321321321)
//	@Param			username	header	string					true	"Authentication header"	default(guest)
//	@Param			serviceId	header	string					true	"Authentication header"	default(GATEWAY)
//	@Param			resellerId	header	string					false	"Authentication header"
//	@Param			identity	header	string					false	"Authentication header"
//	@Param			message		body	entity.SystemParameter	true	"Request Body"
//
//	@Tags			system_parameter
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		400
//	@Router			/{{BASE_PATH}}/system-parameter [post]
func (impl *SystemParameterController) Create(c echo.Context) error {
	ctxOperation := "system_parameter_create"

	mr := util.GetMandatoryRequest(c)

	var request entity.SystemParameter
	body, _ := io.ReadAll(c.Request().Body)
	if err := json.Unmarshal(body, &request); err != nil {
		return nil
	}

	ctx := c.Request().Context()

	impl.deps.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.REQUEST, request))
	result, err := impl.service.Create(ctx, mr, request)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogRest(mr, ctxOperation, shared.RESPONSE, result))
		return c.JSON(400, result)
	}

	response := common.ConstructResponse(shared.SUCCESS_RESPONSE_CODE, shared.SUCCESS_RESPONSE_CODE, nil, result)
	impl.deps.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.RESPONSE, response))
	return c.JSON(200, response)
}

// UpdateByID @Summary      Endpoint to update system parameter by id
//
//	@Description	Update system parameter by id
//
//	@Param			id			path	string					true	"ID system parameter"
//	@Param			storeId		header	string					true	"Authentication header"	default(TIKETCOM)
//	@Param			channelId	header	string					true	"Authentication header"	default(WEB)
//	@Param			requestId	header	string					true	"Authentication header"	default(321321321)
//	@Param			username	header	string					true	"Authentication header"	default(guest)
//	@Param			serviceId	header	string					true	"Authentication header"	default(GATEWAY)
//	@Param			resellerId	header	string					false	"Authentication header"
//	@Param			identity	header	string					false	"Authentication header"
//	@Param			message		body	entity.SystemParameter	true	"Request Body"
//
//	@Tags			system_parameter
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		400
//	@Router			/{{BASE_PATH}}/system-parameter/{id} [put]
func (impl *SystemParameterController) UpdateByID(c echo.Context) error {
	ctxOperation := "system_parameter_update_by_id"

	mr := util.GetMandatoryRequest(c)

	var request entity.SystemParameter
	body, _ := io.ReadAll(c.Request().Body)
	if err := json.Unmarshal(body, &request); err != nil {
		return nil
	}

	id := c.Param("id")
	ctx := c.Request().Context()

	impl.deps.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.REQUEST, request))
	result, err := impl.service.UpdateByID(ctx, mr, id, request)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogRest(mr, ctxOperation, shared.RESPONSE, result))
		return c.JSON(400, result)
	}

	response := common.ConstructResponse(shared.SUCCESS_RESPONSE_CODE, shared.SUCCESS_RESPONSE_CODE, nil, result)
	impl.deps.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.RESPONSE, response))
	return c.JSON(200, response)
}

// DeleteByID @Summary      Endpoint to delete system parameter by id
//
//	@Description	Delete system parameter by id
//
//	@Param			id			path	string	true	"ID system parameter"
//	@Param			storeId		header	string	true	"Authentication header"	default(TIKETCOM)
//	@Param			channelId	header	string	true	"Authentication header"	default(WEB)
//	@Param			requestId	header	string	true	"Authentication header"	default(321321321)
//	@Param			username	header	string	true	"Authentication header"	default(guest)
//	@Param			serviceId	header	string	true	"Authentication header"	default(GATEWAY)
//	@Param			resellerId	header	string	false	"Authentication header"
//	@Param			identity	header	string	false	"Authentication header"
//
//	@Tags			system_parameter
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		400
//	@Router			/{{BASE_PATH}}/system-parameter/{id} [delete]
func (impl *SystemParameterController) DeleteByID(c echo.Context) error {
	ctxOperation := "system_parameter_delete_by_id"

	mr := util.GetMandatoryRequest(c)

	id := c.Param("id")
	ctx := c.Request().Context()

	impl.deps.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.REQUEST, id))
	result, err := impl.service.DeleteById(ctx, mr, id)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogRest(mr, ctxOperation, shared.RESPONSE, result))
		return c.JSON(400, result)
	}

	response := common.ConstructResponse(shared.SUCCESS_RESPONSE_CODE, shared.SUCCESS_RESPONSE_CODE, nil, result)
	impl.deps.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.RESPONSE, response))
	return c.JSON(200, response)
}
