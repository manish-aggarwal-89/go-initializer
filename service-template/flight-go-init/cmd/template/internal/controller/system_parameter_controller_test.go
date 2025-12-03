package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"{{MODULE_NAME}}/internal/service"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/entity"
	"{{MODULE_NAME}}/internal/shared/test_var"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/log/logrus"
)

func TestSystemParameterController_Create(t *testing.T) {
	logger, err := logrus.DefaultLog()
	if err != nil {
		panic(err)
	}

	svc := service.NewMockSystemParameterServiceItf(t)

	controller, _ := NewSystemParameterController(
		deps.Deps{
			Logger: logger,
		},
		svc,
	)

	tests := []struct {
		name    string
		stubs   func()
		assert  func(echo.Context)
		jsonReq string
	}{
		{
			name: "controller system parameter create success",
			stubs: func() {
				svc.EXPECT().Create(mock.Anything, mock.Anything, mock.Anything).
					Return(&entity.SystemParameterRepository{}, nil).
					Once()
			},
			assert: func(ctx echo.Context) {
				assert.Equal(t, ctx.Response().Status, http.StatusOK)
				svc.AssertExpectations(t)
			},
			jsonReq: test_var.SYSTEM_PARAMETER_JSON,
		},
		{
			name: "controller system parameter create failure",
			stubs: func() {
				svc.EXPECT().Create(mock.Anything, mock.Anything, mock.Anything).
					Return(&entity.SystemParameterRepository{}, errors.New("test error")).
					Once()
			},
			assert: func(ctx echo.Context) {
				assert.Equal(t, ctx.Response().Status, http.StatusBadRequest)
				svc.AssertExpectations(t)
			},
			jsonReq: test_var.SYSTEM_PARAMETER_JSON,
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

func TestSystemParameterController_UpdateById(t *testing.T) {
	logger, err := logrus.DefaultLog()
	if err != nil {
		panic(err)
	}

	svc := service.NewMockSystemParameterServiceItf(t)

	controller, _ := NewSystemParameterController(
		deps.Deps{
			Logger: logger,
		},
		svc,
	)

	tests := []struct {
		name    string
		stubs   func()
		assert  func(echo.Context)
		jsonReq string
	}{
		{
			name: "controller system parameter update by id success",
			stubs: func() {
				svc.EXPECT().UpdateByID(mock.Anything, mock.Anything, "", mock.Anything).
					Return(&entity.SystemParameterRepository{}, nil).
					Once()
			},
			assert: func(ctx echo.Context) {
				assert.Equal(t, ctx.Response().Status, http.StatusOK)
				svc.AssertExpectations(t)
			},
			jsonReq: test_var.SYSTEM_PARAMETER_JSON,
		},
		{
			name: "controller system parameter update by id failure",
			stubs: func() {
				svc.EXPECT().UpdateByID(mock.Anything, mock.Anything, "", mock.Anything).
					Return(&entity.SystemParameterRepository{}, errors.New("test error")).
					Once()
			},
			assert: func(ctx echo.Context) {
				assert.Equal(t, ctx.Response().Status, http.StatusBadRequest)
				svc.AssertExpectations(t)
			},
			jsonReq: test_var.SYSTEM_PARAMETER_JSON,
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

func TestSystemParameterController_DeleteById(t *testing.T) {
	logger, err := logrus.DefaultLog()
	if err != nil {
		panic(err)
	}

	svc := service.NewMockSystemParameterServiceItf(t)

	controller, _ := NewSystemParameterController(
		deps.Deps{
			Logger: logger,
		},
		svc,
	)

	tests := []struct {
		name   string
		stubs  func()
		assert func(echo.Context)
	}{
		{
			name: "controller system parameter delete by id success",
			stubs: func() {
				svc.EXPECT().DeleteById(mock.Anything, mock.Anything, "").
					Return(&entity.SystemParameterRepository{}, nil).
					Once()
			},
			assert: func(ctx echo.Context) {
				assert.Equal(t, ctx.Response().Status, http.StatusOK)
				svc.AssertExpectations(t)
			},
		},
		{
			name: "controller system parameter delete by id failure",
			stubs: func() {
				svc.EXPECT().DeleteById(mock.Anything, mock.Anything, "").
					Return(&entity.SystemParameterRepository{}, errors.New("test error")).
					Once()
			},
			assert: func(ctx echo.Context) {
				assert.Equal(t, ctx.Response().Status, http.StatusBadRequest)
				svc.AssertExpectations(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.stubs != nil {
				tt.stubs()
			}

			ctx := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			rec := httptest.NewRecorder()
			c := ctx.NewContext(req, rec)

			assert.NoError(t, controller.DeleteByID(c))
			tt.assert(c)
		})
	}
}

func TestSystemParameterController_FindById(t *testing.T) {
	logger, err := logrus.DefaultLog()
	if err != nil {
		panic(err)
	}

	svc := service.NewMockSystemParameterServiceItf(t)

	controller, _ := NewSystemParameterController(
		deps.Deps{
			Logger: logger,
		},
		svc,
	)

	tests := []struct {
		name   string
		stubs  func()
		assert func(echo.Context)
	}{
		{
			name: "controller system parameter find by id success",
			stubs: func() {
				svc.EXPECT().FindByID(mock.Anything, mock.Anything, "").
					Return(&entity.SystemParameterRepository{}, nil).
					Once()
			},
			assert: func(ctx echo.Context) {
				assert.Equal(t, ctx.Response().Status, http.StatusOK)
				svc.AssertExpectations(t)
			},
		},
		{
			name: "controller system parameter find by id failure",
			stubs: func() {
				svc.EXPECT().FindByID(mock.Anything, mock.Anything, "").
					Return(&entity.SystemParameterRepository{}, errors.New("test error")).
					Once()
			},
			assert: func(ctx echo.Context) {
				assert.Equal(t, ctx.Response().Status, http.StatusBadRequest)
				svc.AssertExpectations(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.stubs != nil {
				tt.stubs()
			}

			ctx := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := ctx.NewContext(req, rec)

			assert.NoError(t, controller.FindByID(c))
			tt.assert(c)
		})
	}
}
