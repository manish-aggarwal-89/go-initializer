package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"{{MODULE_NAME}}/internal/repository"
	"{{MODULE_NAME}}/internal/service"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/util"

	commonModel "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
)

type CredentialMasterDataController struct {
	deps    deps.Deps
	service service.CredentialMasterDataServiceItf
}

func NewCredentialMasterDataController(deps deps.Deps,
	service service.CredentialMasterDataServiceItf) *CredentialMasterDataController {
	return &CredentialMasterDataController{
		deps:    deps,
		service: service,
	}
}

// @Summary		Create credential master data endpoint
// @Description	endpoint for creting the credential master data
//
// @Param			storeId		header	string							true	"Authentication header"	default(TIKETCOM)
// @Param			channelId	header	string							true	"Authentication header"	default(WEB)
// @Param			requestId	header	string							true	"Authentication header"	default(321321321)
// @Param			username	header	string							true	"Authentication header"	default(guest)
// @Param			serviceId	header	string							true	"Authentication header"	default(GATEWAY)
// @Param			resellerId	header	string							false	"Authentication header"
// @Param			identity	header	string							false	"Authentication header"
//
// @Param			message		body	repository.CredentialMasterData	true	"Request Body"
//
// @Tags			credential-master-data
// @Accept			json
// @Produce		json
// @Success		200
// @Failure		400
// @Router			/{{BASE_PATH}}/credential-master-data [post]
func (impl *CredentialMasterDataController) Create(c echo.Context) error {
	ctxOperation := "credential_master_data_create"
	mr := util.GetMandatoryRequest(c)

	var request repository.CredentialMasterData
	body, _ := io.ReadAll(c.Request().Body)

	if err := json.Unmarshal(body, &request); err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	res, err := impl.service.Create(c.Request().Context(), mr, request)
	if err != nil {
		errRes := commonModel.ConstructResponse(shared.FAILED_INSERT_DATA_RESPONSE_CODE, err.Error(), []string{err.Error()}, res)
		impl.deps.GetLogger(c.Request().Context()).Error(util.LogRest(mr, ctxOperation, shared.RESPONSE, errRes))
		return c.JSON(http.StatusOK, errRes)
	}

	result := commonModel.ConstructResponse(shared.SUCCESS_RESPONSE_CODE, shared.SUCCESS_RESPONSE_CODE, nil, res)
	return c.JSON(http.StatusOK, result)
}

// @Summary		Endpoint to find all credential master data map
// @Description	Find all credential master data map
//
// @Param			storeId		header	string	true	"Authentication header"	default(TIKETCOM)
// @Param			channelId	header	string	true	"Authentication header"	default(WEB)
// @Param			requestId	header	string	true	"Authentication header"	default(321321321)
// @Param			username	header	string	true	"Authentication header"	default(guest)
// @Param			serviceId	header	string	true	"Authentication header"	default(GATEWAY)
// @Param			resellerId	header	string	false	"Authentication header"
// @Param			identity	header	string	false	"Authentication header"
//
// @Tags			credential-master-data
// @Accept			json
// @Produce		json
// @Success		200
// @Failure		400
// @Router			/{{BASE_PATH}}/credential-master-data [get]
func (impl *CredentialMasterDataController) FindAllMap(c echo.Context) error {
	ctxOperation := "credential_master_data_find_all_map"
	mr := util.GetMandatoryRequest(c)

	res, err := impl.service.FindAllMap(c.Request().Context(), mr)
	if err != nil {
		errRes := commonModel.ConstructResponse(shared.DATA_NOT_EXIST_RESPONSE_CODE, err.Error(), []string{err.Error()}, res)
		impl.deps.GetLogger(c.Request().Context()).Error(util.LogRest(mr, ctxOperation, shared.RESPONSE, errRes))
		return c.JSON(http.StatusOK, errRes)
	}

	result := commonModel.ConstructResponse(shared.SUCCESS_RESPONSE_CODE, shared.SUCCESS_RESPONSE_CODE, nil, res)
	return c.JSON(http.StatusOK, result)
}

// @Summary		Endpoint to find all credential master data list
// @Description	Find all credential master data list
//
// @Param			storeId		header	string	true	"Authentication header"	default(TIKETCOM)
// @Param			channelId	header	string	true	"Authentication header"	default(WEB)
// @Param			requestId	header	string	true	"Authentication header"	default(321321321)
// @Param			username	header	string	true	"Authentication header"	default(guest)
// @Param			serviceId	header	string	true	"Authentication header"	default(GATEWAY)
// @Param			resellerId	header	string	false	"Authentication header"
// @Param			identity	header	string	false	"Authentication header"
//
// @Tags			credential-master-data
// @Accept			json
// @Produce		json
// @Success		200
// @Failure		400
// @Router			/{{BASE_PATH}}/credential-master-data/list [get]
func (impl *CredentialMasterDataController) FindAllList(c echo.Context) error {
	ctxOperation := "credential_master_data_find_all_list"
	mr := util.GetMandatoryRequest(c)

	res, err := impl.service.FindAllList(c.Request().Context(), mr)
	if err != nil {
		errRes := commonModel.ConstructResponse(shared.DATA_NOT_EXIST_RESPONSE_CODE, err.Error(), []string{err.Error()}, res)
		impl.deps.GetLogger(c.Request().Context()).Error(util.LogRest(mr, ctxOperation, shared.RESPONSE, errRes))
		return c.JSON(http.StatusOK, errRes)
	}

	result := commonModel.ConstructResponse(shared.SUCCESS_RESPONSE_CODE, shared.SUCCESS_RESPONSE_CODE, nil, res)
	return c.JSON(http.StatusOK, result)
}

// @Summary		Endpoint to update credential master data by id
// @Description	Update credential master data by id
//
// @Param			id			path	string							true	"ID credential master data"
// @Param			storeId		header	string							true	"Authentication header"	default(TIKETCOM)
// @Param			channelId	header	string							true	"Authentication header"	default(WEB)
// @Param			requestId	header	string							true	"Authentication header"	default(321321321)
// @Param			username	header	string							true	"Authentication header"	default(guest)
// @Param			serviceId	header	string							true	"Authentication header"	default(GATEWAY)
// @Param			resellerId	header	string							false	"Authentication header"
// @Param			identity	header	string							false	"Authentication header"
// @Param			message		body	repository.CredentialMasterData	true	"Request Body"
//
// @Tags			credential-master-data
// @Accept			json
// @Produce		json
// @Success		200
// @Failure		400
// @Router			/{{BASE_PATH}}/credential-master-data/{id} [put]
func (impl *CredentialMasterDataController) UpdateByID(c echo.Context) error {

	ctxOperation := "credential_master_update_by_id"
	mr := util.GetMandatoryRequest(c)

	var request repository.CredentialMasterData
	body, _ := io.ReadAll(c.Request().Body)
	paramId := c.Param("id")

	if paramId == "" {
		return c.JSON(http.StatusBadRequest, "id must provided")
	}

	if err := json.Unmarshal(body, &request); err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	res, err := impl.service.UpdateByID(c.Request().Context(), mr, paramId, request)
	if err != nil {
		errRes := commonModel.ConstructResponse(shared.FAILED_UPDATE_DATA_RESPONSE_CODE, err.Error(), []string{err.Error()}, res)
		impl.deps.GetLogger(c.Request().Context()).Error(util.LogRest(mr, ctxOperation, shared.RESPONSE, errRes))
		return c.JSON(http.StatusOK, errRes)
	}

	result := commonModel.ConstructResponse(shared.SUCCESS_RESPONSE_CODE, shared.SUCCESS_RESPONSE_CODE, nil, res)
	return c.JSON(http.StatusOK, result)
}

// @Summary		Endpoint to delete credential master data by id
// @Description	Delete credential master data by id
//
// @Param			id			path	string	true	"ID Credential Master Data"
// @Param			storeId		header	string	true	"Authentication header"	default(TIKETCOM)
// @Param			channelId	header	string	true	"Authentication header"	default(WEB)
// @Param			requestId	header	string	true	"Authentication header"	default(321321321)
// @Param			username	header	string	true	"Authentication header"	default(guest)
// @Param			serviceId	header	string	true	"Authentication header"	default(GATEWAY)
// @Param			resellerId	header	string	false	"Authentication header"
// @Param			identity	header	string	false	"Authentication header"
//
// @Tags			credential-master-data
// @Accept			json
// @Produce		json
// @Success		200
// @Failure		400
// @Router			/{{BASE_PATH}}/credential-master-data/{id} [delete]
func (impl *CredentialMasterDataController) DeleteByID(c echo.Context) error {

	ctxOperation := "credential_master_delete_by_id"
	mr := util.GetMandatoryRequest(c)

	paramId := c.Param("id")
	if paramId == "" {
		return c.JSON(http.StatusBadRequest, "id must provided")
	}

	res, err := impl.service.DeleteById(c.Request().Context(), mr, paramId)
	if err != nil {
		errRes := commonModel.ConstructResponse(shared.FAILED_DELETED_DATA_RESPONSE_CODE, err.Error(), []string{err.Error()}, res)
		impl.deps.GetLogger(c.Request().Context()).Error(util.LogRest(mr, ctxOperation, shared.RESPONSE, errRes))
		return c.JSON(http.StatusOK, errRes)
	}

	result := commonModel.ConstructResponse(shared.SUCCESS_RESPONSE_CODE, shared.SUCCESS_RESPONSE_CODE, nil, res)
	return c.JSON(http.StatusOK, result)
}
