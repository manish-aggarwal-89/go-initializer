package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"{{MODULE_NAME}}/internal/repository"
	"{{MODULE_NAME}}/internal/service"
	"{{MODULE_NAME}}/internal/shared/deps"

	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/log/logrus"
	commonModel "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	StoreId      = "STORE_ID"
	ChannelId    = "CHANNEL_ID"
	RequestId    = "REQUEST_ID"
	ServiceId    = "SERVICE_ID"
	ResellerId   = "RESELLER_ID"
	Identity     = "IDENTITY"
	Username     = "USERNAME"
	MandatoryReq = commonModel.MandatoryRequest{
		StoreId:    StoreId,
		ChannelId:  ChannelId,
		RequestId:  RequestId,
		ServiceId:  ServiceId,
		ResellerId: ResellerId,
		Identity:   Identity,
		Username:   Username,
	}
	CredentialMasterDataId = "11223344"
	ReqFail                = `abc`
	JsonReq                = `{
		"defaultValues": [
		  "string"
		],
		"inputFieldType": "string",
		"isReadOnly": true,
		"isRequired": true,
		"key": "string",
		"label": "string"
	}`
	CredentialMasterData = repository.CredentialMasterData{
		DefaultValues:  []string{"string"},
		InputFieldType: "string",
		Key:            "string",
		Label:          "string",
		IsReadOnly:     true,
		IsRequired:     true,
	}
)

func TestCredentialMasterDataControllerCreate(t *testing.T) {
	logger, err := logrus.DefaultLog()
	if err != nil {
		panic(err)
	}
	dep := deps.Deps{
		Logger: logger,
	}

	// create success
	credentialMasterDataSvcMock := service.NewMockCredentialMasterDataServiceItf(t)
	credentialMasterDataSvcMock.EXPECT().Create(mock.Anything, MandatoryReq, CredentialMasterData).Return(&repository.CredentialMasterDataRepository{}, nil).Once()

	ctrl := NewCredentialMasterDataController(dep, credentialMasterDataSvcMock)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(JsonReq))
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Request().Header.Add("storeId", StoreId)
	ctx.Request().Header.Add("channelId", ChannelId)
	ctx.Request().Header.Add("requestId", RequestId)
	ctx.Request().Header.Add("serviceId", ServiceId)
	ctx.Request().Header.Add("resellerId", ResellerId)
	ctx.Request().Header.Add("identity", Identity)
	ctx.Request().Header.Add("username", Username)

	assert.NoError(t, ctrl.Create(ctx))
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
	credentialMasterDataSvcMock.AssertExpectations(t)

	// create fail
	credentialMasterDataSvcMock = service.NewMockCredentialMasterDataServiceItf(t)
	credentialMasterDataSvcMock.EXPECT().Create(mock.Anything, MandatoryReq, CredentialMasterData).Return(&repository.CredentialMasterDataRepository{}, errors.New("ERROR")).Once()

	ctrl = NewCredentialMasterDataController(dep, credentialMasterDataSvcMock)

	e = echo.New()
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(JsonReq))
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	ctx.Request().Header.Add("storeId", StoreId)
	ctx.Request().Header.Add("channelId", ChannelId)
	ctx.Request().Header.Add("requestId", RequestId)
	ctx.Request().Header.Add("serviceId", ServiceId)
	ctx.Request().Header.Add("resellerId", ResellerId)
	ctx.Request().Header.Add("identity", Identity)
	ctx.Request().Header.Add("username", Username)

	assert.NoError(t, ctrl.Create(ctx))
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
	credentialMasterDataSvcMock.AssertExpectations(t)

	// error unmarshal
	credentialMasterDataSvcMock = service.NewMockCredentialMasterDataServiceItf(t)

	ctrl = NewCredentialMasterDataController(dep, credentialMasterDataSvcMock)

	e = echo.New()
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(ReqFail))
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	ctx.Request().Header.Add("storeId", StoreId)
	ctx.Request().Header.Add("channelId", ChannelId)
	ctx.Request().Header.Add("requestId", RequestId)
	ctx.Request().Header.Add("serviceId", ServiceId)
	ctx.Request().Header.Add("resellerId", ResellerId)
	ctx.Request().Header.Add("identity", Identity)
	ctx.Request().Header.Add("username", Username)

	assert.NoError(t, ctrl.Create(ctx))
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
	credentialMasterDataSvcMock.AssertExpectations(t)
}

func TestCredentialMasterDataControllerFindAllMap(t *testing.T) {
	logger, err := logrus.DefaultLog()
	if err != nil {
		panic(err)
	}
	dep := deps.Deps{
		Logger: logger,
	}

	// find all map success
	credentialMasterDataSvcMock := service.NewMockCredentialMasterDataServiceItf(t)
	credentialMasterDataSvcMock.EXPECT().FindAllMap(mock.Anything, MandatoryReq).Return(map[string]repository.CredentialMasterData{}, nil).Once()

	ctrl := NewCredentialMasterDataController(dep, credentialMasterDataSvcMock)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Request().Header.Add("storeId", StoreId)
	ctx.Request().Header.Add("channelId", ChannelId)
	ctx.Request().Header.Add("requestId", RequestId)
	ctx.Request().Header.Add("serviceId", ServiceId)
	ctx.Request().Header.Add("resellerId", ResellerId)
	ctx.Request().Header.Add("identity", Identity)
	ctx.Request().Header.Add("username", Username)

	assert.NoError(t, ctrl.FindAllMap(ctx))
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
	credentialMasterDataSvcMock.AssertExpectations(t)

	// find all map fail
	credentialMasterDataSvcMock = service.NewMockCredentialMasterDataServiceItf(t)
	credentialMasterDataSvcMock.EXPECT().FindAllMap(mock.Anything, MandatoryReq).Return(map[string]repository.CredentialMasterData{}, errors.New("ERROR")).Once()

	ctrl = NewCredentialMasterDataController(dep, credentialMasterDataSvcMock)

	e = echo.New()
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	ctx.Request().Header.Add("storeId", StoreId)
	ctx.Request().Header.Add("channelId", ChannelId)
	ctx.Request().Header.Add("requestId", RequestId)
	ctx.Request().Header.Add("serviceId", ServiceId)
	ctx.Request().Header.Add("resellerId", ResellerId)
	ctx.Request().Header.Add("identity", Identity)
	ctx.Request().Header.Add("username", Username)

	assert.NoError(t, ctrl.FindAllMap(ctx))
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
	credentialMasterDataSvcMock.AssertExpectations(t)
}

func TestCredentialMasterDataControllerFindAllList(t *testing.T) {
	logger, err := logrus.DefaultLog()
	if err != nil {
		panic(err)
	}
	dep := deps.Deps{
		Logger: logger,
	}

	// find all map success
	credentialMasterDataSvcMock := service.NewMockCredentialMasterDataServiceItf(t)
	credentialMasterDataSvcMock.EXPECT().FindAllList(mock.Anything, MandatoryReq).Return([]repository.CredentialMasterData{}, nil).Once()

	ctrl := NewCredentialMasterDataController(dep, credentialMasterDataSvcMock)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Request().Header.Add("storeId", StoreId)
	ctx.Request().Header.Add("channelId", ChannelId)
	ctx.Request().Header.Add("requestId", RequestId)
	ctx.Request().Header.Add("serviceId", ServiceId)
	ctx.Request().Header.Add("resellerId", ResellerId)
	ctx.Request().Header.Add("identity", Identity)
	ctx.Request().Header.Add("username", Username)

	assert.NoError(t, ctrl.FindAllList(ctx))
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
	credentialMasterDataSvcMock.AssertExpectations(t)

	// find all map fail
	credentialMasterDataSvcMock = service.NewMockCredentialMasterDataServiceItf(t)
	credentialMasterDataSvcMock.EXPECT().FindAllList(mock.Anything, MandatoryReq).Return([]repository.CredentialMasterData{}, errors.New("ERROR")).Once()

	ctrl = NewCredentialMasterDataController(dep, credentialMasterDataSvcMock)

	e = echo.New()
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	ctx.Request().Header.Add("storeId", StoreId)
	ctx.Request().Header.Add("channelId", ChannelId)
	ctx.Request().Header.Add("requestId", RequestId)
	ctx.Request().Header.Add("serviceId", ServiceId)
	ctx.Request().Header.Add("resellerId", ResellerId)
	ctx.Request().Header.Add("identity", Identity)
	ctx.Request().Header.Add("username", Username)

	assert.NoError(t, ctrl.FindAllList(ctx))
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
	credentialMasterDataSvcMock.AssertExpectations(t)
}

func TestCredentialMasterDataControllerUpdateById(t *testing.T) {
	logger, err := logrus.DefaultLog()
	if err != nil {
		panic(err)
	}
	dep := deps.Deps{
		Logger: logger,
	}

	// create success
	credentialMasterDataSvcMock := service.NewMockCredentialMasterDataServiceItf(t)
	credentialMasterDataSvcMock.EXPECT().UpdateByID(mock.Anything, MandatoryReq, CredentialMasterDataId, CredentialMasterData).Return(&repository.CredentialMasterDataRepository{}, nil).Once()

	ctrl := NewCredentialMasterDataController(dep, credentialMasterDataSvcMock)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(JsonReq))
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(CredentialMasterDataId)
	ctx.Request().Header.Add("storeId", StoreId)
	ctx.Request().Header.Add("channelId", ChannelId)
	ctx.Request().Header.Add("requestId", RequestId)
	ctx.Request().Header.Add("serviceId", ServiceId)
	ctx.Request().Header.Add("resellerId", ResellerId)
	ctx.Request().Header.Add("identity", Identity)
	ctx.Request().Header.Add("username", Username)

	assert.NoError(t, ctrl.UpdateByID(ctx))
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
	credentialMasterDataSvcMock.AssertExpectations(t)

	// create fail
	credentialMasterDataSvcMock = service.NewMockCredentialMasterDataServiceItf(t)
	credentialMasterDataSvcMock.EXPECT().UpdateByID(mock.Anything, MandatoryReq, CredentialMasterDataId, CredentialMasterData).Return(&repository.CredentialMasterDataRepository{}, errors.New("ERROR")).Once()

	ctrl = NewCredentialMasterDataController(dep, credentialMasterDataSvcMock)

	e = echo.New()
	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(JsonReq))
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	ctx.SetPath("/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(CredentialMasterDataId)
	ctx.Request().Header.Add("storeId", StoreId)
	ctx.Request().Header.Add("channelId", ChannelId)
	ctx.Request().Header.Add("requestId", RequestId)
	ctx.Request().Header.Add("serviceId", ServiceId)
	ctx.Request().Header.Add("resellerId", ResellerId)
	ctx.Request().Header.Add("identity", Identity)
	ctx.Request().Header.Add("username", Username)

	assert.NoError(t, ctrl.UpdateByID(ctx))
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
	credentialMasterDataSvcMock.AssertExpectations(t)

	// error unmarshal
	credentialMasterDataSvcMock = service.NewMockCredentialMasterDataServiceItf(t)

	ctrl = NewCredentialMasterDataController(dep, credentialMasterDataSvcMock)

	e = echo.New()
	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(ReqFail))
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	ctx.SetPath("/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(CredentialMasterDataId)
	ctx.Request().Header.Add("storeId", StoreId)
	ctx.Request().Header.Add("channelId", ChannelId)
	ctx.Request().Header.Add("requestId", RequestId)
	ctx.Request().Header.Add("serviceId", ServiceId)
	ctx.Request().Header.Add("resellerId", ResellerId)
	ctx.Request().Header.Add("identity", Identity)
	ctx.Request().Header.Add("username", Username)

	assert.NoError(t, ctrl.UpdateByID(ctx))
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
	credentialMasterDataSvcMock.AssertExpectations(t)

	// error id not provided
	credentialMasterDataSvcMock = service.NewMockCredentialMasterDataServiceItf(t)

	ctrl = NewCredentialMasterDataController(dep, credentialMasterDataSvcMock)

	e = echo.New()
	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(JsonReq))
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	ctx.Request().Header.Add("storeId", StoreId)
	ctx.Request().Header.Add("channelId", ChannelId)
	ctx.Request().Header.Add("requestId", RequestId)
	ctx.Request().Header.Add("serviceId", ServiceId)
	ctx.Request().Header.Add("resellerId", ResellerId)
	ctx.Request().Header.Add("identity", Identity)
	ctx.Request().Header.Add("username", Username)

	assert.NoError(t, ctrl.UpdateByID(ctx))
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
	credentialMasterDataSvcMock.AssertExpectations(t)
}

func TestCredentialMasterDataControllerDeleteById(t *testing.T) {
	logger, err := logrus.DefaultLog()
	if err != nil {
		panic(err)
	}
	dep := deps.Deps{
		Logger: logger,
	}

	// delete success
	credentialMasterDataSvcMock := service.NewMockCredentialMasterDataServiceItf(t)
	credentialMasterDataSvcMock.EXPECT().DeleteById(mock.Anything, MandatoryReq, CredentialMasterDataId).Return(&repository.CredentialMasterDataRepository{}, nil).Once()

	ctrl := NewCredentialMasterDataController(dep, credentialMasterDataSvcMock)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(JsonReq))
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(CredentialMasterDataId)
	ctx.Request().Header.Add("storeId", StoreId)
	ctx.Request().Header.Add("channelId", ChannelId)
	ctx.Request().Header.Add("requestId", RequestId)
	ctx.Request().Header.Add("serviceId", ServiceId)
	ctx.Request().Header.Add("resellerId", ResellerId)
	ctx.Request().Header.Add("identity", Identity)
	ctx.Request().Header.Add("username", Username)

	assert.NoError(t, ctrl.DeleteByID(ctx))
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
	credentialMasterDataSvcMock.AssertExpectations(t)

	// create fail
	credentialMasterDataSvcMock = service.NewMockCredentialMasterDataServiceItf(t)
	credentialMasterDataSvcMock.EXPECT().DeleteById(mock.Anything, MandatoryReq, CredentialMasterDataId).Return(&repository.CredentialMasterDataRepository{}, errors.New("ERROR")).Once()

	ctrl = NewCredentialMasterDataController(dep, credentialMasterDataSvcMock)

	e = echo.New()
	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(""))
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	ctx.SetPath("/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(CredentialMasterDataId)
	ctx.Request().Header.Add("storeId", StoreId)
	ctx.Request().Header.Add("channelId", ChannelId)
	ctx.Request().Header.Add("requestId", RequestId)
	ctx.Request().Header.Add("serviceId", ServiceId)
	ctx.Request().Header.Add("resellerId", ResellerId)
	ctx.Request().Header.Add("identity", Identity)
	ctx.Request().Header.Add("username", Username)

	assert.NoError(t, ctrl.DeleteByID(ctx))
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
	credentialMasterDataSvcMock.AssertExpectations(t)

	// error id not provided
	credentialMasterDataSvcMock = service.NewMockCredentialMasterDataServiceItf(t)

	ctrl = NewCredentialMasterDataController(dep, credentialMasterDataSvcMock)

	e = echo.New()
	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(""))
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	ctx.Request().Header.Add("storeId", StoreId)
	ctx.Request().Header.Add("channelId", ChannelId)
	ctx.Request().Header.Add("requestId", RequestId)
	ctx.Request().Header.Add("serviceId", ServiceId)
	ctx.Request().Header.Add("resellerId", ResellerId)
	ctx.Request().Header.Add("identity", Identity)
	ctx.Request().Header.Add("username", Username)

	assert.NoError(t, ctrl.DeleteByID(ctx))
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
	credentialMasterDataSvcMock.AssertExpectations(t)
}
