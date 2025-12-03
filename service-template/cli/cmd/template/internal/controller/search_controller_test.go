package controller

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"{{MODULE_NAME}}/internal/service"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/test_var"
	"{{MODULE_NAME}}/internal/shared/test_var/search"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/log/logrus"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/constant/enum"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/fare"
)

func TestSearchController(t *testing.T) {
	logger, err := logrus.DefaultLog()
	if err != nil {
		panic(err)
	}
	searchSvcMock := service.NewMockSearchServiceItf(t)

	ctrl := NewSearchController(
		deps.Deps{
			Logger: logger,
		},
		searchSvcMock,
	)

	tests := []struct {
		name    string
		stubs   func()
		assert  func(rec *httptest.ResponseRecorder)
		jsonReq string
	}{
		{
			name: "success",
			stubs: func() {
				searchSvcMock.EXPECT().SearchWithMetric(mock.Anything, mock.Anything, mock.Anything, enum.ROUND_TRIP, mock.Anything).
					Return(fare.FlightIntegratorSearchResponse{}, nil).
					Once()
			},
			assert: func(rec *httptest.ResponseRecorder) {
				searchSvcMock.AssertExpectations(t)

				assert.Equal(t, rec.Result().StatusCode, http.StatusOK)
			},
			jsonReq: test_var.Marshal(search.INTEGRATOR_FARE_REQUEST_CGK_DPS_2A2C_RT),
		},
		{
			name: "search failed should return 200 status code with error response",
			stubs: func() {
				err = shared.ErrSearchFailed

				searchSvcMock.EXPECT().SearchWithMetric(mock.Anything, mock.Anything, mock.Anything, enum.ROUND_TRIP, mock.Anything).
					Return(fare.FlightIntegratorSearchResponse{}, err).
					Once()
			},
			assert: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, rec.Result().StatusCode, http.StatusOK)

				searchSvcMock.AssertExpectations(t)
			},
			jsonReq: test_var.Marshal(search.INTEGRATOR_FARE_REQUEST_CGK_DPS_2A2C_RT),
		},
		{
			name:  "bad request",
			stubs: func() {},
			assert: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, rec.Result().StatusCode, http.StatusBadRequest)
			},
			jsonReq: `xyz`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stubs()

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.jsonReq))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			assert.NoError(t, ctrl.Search(c))
			tt.assert(rec)
		})
	}
}

func TestSearchV2Controller(t *testing.T) {
	logger, err := logrus.DefaultLog()
	if err != nil {
		panic(err)
	}
	searchSvcMock := service.NewMockSearchServiceItf(t)

	ctrl := NewSearchController(
		deps.Deps{
			Logger: logger,
		},
		searchSvcMock,
	)

	tests := []struct {
		name    string
		stubs   func()
		assert  func(rec *httptest.ResponseRecorder)
		jsonReq string
	}{
		{
			name: "success",
			stubs: func() {
				searchSvcMock.EXPECT().SearchAndPublish(mock.Anything, mock.Anything, mock.Anything, enum.ROUND_TRIP, mock.Anything).
					Return(fare.FlightIntegratorSearchResponse{}, nil).
					Once()
			},
			assert: func(rec *httptest.ResponseRecorder) {
				searchSvcMock.AssertExpectations(t)

				assert.Equal(t, rec.Result().StatusCode, http.StatusOK)
			},
			jsonReq: test_var.Marshal(search.INTEGRATOR_FARE_REQUEST_CGK_DPS_2A2C_RT),
		},
		{
			name: "search failed should return 200 status code with error response",
			stubs: func() {
				err = shared.ErrSearchFailed

				searchSvcMock.EXPECT().SearchAndPublish(mock.Anything, mock.Anything, mock.Anything, enum.ROUND_TRIP, mock.Anything).
					Return(fare.FlightIntegratorSearchResponse{}, err).
					Once()
			},
			assert: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, rec.Result().StatusCode, http.StatusOK)

				searchSvcMock.AssertExpectations(t)
			},
			jsonReq: test_var.Marshal(search.INTEGRATOR_FARE_REQUEST_CGK_DPS_2A2C_RT),
		},
		{
			name:  "bad request",
			stubs: func() {},
			assert: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, rec.Result().StatusCode, http.StatusBadRequest)
			},
			jsonReq: `xyz`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stubs()

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.jsonReq))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			assert.NoError(t, ctrl.SearchV2(c))
			tt.assert(rec)
		})
	}
}
