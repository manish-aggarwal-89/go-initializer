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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/log/logrus"
	bookRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/book"
	bookRS "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/book"
)

func TestBookingController(t *testing.T) {
	logger, err := logrus.DefaultLog()
	if err != nil {
		panic(err)
	}

	dep := deps.Deps{
		Logger: logger,
	}

	bookingSvcMock := service.NewMockBookingServiceItf(t)
	bauBookingSvcMock := bau.NewMockAnalyticsService[bookRQ.IntegratorBookRequest](t)

	ctrl := NewBookingController(dep, bookingSvcMock, bauBookingSvcMock)

	tests := []struct {
		name    string
		stubs   func()
		assert  func(rec *httptest.ResponseRecorder)
		jsonReq string
	}{
		{
			name: "success",
			stubs: func() {
				bauBookingSvcMock.EXPECT().CreateToken(mock.Anything, mock.Anything).
					Return(&bau.AnalyticBookingToken{}).
					Once()
				bauBookingSvcMock.EXPECT().Complete(mock.Anything, mock.Anything, nil).
					Once()

				bookingSvcMock.EXPECT().Booking(mock.Anything, mock.Anything, mock.Anything).
					Return(bookRS.BookResponse{}, nil).
					Once()
			},
			assert: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, rec.Result().StatusCode, http.StatusOK)

				bookingSvcMock.AssertExpectations(t)
			},
			jsonReq: `{}`,
		},
		{
			name: "booking failed with http status code 200",
			stubs: func() {
				bauBookingSvcMock.EXPECT().CreateToken(mock.Anything, mock.Anything).
					Return(&bau.AnalyticBookingToken{}).
					Once()
				bauBookingSvcMock.EXPECT().Complete(mock.Anything, mock.Anything, errors.New("error")).
					Once()

				bookingSvcMock.EXPECT().Booking(mock.Anything, mock.Anything, mock.Anything).
					Return(bookRS.BookResponse{}, errors.New("error")).
					Once()
			},
			assert: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, rec.Result().StatusCode, http.StatusOK)

				bookingSvcMock.AssertExpectations(t)
			},
			jsonReq: `{}`,
		},
		{
			name:  "bad_request",
			stubs: func() {},
			assert: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, rec.Result().StatusCode, http.StatusBadRequest)

				bookingSvcMock.AssertExpectations(t)
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

			require.NoError(t, ctrl.Booking(c))
			tt.assert(rec)
		})
	}
}
