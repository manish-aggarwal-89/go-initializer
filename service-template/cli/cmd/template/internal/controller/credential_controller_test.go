package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"{{MODULE_NAME}}/internal/service"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/dto"
	"{{MODULE_NAME}}/internal/shared/entity"
	"{{MODULE_NAME}}/internal/shared/test_var"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/log/logrus"
)

func TestCredentialController_Create(t *testing.T) {
	logger, err := logrus.DefaultLog()
	if err != nil {
		panic(err)
	}

	svc := service.NewMockCredentialServiceItf(t)
	svclog := service.NewMockCredentialLogServiceItf(t)

	controller := NewCredentialController(
		deps.Deps{
			Logger: logger,
		},
		svc,
		svclog,
	)

	tests := []struct {
		name    string
		stubs   func()
		assert  func(echo.Context)
		jsonReq string
	}{
		{
			name: "controller credential create success",
			stubs: func() {
				svc.EXPECT().Create(mock.Anything, mock.Anything, mock.Anything).
					Return(entity.CredentialRepository{}, nil).
					Once()
			},
			assert: func(ctx echo.Context) {
				assert.Equal(t, ctx.Response().Status, http.StatusOK)
				svc.AssertExpectations(t)
			},
			jsonReq: test_var.CredentialJson,
		},
		{
			name: "controller credential create failure",
			stubs: func() {
				svc.EXPECT().Create(mock.Anything, mock.Anything, mock.Anything).
					Return(entity.CredentialRepository{}, errors.New("test error")).
					Once()
			},
			assert: func(ctx echo.Context) {
				assert.Equal(t, ctx.Response().Status, http.StatusOK)
				svc.AssertExpectations(t)
			},
			jsonReq: test_var.CredentialJson,
		},
		{
			name: "controller credential create empty request",
			assert: func(ctx echo.Context) {
				assert.Equal(t, ctx.Response().Status, http.StatusBadRequest)
			},
			jsonReq: `{}`,
		},
		{
			name: "controller credential create bad request",
			assert: func(ctx echo.Context) {
				assert.Equal(t, ctx.Response().Status, http.StatusBadRequest)
			},
			jsonReq: `xyz`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.stubs != nil {
				tt.stubs()
			}

			ctx := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.jsonReq))
			rec := httptest.NewRecorder()
			c := ctx.NewContext(req, rec)

			assert.NoError(t, controller.Create(c))
			tt.assert(c)
		})
	}
}

func TestCredentialController_Update(t *testing.T) {
	logger, err := logrus.DefaultLog()
	if err != nil {
		panic(err)
	}

	svc := service.NewMockCredentialServiceItf(t)
	svclog := service.NewMockCredentialLogServiceItf(t)

	controller := NewCredentialController(
		deps.Deps{
			Logger: logger,
		},
		svc,
		svclog,
	)

	tests := []struct {
		name    string
		stubs   func()
		assert  func(echo.Context)
		jsonReq string
	}{
		{
			name: "controller credential update success",
			stubs: func() {
				svc.EXPECT().UpdateByID(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(entity.CredentialRepository{}, nil).
					Once()
			},
			assert: func(ctx echo.Context) {
				assert.Equal(t, ctx.Response().Status, http.StatusOK)
				svc.AssertExpectations(t)
			},
			jsonReq: test_var.CredentialJson,
		},
		{
			name: "controller credential update failure",
			stubs: func() {
				svc.EXPECT().UpdateByID(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(entity.CredentialRepository{}, errors.New("test error")).
					Once()
			},
			assert: func(ctx echo.Context) {
				assert.Equal(t, ctx.Response().Status, http.StatusOK)
				svc.AssertExpectations(t)
			},
			jsonReq: test_var.CredentialJson,
		},
		{
			name: "controller credential update empty request",
			assert: func(ctx echo.Context) {
				assert.Equal(t, ctx.Response().Status, http.StatusBadRequest)
			},
			jsonReq: `{}`,
		},
		{
			name: "controller credential update bad request",
			assert: func(ctx echo.Context) {
				assert.Equal(t, ctx.Response().Status, http.StatusBadRequest)
			},
			jsonReq: `xyz`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.stubs != nil {
				tt.stubs()
			}

			ctx := echo.New()
			req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(tt.jsonReq))
			rec := httptest.NewRecorder()
			c := ctx.NewContext(req, rec)

			assert.NoError(t, controller.UpdateByID(c))
			tt.assert(c)
		})
	}
}

func TestCredentialController_Delete(t *testing.T) {
	logger, err := logrus.DefaultLog()
	if err != nil {
		panic(err)
	}

	svc := service.NewMockCredentialServiceItf(t)
	svclog := service.NewMockCredentialLogServiceItf(t)

	controller := NewCredentialController(
		deps.Deps{
			Logger: logger,
		},
		svc,
		svclog,
	)

	tests := []struct {
		name    string
		stubs   func()
		assert  func(echo.Context)
		jsonReq string
	}{
		{
			name: "controller credential delete success",
			stubs: func() {
				svc.EXPECT().DeleteById(mock.Anything, mock.Anything, mock.Anything).
					Return(entity.CredentialRepository{}, nil).
					Once()
			},
			assert: func(ctx echo.Context) {
				assert.Equal(t, ctx.Response().Status, http.StatusOK)
				svc.AssertExpectations(t)
			},
			jsonReq: `{}`,
		},
		{
			name: "controller credential delete failure",
			stubs: func() {
				svc.EXPECT().DeleteById(mock.Anything, mock.Anything, mock.Anything).
					Return(entity.CredentialRepository{}, errors.New("test error")).
					Once()
			},
			assert: func(ctx echo.Context) {
				assert.Equal(t, ctx.Response().Status, http.StatusOK)
				svc.AssertExpectations(t)
			},
			jsonReq: `{}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.stubs != nil {
				tt.stubs()
			}

			ctx := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/", strings.NewReader(tt.jsonReq))
			rec := httptest.NewRecorder()
			c := ctx.NewContext(req, rec)

			assert.NoError(t, controller.DeleteByID(c))
			tt.assert(c)
		})
	}
}

func TestCredentialController_FindById(t *testing.T) {
	logger, err := logrus.DefaultLog()
	if err != nil {
		panic(err)
	}

	svc := service.NewMockCredentialServiceItf(t)
	svclog := service.NewMockCredentialLogServiceItf(t)

	controller := NewCredentialController(
		deps.Deps{
			Logger: logger,
		},
		svc,
		svclog,
	)

	tests := []struct {
		name    string
		stubs   func()
		assert  func(echo.Context)
		jsonReq string
	}{
		{
			name: "controller find credential by id success",
			stubs: func() {
				svc.EXPECT().FindByID(mock.Anything, mock.Anything, mock.Anything).
					Return(entity.CredentialRepository{}, nil).
					Once()
			},
			assert: func(ctx echo.Context) {
				assert.Equal(t, ctx.Response().Status, http.StatusOK)
				svc.AssertExpectations(t)
			},
			jsonReq: `{}`,
		},
		{
			name: "controller find credential by id failure",
			stubs: func() {
				svc.EXPECT().FindByID(mock.Anything, mock.Anything, mock.Anything).
					Return(entity.CredentialRepository{}, errors.New("test error")).
					Once()
			},
			assert: func(ctx echo.Context) {
				assert.Equal(t, ctx.Response().Status, http.StatusOK)
				svc.AssertExpectations(t)
			},
			jsonReq: `{}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.stubs != nil {
				tt.stubs()
			}

			ctx := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(tt.jsonReq))
			rec := httptest.NewRecorder()
			c := ctx.NewContext(req, rec)

			assert.NoError(t, controller.FindByID(c))
			tt.assert(c)
		})
	}
}

func TestCredentialController_FindWithPaginated(t *testing.T) {
	logger, err := logrus.DefaultLog()
	if err != nil {
		panic(err)
	}

	svc := service.NewMockCredentialServiceItf(t)
	svclog := service.NewMockCredentialLogServiceItf(t)

	controller := NewCredentialController(
		deps.Deps{
			Logger: logger,
		},
		svc,
		svclog,
	)

	tests := []struct {
		name    string
		stubs   func()
		assert  func(echo.Context)
		jsonReq string
	}{
		{
			name: "controller find credential with paginated success",
			stubs: func() {
				svc.EXPECT().FindWithPaginated(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(service.CredentialPaginatedResponse{}, nil).
					Once()
			},
			assert: func(ctx echo.Context) {
				assert.Equal(t, ctx.Response().Status, http.StatusOK)
				svc.AssertExpectations(t)
			},
			jsonReq: `{}`,
		},
		{
			name: "controller find credential with paginated failure",
			stubs: func() {
				svc.EXPECT().FindWithPaginated(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(service.CredentialPaginatedResponse{}, errors.New("test error")).
					Once()
			},
			assert: func(ctx echo.Context) {
				assert.Equal(t, ctx.Response().Status, http.StatusOK)
				svc.AssertExpectations(t)
			},
			jsonReq: `{}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.stubs != nil {
				tt.stubs()
			}

			ctx := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(tt.jsonReq))
			rec := httptest.NewRecorder()
			c := ctx.NewContext(req, rec)

			assert.NoError(t, controller.FindWithPaginated(c))
			tt.assert(c)
		})
	}
}

func TestCredentialController_FindLogById(t *testing.T) {
	logger, err := logrus.DefaultLog()
	if err != nil {
		panic(err)
	}

	svc := service.NewMockCredentialServiceItf(t)
	svclog := service.NewMockCredentialLogServiceItf(t)

	controller := NewCredentialController(
		deps.Deps{
			Logger: logger,
		},
		svc,
		svclog,
	)

	tests := []struct {
		name    string
		stubs   func()
		assert  func(echo.Context)
		jsonReq string
	}{
		{
			name: "controller find log credential by id success",
			stubs: func() {
				svclog.EXPECT().FindByID(mock.Anything, mock.Anything, mock.Anything).
					Return([]dto.CredentialLogResponse{}, nil).
					Once()
			},
			assert: func(ctx echo.Context) {
				assert.Equal(t, ctx.Response().Status, http.StatusOK)
				svclog.AssertExpectations(t)
			},
			jsonReq: `{}`,
		},
		{
			name: "controller find log credential by id failure",
			stubs: func() {
				svclog.EXPECT().FindByID(mock.Anything, mock.Anything, mock.Anything).
					Return([]dto.CredentialLogResponse{}, errors.New("test error")).
					Once()
			},
			assert: func(ctx echo.Context) {
				assert.Equal(t, ctx.Response().Status, http.StatusOK)
				svclog.AssertExpectations(t)
			},
			jsonReq: `{}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.stubs != nil {
				tt.stubs()
			}

			ctx := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(tt.jsonReq))
			rec := httptest.NewRecorder()
			c := ctx.NewContext(req, rec)

			assert.NoError(t, controller.FindLogByID(c))
			tt.assert(c)
		})
	}
}

func TestCredentialController_SendAlert(t *testing.T) {

}
