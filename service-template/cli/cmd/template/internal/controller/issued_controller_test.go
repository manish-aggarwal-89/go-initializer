package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"{{MODULE_NAME}}/internal/service"
	"{{MODULE_NAME}}/internal/shared/bau"
	"{{MODULE_NAME}}/internal/shared/deps"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/log/logrus"
	issuedRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/issued"
	issuedRS "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/issued"
)

func TestIssuedController(t *testing.T) {
	logger, err := logrus.DefaultLog()
	if err != nil {
		panic(err)
	}

	dep := deps.Deps{
		Logger: logger,
	}

	issuedSvcMock := service.NewMockIssuedServiceItf(t)
	bauIssuedSvcMock := bau.NewMockAnalyticsService[issuedRQ.IntegratorIssuedRequest](t)

	ctrl := NewIssuedController(dep, issuedSvcMock, bauIssuedSvcMock)

	tests := []struct {
		name    string
		stubs   func()
		assert  func(rec *httptest.ResponseRecorder)
		jsonReq string
	}{
		{
			name: "success",
			stubs: func() {
				bauIssuedSvcMock.EXPECT().CreateToken(mock.Anything, mock.Anything).
					Return(&bau.AnalyticsIssuanceToken{}).
					Once()
				bauIssuedSvcMock.EXPECT().Complete(mock.Anything, mock.Anything, nil).
					Once()

				issuedSvcMock.EXPECT().Issued(mock.Anything, mock.Anything, mock.Anything).
					Return(issuedRS.IssuedResponse{}, nil).
					Once()
			},
			assert: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, rec.Result().StatusCode, http.StatusOK)

				issuedSvcMock.AssertExpectations(t)
			},
			jsonReq: `{}`,
		},
		{
			name: "failed should return 200 status code with error response",
			stubs: func() {
				bauIssuedSvcMock.EXPECT().CreateToken(mock.Anything, mock.Anything).
					Return(&bau.AnalyticsIssuanceToken{}).
					Once()
				bauIssuedSvcMock.EXPECT().Complete(mock.Anything, mock.Anything, errors.New("error")).
					Once()

				issuedSvcMock.EXPECT().Issued(mock.Anything, mock.Anything, mock.Anything).
					Return(issuedRS.IssuedResponse{}, errors.New("error")).
					Once()
			},
			assert: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, rec.Result().StatusCode, http.StatusOK)

				issuedSvcMock.AssertExpectations(t)
			},
			jsonReq: `{}`,
		},
		{
			name:  "bad request",
			stubs: func() {},
			assert: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, rec.Result().StatusCode, http.StatusBadRequest)

				issuedSvcMock.AssertExpectations(t)
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

			assert.NoError(t, ctrl.Issued(c))
			tt.assert(rec)
		})
	}
}
