package controller

import (
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/log/logrus"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"{{MODULE_NAME}}/internal/service"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	cancelRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/cancel-book"
	cancelRS "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/cancel-book"
)

func TestCancelBookingController(t *testing.T) {
	logger, err := logrus.DefaultLog()
	if err != nil {
		panic(err)
	}

	dep := deps.Deps{
		Logger: logger,
	}

	cancelBookingSvcMock := service.NewMockCancelBookingService(t)
	ctrl := NewCancelBookingController(dep, cancelBookingSvcMock)

	tests := []struct {
		name    string
		stubs   func()
		assert  func(rec *httptest.ResponseRecorder)
		jsonReq string
	}{
		{
			name: "success",
			stubs: func() {
				cancelBookingSvcMock.EXPECT().CancelBooking(mock.Anything, mock.Anything, cancelRQ.IntegratorCancelBookRequest{
					BookingCode: "ABC123",
				}).
					Return(cancelRS.CancelBookResponse{}, nil).
					Once()
			},
			assert: func(rec *httptest.ResponseRecorder) {
				cancelBookingSvcMock.AssertExpectations(t)

				assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
			},
			jsonReq: `{"bookingCode": "ABC123"}`,
		}, {
			name:  "failed marshal",
			stubs: func() {},
			assert: func(rec *httptest.ResponseRecorder) {
				cancelBookingSvcMock.AssertExpectations(t)

				assert.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
			},
			jsonReq: `"bookingCode": "ABC123"}`,
		}, {
			name: "cancel-booking error",
			stubs: func() {
				cancelBookingSvcMock.EXPECT().CancelBooking(mock.Anything, mock.Anything, cancelRQ.IntegratorCancelBookRequest{
					BookingCode: "ABC123",
				}).
					Return(cancelRS.CancelBookResponse{}, shared.ErrCancelBookingFailed).
					Once()
			},
			assert: func(rec *httptest.ResponseRecorder) {
				cancelBookingSvcMock.AssertExpectations(t)
				assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
			},
			jsonReq: `{"bookingCode": "ABC123"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stubs()

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.jsonReq))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			assert.NoError(t, ctrl.CancelBook(c))
			tt.assert(rec)
		})
	}
}
