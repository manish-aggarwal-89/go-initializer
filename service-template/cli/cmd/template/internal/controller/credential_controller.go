package controller

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"{{MODULE_NAME}}/internal/service"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/dto"
	"{{MODULE_NAME}}/internal/shared/entity"
	"{{MODULE_NAME}}/internal/shared/util"

	commonModel "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
)

type (
	CredentialController struct {
		dep        deps.Deps
		service    service.CredentialServiceItf
		logService service.CredentialLogServiceItf
	}
)

func NewCredentialController(deps deps.Deps, service service.CredentialServiceItf, logService service.CredentialLogServiceItf) *CredentialController {
	return &CredentialController{dep: deps, service: service, logService: logService}
}

// @Summary		Endpoint to find credential with paginated
// @Description	Find credential with paginated
//
// @Param			storeId			header	string					true	"Authentication header"	default(TIKETCOM)
// @Param			channelId		header	string					true	"Authentication header"	default(WEB)
// @Param			requestId		header	string					true	"Authentication header"	default(321321321)
// @Param			username		header	string					true	"Authentication header"	default(guest)
// @Param			serviceId		header	string					true	"Authentication header"	default(GATEWAY)
// @Param			resellerId		header	string					false	"Authentication header"
// @Param			identity		header	string					false	"Authentication header"
// @Param			page			query	int						true	"Page"				default(0)
// @Param			size			query	int						true	"Size"				default(10)
// @Param			sort			query	dto.SortType			true	"Sort"				default(ID)
// @Param			sortDirection	query	dto.OrderType			true	"Sort Direction"	default(DESC)
// @Param			filter			query	dto.CredentialFilter	true	"Filter"
//
// @Tags			credential
// @Accept			json
// @Produce		json
// @Success		200
// @Failure		400
// @Router			/{{BASE_PATH}}/credential [get]
func (impl *CredentialController) FindWithPaginated(c echo.Context) error {
	ctxOperation := "credential_find_with_paginated"

	mr := util.GetMandatoryRequest(c)
	ctx := c.Request().Context()

	var (
		page       int64 = 0
		size       int64 = 10
		sort             = ""
		sortDirect       = ""
	)

	if c.QueryParams().Has("page") {
		page, _ = strconv.ParseInt(c.QueryParams().Get("page"), 10, 64)
	}
	if c.QueryParams().Has("size") {
		size, _ = strconv.ParseInt(c.QueryParams().Get("size"), 10, 64)
	}
	if c.QueryParams().Has("sort") {
		sort = c.QueryParams().Get("sort")
	}
	if c.QueryParams().Has("sortDirection") {
		sortDirect = c.QueryParams().Get("sortDirection")
	}

	filter := getFilter(c)

	impl.dep.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.REQUEST, filter))
	result, err := impl.service.FindWithPaginated(ctx, mr, filter, page, size, sort, sortDirect)
	if err != nil {
		errRes := commonModel.ConstructResponse(shared.DATA_NOT_EXIST_RESPONSE_CODE, err.Error(), []string{err.Error()}, result)
		impl.dep.GetLogger(ctx).Error(util.LogRest(mr, ctxOperation, shared.RESPONSE, errRes))
		return c.JSON(http.StatusOK, errRes)
	}

	res := commonModel.ConstructResponse(shared.SUCCESS_RESPONSE_CODE, shared.SUCCESS_RESPONSE_CODE, nil, result)
	impl.dep.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.RESPONSE, res))
	return c.JSON(http.StatusOK, res)
}

// @Summary		Endpoint to find credential by id
// @Description	Find credential by id
//
// @Param			id			path	string	true	"ID credential"
// @Param			storeId		header	string	true	"Authentication header"	default(TIKETCOM)
// @Param			channelId	header	string	true	"Authentication header"	default(WEB)
// @Param			requestId	header	string	true	"Authentication header"	default(321321321)
// @Param			username	header	string	true	"Authentication header"	default(guest)
// @Param			serviceId	header	string	true	"Authentication header"	default(GATEWAY)
// @Param			resellerId	header	string	false	"Authentication header"
// @Param			identity	header	string	false	"Authentication header"
//
// @Tags			credential
// @Accept			json
// @Produce		json
// @Success		200
// @Failure		400
// @Router			/{{BASE_PATH}}/credential/{id} [get]
func (impl *CredentialController) FindByID(c echo.Context) error {
	ctxOperation := "credential_find_by_id"

	mr := util.GetMandatoryRequest(c)
	ctx := c.Request().Context()

	id := c.Param("id")

	impl.dep.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.REQUEST, id))
	result, err := impl.service.FindByID(ctx, mr, id)
	if err != nil {
		errRes := commonModel.ConstructResponse(shared.DATA_NOT_EXIST_RESPONSE_CODE, err.Error(), []string{err.Error()}, result)
		impl.dep.GetLogger(ctx).Error(util.LogRest(mr, ctxOperation, shared.RESPONSE, errRes))
		return c.JSON(http.StatusOK, errRes)
	}

	res := commonModel.ConstructResponse(shared.SUCCESS_RESPONSE_CODE, shared.SUCCESS_RESPONSE_CODE, nil, result)
	impl.dep.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.RESPONSE, res))
	return c.JSON(http.StatusOK, res)
}

// @Summary		Endpoint to create credential
// @Description	Create credential
//
// @Param			storeId		header	string				true	"Authentication header"	default(TIKETCOM)
// @Param			channelId	header	string				true	"Authentication header"	default(WEB)
// @Param			requestId	header	string				true	"Authentication header"	default(321321321)
// @Param			username	header	string				true	"Authentication header"	default(guest)
// @Param			serviceId	header	string				true	"Authentication header"	default(GATEWAY)
// @Param			resellerId	header	string				false	"Authentication header"
// @Param			identity	header	string				false	"Authentication header"
// @Param			message		body	entity.Credential	true	"Request Body"
//
// @Tags			credential
// @Accept			json
// @Produce		json
// @Success		200
// @Failure		400
// @Router			/{{BASE_PATH}}/credential [post]
func (impl *CredentialController) Create(c echo.Context) error {
	ctxOperation := "credential_create"

	mr := util.GetMandatoryRequest(c)
	ctx := c.Request().Context()

	validate := validator.New()
	var request entity.Credential
	body, _ := io.ReadAll(c.Request().Body)
	if err := json.Unmarshal(body, &request); err != nil {
		return c.JSON(http.StatusBadRequest, commonModel.ConstructResponse(shared.FAILURE_RESPONSE_CODE, err.Error(), []string{err.Error()}, &entity.CredentialRepository{}))
	}
	err := validate.Struct(request)
	if err != nil {
		var errs validator.ValidationErrors
		errors.As(err, &errs)
		return c.JSON(http.StatusBadRequest, commonModel.ConstructResponse(shared.FAILURE_RESPONSE_CODE, err.Error(), []string{err.Error()}, &entity.CredentialRepository{}))
	}

	impl.dep.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.REQUEST, request))
	result, err := impl.service.Create(ctx, mr, request)
	if err != nil {

		errRes := commonModel.ConstructResponse(shared.FAILED_INSERT_DATA_RESPONSE_CODE, err.Error(), []string{err.Error()}, result)
		impl.dep.GetLogger(ctx).Error(util.LogRest(mr, ctxOperation, shared.RESPONSE, errRes))
		return c.JSON(http.StatusOK, result)
	}

	res := commonModel.ConstructResponse(shared.SUCCESS_RESPONSE_CODE, shared.SUCCESS_RESPONSE_CODE, nil, result)
	impl.dep.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.RESPONSE, res))
	return c.JSON(http.StatusOK, res)
}

// @Summary		Endpoint to update credential by id
// @Description	Update credential by id
//
// @Param			id			path	string				true	"ID credential"
// @Param			storeId		header	string				true	"Authentication header"	default(TIKETCOM)
// @Param			channelId	header	string				true	"Authentication header"	default(WEB)
// @Param			requestId	header	string				true	"Authentication header"	default(321321321)
// @Param			username	header	string				true	"Authentication header"	default(guest)
// @Param			serviceId	header	string				true	"Authentication header"	default(GATEWAY)
// @Param			resellerId	header	string				false	"Authentication header"
// @Param			identity	header	string				false	"Authentication header"
// @Param			message		body	entity.Credential	true	"Request Body"
//
// @Tags			credential
// @Accept			json
// @Produce		json
// @Success		200
// @Failure		400
// @Router			/{{BASE_PATH}}/credential/{id} [put]
func (impl *CredentialController) UpdateByID(c echo.Context) error {
	ctxOperation := "credential_update_by_id"

	mr := util.GetMandatoryRequest(c)
	ctx := c.Request().Context()

	validate := validator.New()
	var request entity.Credential
	body, _ := io.ReadAll(c.Request().Body)
	if err := json.Unmarshal(body, &request); err != nil {
		return c.JSON(http.StatusBadRequest, commonModel.ConstructResponse(shared.FAILURE_RESPONSE_CODE, err.Error(), []string{err.Error()}, &entity.CredentialRepository{}))
	}
	err := validate.Struct(request)
	if err != nil {
		var errs validator.ValidationErrors
		errors.As(err, &errs)
		return c.JSON(http.StatusBadRequest, commonModel.ConstructResponse(shared.FAILURE_RESPONSE_CODE, err.Error(), []string{err.Error()}, &entity.CredentialRepository{}))
	}

	id := c.Param("id")
	impl.dep.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.REQUEST, request))
	result, err := impl.service.UpdateByID(ctx, mr, id, request)
	if err != nil {
		errRes := commonModel.ConstructResponse(shared.FAILED_UPDATE_DATA_RESPONSE_CODE, err.Error(), []string{err.Error()}, result)
		impl.dep.GetLogger(ctx).Error(util.LogRest(mr, ctxOperation, shared.RESPONSE, errRes))
		return c.JSON(http.StatusOK, errRes)
	}

	res := commonModel.ConstructResponse(shared.SUCCESS_RESPONSE_CODE, shared.SUCCESS_RESPONSE_CODE, nil, result)
	impl.dep.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.RESPONSE, res))
	return c.JSON(http.StatusOK, res)
}

// @Summary		Endpoint to delete credential by id
// @Description	Delete credential by id
//
// @Param			id			path	string	true	"ID credential"
// @Param			storeId		header	string	true	"Authentication header"	default(TIKETCOM)
// @Param			channelId	header	string	true	"Authentication header"	default(WEB)
// @Param			requestId	header	string	true	"Authentication header"	default(321321321)
// @Param			username	header	string	true	"Authentication header"	default(guest)
// @Param			serviceId	header	string	true	"Authentication header"	default(GATEWAY)
// @Param			resellerId	header	string	false	"Authentication header"
// @Param			identity	header	string	false	"Authentication header"
//
// @Tags			credential
// @Accept			json
// @Produce		json
// @Success		200
// @Failure		400
// @Router			/{{BASE_PATH}}/credential/{id} [delete]
func (impl *CredentialController) DeleteByID(c echo.Context) error {
	ctxOperation := "credential_delete_by_id"

	mr := util.GetMandatoryRequest(c)
	ctx := c.Request().Context()

	id := c.Param("id")
	impl.dep.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.REQUEST, id))
	result, err := impl.service.DeleteById(ctx, mr, id)
	if err != nil {
		errRes := commonModel.ConstructResponse(shared.FAILED_DELETED_DATA_RESPONSE_CODE, err.Error(), []string{err.Error()}, result)
		impl.dep.GetLogger(ctx).Error(util.LogRest(mr, ctxOperation, shared.RESPONSE, errRes))
		return c.JSON(http.StatusOK, errRes)
	}

	res := commonModel.ConstructResponse(shared.SUCCESS_RESPONSE_CODE, shared.SUCCESS_RESPONSE_CODE, nil, result)
	impl.dep.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.RESPONSE, res))
	return c.JSON(http.StatusOK, res)
}

// @Summary		Endpoint to find credential log by id
// @Description	Find credential log by id
//
// @Param			id			path	string	true	"ID credential"
// @Param			storeId		header	string	true	"Authentication header"	default(TIKETCOM)
// @Param			channelId	header	string	true	"Authentication header"	default(WEB)
// @Param			requestId	header	string	true	"Authentication header"	default(321321321)
// @Param			username	header	string	true	"Authentication header"	default(guest)
// @Param			serviceId	header	string	true	"Authentication header"	default(GATEWAY)
// @Param			resellerId	header	string	false	"Authentication header"
// @Param			identity	header	string	false	"Authentication header"
//
// @Tags			credential
// @Accept			json
// @Produce		json
// @Success		200
// @Failure		400
// @Router			/{{BASE_PATH}}/credential/log/{id} [get]
func (impl *CredentialController) FindLogByID(c echo.Context) error {
	ctxOperation := "credential_log_find_by_id"

	mr := util.GetMandatoryRequest(c)

	id := c.Param("id")
	ctx := c.Request().Context()
	impl.dep.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.REQUEST, id))
	result, err := impl.logService.FindByID(ctx, mr, id)
	if err != nil {
		errRes := commonModel.ConstructResponse(shared.DATA_NOT_EXIST_RESPONSE_CODE, err.Error(), []string{err.Error()}, result)
		impl.dep.GetLogger(ctx).Error(util.LogRest(mr, ctxOperation, shared.RESPONSE, errRes))
		return c.JSON(http.StatusOK, errRes)
	}
	res := commonModel.ConstructResponse(shared.SUCCESS_RESPONSE_CODE, shared.SUCCESS_RESPONSE_CODE, nil, result)
	impl.dep.GetLogger(ctx).Info(util.LogRest(mr, ctxOperation, shared.RESPONSE, res))
	return c.JSON(http.StatusOK, res)
}

// @Summary		Endpoint to send the alert
// @Description	Find credential log by id
//
// @Param			storeId		header	string	true	"Authentication header"	default(TIKETCOM)
// @Param			channelId	header	string	true	"Authentication header"	default(WEB)
// @Param			requestId	header	string	true	"Authentication header"	default(321321321)
// @Param			username	header	string	true	"Authentication header"	default(guest)
// @Param			serviceId	header	string	true	"Authentication header"	default(GATEWAY)
// @Param			resellerId	header	string	false	"Authentication header"
// @Param			identity	header	string	false	"Authentication header"
//
// @Tags			credential
// @Accept			json
// @Produce		json
// @Success		200
// @Failure		400
// @Router			/{{BASE_PATH}}/credential/alert [get]
func (impl *CredentialController) SendAlert(c echo.Context) error {
	// ctxOperation := "credential_alert"

	// mr := util.GetMandatoryRequest(c)
	// err := impl.alertService.Alert(ctx, mr)
	// if err != nil {
	// 	impl.dep.Logger.Error(util.LogRest(mr, ctxOperation, shared.RESPONSE, err))
	// 	return c.JSON(http.StatusInternalServerError, common.ResponseError("ERROR", err))
	// }

	return c.JSON(http.StatusOK, commonModel.ConstructResponse(shared.SUCCESS_RESPONSE_CODE, shared.SUCCESS_RESPONSE_CODE, nil, "SUCCESS"))
}

func getFilter(c echo.Context) dto.CredentialFilter {
	var (
		distribution = ""
		supplier     = ""
		username     = ""
	)

	if c.QueryParams().Has("distributionType") {
		distribution = c.QueryParams().Get("distributionType")
	}

	if c.QueryParams().Has("supplier") {
		supplier = c.QueryParams().Get("supplier")
	}

	if c.QueryParams().Has("username") {
		username = c.QueryParams().Get("username")
	}

	return dto.CredentialFilter{
		DistributionType: distribution,
		Supplier:         supplier,
		Username:         username,
	}
}
