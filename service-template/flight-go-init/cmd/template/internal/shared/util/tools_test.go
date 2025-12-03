package util

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	commonModel "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
)

func TestObjToJson(t *testing.T) {
	type objData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type args struct {
		data interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1. Positive",
			args: args{
				objData{
					Username: "admin",
					Password: "123456",
				},
			},
			want: "{\"username\":\"admin\",\"password\":\"123456\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ObjToJson(tt.args.data); got != tt.want {
				t.Errorf("ObjToJson() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMandatoryRequest(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("requestId", "12345678")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name string
		args args
		want commonModel.MandatoryRequest
	}{
		{
			name: "1. mandatory request",
			args: args{
				c: c,
			},
			want: commonModel.MandatoryRequest{
				RequestId: "12345678",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMandatoryRequest(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMandatoryRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObjectIDFromString(t *testing.T) {
	dataPrimitive, _ := primitive.ObjectIDFromHex("637c35a849ca3b7e655b2698")

	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    primitive.ObjectID
		wantErr bool
	}{
		{
			name: "1. using correct id mongo",
			args: args{
				s: "637c35a849ca3b7e655b2698",
			},
			want:    dataPrimitive,
			wantErr: false,
		},
		{
			name: "2. using negative case",
			args: args{
				s: "asasasas",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ObjectIDFromString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ObjectIDFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ObjectIDFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}
