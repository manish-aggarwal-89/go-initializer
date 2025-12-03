package service

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
	"{{MODULE_NAME}}/config"
	"{{MODULE_NAME}}/internal/predefine"
	"{{MODULE_NAME}}/internal/repository"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/entity"
	"{{MODULE_NAME}}/internal/shared/test_var"
	"{{MODULE_NAME}}/internal/shared/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/log/logrus"
)

var (
	systemParameterSvc  SystemParameterServiceItf
	systemParameterRepo = new(repository.MockSystemParameterInterface)
	sysParamPredefine   = new(predefine.MockPredefined)
)

func init() {
	logger, err := logrus.DefaultLog()
	if err != nil {
		panic(err)
	}
	systemParameterSvc, err = NewSystemParameterServiceImpl(
		deps.Deps{
			Logger: logger,
			Config: &config.Config{
				CacheConfig: config.CacheConfig{
					DefaultTtlDuration:             time.Hour,
					DefaultCleanupIntervalDuration: time.Hour,
				},
			},
		},
		systemParameterRepo,
		sysParamPredefine,
	)
	if err != nil {
		panic(err)
	}
}

func TestSystemParameterSvc_FindWithPaginated(t *testing.T) {
	tests := []struct {
		name   string
		stubs  func()
		assert func(SystemParameterPage, error)
	}{
		{
			name: "1. find all success",
			stubs: func() {
				systemParameterRepo.EXPECT().FindAllPaginate(mock.Anything, int64(1), int64(10)).Return(test_var.LIST_SYSTEM_PARAMETER, nil).Once()
				systemParameterRepo.EXPECT().Count(mock.Anything).Return(int64(1), nil).Once()
			},
			assert: func(res SystemParameterPage, err error) {
				assert.NoError(t, err)
				if !assert.True(t, reflect.DeepEqual(test_var.LIST_SYSTEM_PARAMETER, res.Content)) {
					t.Errorf("error with actual = %v, want %v", util.ObjToJson(res), util.ObjToJson(test_var.LIST_SYSTEM_PARAMETER))
				}
				systemParameterRepo.AssertExpectations(t)
			},
		},
		{
			name: "2. find all empty data",
			stubs: func() {
				systemParameterRepo.EXPECT().FindAllPaginate(mock.Anything, int64(1), int64(10)).Return([]entity.SystemParameterRepository{}, errors.New("Error")).Once()
			},
			assert: func(res SystemParameterPage, err error) {
				assert.Error(t, err, shared.ErrDataNotExist)
				systemParameterRepo.AssertExpectations(t)
			},
		},
		{
			name: "3. find all error count data",
			stubs: func() {
				systemParameterRepo.EXPECT().FindAllPaginate(mock.Anything, int64(1), int64(10)).Return(test_var.LIST_SYSTEM_PARAMETER, nil).Once()
				systemParameterRepo.EXPECT().Count(mock.Anything).Return(0, errors.New("Error")).Once()
			},
			assert: func(res SystemParameterPage, err error) {
				assert.Error(t, err, errors.New("COUNT DOCUMENT ERROR"))
				systemParameterRepo.AssertExpectations(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stubs()
			res, err := systemParameterSvc.FindAllPaginated(context.Background(), test_var.MandatoryRequest, 1, 10)
			tt.assert(res, err)
		})
	}
}

func TestSystemParameterService_FindById(t *testing.T) {
	tests := []struct {
		name   string
		req    string
		stubs  func()
		assert func(*entity.SystemParameterRepository, error)
	}{
		{
			name: "1. find by id success",
			req:  test_var.SYSTEM_PARAMETER_ID_STR,
			stubs: func() {
				systemParameterRepo.EXPECT().FindByID(mock.Anything, test_var.SYSTEM_PARAMETER_ID).Return(&test_var.SYSTEM_PARAMETER_PUBLISH_ON_BOOKING_FLAG_FALSE, nil).Once()
			},
			assert: func(res *entity.SystemParameterRepository, err error) {
				assert.Equal(t, &test_var.SYSTEM_PARAMETER_PUBLISH_ON_BOOKING_FLAG_FALSE, res)
				assert.NoError(t, err)
				repo.AssertExpectations(t)
			},
		},
		{
			name: "2. find by id failed",
			req:  test_var.SYSTEM_PARAMETER_ID_ERROR_STR,
			stubs: func() {
				systemParameterRepo.EXPECT().FindByID(mock.Anything, test_var.SYSTEM_PARAMETER_ID_ERROR).Return(nil, errors.New("ERROR")).Once()
			},
			assert: func(res *entity.SystemParameterRepository, err error) {
				assert.Nil(t, res)
				assert.Error(t, shared.ErrDataNotExist, err)
				repo.AssertExpectations(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stubs()
			res, err := systemParameterSvc.FindByID(context.Background(), test_var.MandatoryRequest, tt.req)
			tt.assert(res, err)
		})
	}
}

func TestSystemParameterService_FindByVariable(t *testing.T) {
	tests := []struct {
		name   string
		req    string
		stubs  func()
		assert func(*entity.SystemParameterRepository, error)
	}{
		{
			name: "1. find by variable success",
			req:  test_var.VARIABLE,
			stubs: func() {
				systemParameterRepo.EXPECT().FindByVariable(mock.Anything, test_var.VARIABLE).Return(&test_var.SYSTEM_PARAMETER_PUBLISH_ON_BOOKING_FLAG_FALSE, nil).Once()
			},
			assert: func(res *entity.SystemParameterRepository, err error) {
				assert.Equal(t, &test_var.SYSTEM_PARAMETER_PUBLISH_ON_BOOKING_FLAG_FALSE, res)
				assert.NoError(t, err)
				repo.AssertExpectations(t)
			},
		},
		{
			name: "2. find by id failed",
			req:  test_var.VARIABLE,
			stubs: func() {
				systemParameterRepo.EXPECT().FindByVariable(mock.Anything, test_var.VARIABLE).Return(nil, errors.New("ERROR")).Once()
			},
			assert: func(res *entity.SystemParameterRepository, err error) {
				assert.Equal(t, &entity.SystemParameterRepository{}, res)
				assert.Error(t, shared.ErrDataNotExist, err)
				repo.AssertExpectations(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stubs()
			res, err := systemParameterSvc.FindByVariable(context.Background(), test_var.MandatoryRequest, tt.req)
			tt.assert(res, err)
		})
	}
}
func TestSystemParameterService_FindByVariableCache(t *testing.T) {
	tests := []struct {
		name   string
		req    string
		stubs  func()
		assert func(entity.SystemParameterRepository, error)
	}{
		{
			name: "1. find by variable from empty cache",
			req:  test_var.VARIABLE,
			stubs: func() {
				sysParamPredefine.EXPECT().GetDataByID(mock.Anything, test_var.VARIABLE).Return(nil, shared.ErrDataIsExist).Once()
				systemParameterRepo.EXPECT().FindByVariable(mock.Anything, test_var.VARIABLE).Return(&test_var.SYSTEM_PARAMETER_PUBLISH_ON_BOOKING_FLAG_FALSE, nil).Once()
				sysParamPredefine.EXPECT().SetData(mock.Anything, test_var.VARIABLE, test_var.SYSTEM_PARAMETER_PUBLISH_ON_BOOKING_FLAG_FALSE, mock.Anything).Return(nil).Once()
			},
			assert: func(res entity.SystemParameterRepository, err error) {
				assert.Equal(t, test_var.SYSTEM_PARAMETER_PUBLISH_ON_BOOKING_FLAG_FALSE, res)
				assert.NoError(t, err)
				repo.AssertExpectations(t)
			},
		},
		{
			name: "2. find by variable from existing cache",
			req:  test_var.VARIABLE,
			stubs: func() {
				sysParamPredefine.EXPECT().GetDataByID(mock.Anything, test_var.VARIABLE).Return(test_var.SYSTEM_PARAMETER_PUBLISH_ON_BOOKING_FLAG_FALSE, nil).Once()
			},
			assert: func(res entity.SystemParameterRepository, err error) {
				assert.Equal(t, test_var.SYSTEM_PARAMETER_PUBLISH_ON_BOOKING_FLAG_FALSE, res)
				assert.NoError(t, err)
				repo.AssertExpectations(t)
			},
		},
		{
			name: "3. find by variable from cache failed",
			req:  test_var.VARIABLE,
			stubs: func() {
				sysParamPredefine.EXPECT().GetDataByID(mock.Anything, test_var.VARIABLE).Return(nil, shared.ErrDataIsExist).Once()
				systemParameterRepo.EXPECT().FindByVariable(mock.Anything, test_var.VARIABLE).Return(nil, errors.New("ERROR")).Once()
			},
			assert: func(res entity.SystemParameterRepository, err error) {
				assert.Equal(t, entity.SystemParameterRepository{}, res)
				assert.Error(t, shared.ErrDataNotExist, err)
				repo.AssertExpectations(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stubs()
			res, err := systemParameterSvc.FindByVariableByCache(context.Background(), test_var.MandatoryRequest, tt.req)
			tt.assert(res, err)
		})
	}
}

func TestSystemParameterService_Create(t *testing.T) {
	tests := []struct {
		name   string
		req    entity.SystemParameter
		stubs  func()
		assert func(*entity.SystemParameterRepository, error)
	}{
		{
			name: "1. create system parameter success",
			req:  test_var.SYSTEM_PARAMETER_REQUEST,
			stubs: func() {
				systemParameterRepo.EXPECT().FindByVariable(mock.Anything, test_var.SYSTEM_PARAMETER_REQUEST.Variable).Return(nil, shared.ErrDataNotExist).Once()
				systemParameterRepo.EXPECT().Create(mock.Anything, mock.Anything, mock.Anything).Return(&test_var.SYSTEM_PARAMETER_PUBLISH_ON_BOOKING_FLAG_FALSE, nil).Once()
			},
			assert: func(res *entity.SystemParameterRepository, err error) {
				assert.Equal(t, &test_var.SYSTEM_PARAMETER_PUBLISH_ON_BOOKING_FLAG_FALSE, res)
				assert.NoError(t, err)
				repo.AssertExpectations(t)
			},
		},
		{
			name: "2. create system parameter duplicate data",
			req:  test_var.SYSTEM_PARAMETER_REQUEST,
			stubs: func() {
				systemParameterRepo.EXPECT().FindByVariable(mock.Anything, test_var.SYSTEM_PARAMETER_REQUEST.Variable).Return(nil, nil).Once()
			},
			assert: func(res *entity.SystemParameterRepository, err error) {
				assert.Error(t, err, shared.ErrDuplicateData)
				repo.AssertExpectations(t)
			},
		}, {
			name: "3. create system parameter error",
			req:  test_var.SYSTEM_PARAMETER_REQUEST,
			stubs: func() {
				systemParameterRepo.EXPECT().FindByVariable(mock.Anything, test_var.SYSTEM_PARAMETER_REQUEST.Variable).Return(nil, shared.ErrDataNotExist).Once()
				systemParameterRepo.EXPECT().Create(mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("err")).Once()
			},
			assert: func(res *entity.SystemParameterRepository, err error) {
				assert.Error(t, err, errors.New("err"))
				repo.AssertExpectations(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stubs()
			res, err := systemParameterSvc.Create(context.Background(), test_var.MandatoryRequest, tt.req)
			tt.assert(res, err)
		})
	}
}

func TestSystemParameterService_UpdateById(t *testing.T) {
	tests := []struct {
		name   string
		id     string
		req    entity.SystemParameter
		stubs  func()
		assert func(*entity.SystemParameterRepository, error)
	}{
		{
			name: "1. update by id success",
			id:   test_var.SYSTEM_PARAMETER_ID_STR,
			req:  test_var.SYSTEM_PARAMETER_REQUEST,
			stubs: func() {
				systemParameterRepo.EXPECT().UpdateByID(mock.Anything, test_var.SYSTEM_PARAMETER_ID, test_var.SYSTEM_PARAMETER_REQUEST, mock.Anything).Return(&test_var.SYSTEM_PARAMETER_PUBLISH_ON_BOOKING_FLAG_FALSE, nil).Once()
			},
			assert: func(res *entity.SystemParameterRepository, err error) {
				assert.Equal(t, &test_var.SYSTEM_PARAMETER_PUBLISH_ON_BOOKING_FLAG_FALSE, res)
				assert.NoError(t, err)
				repo.AssertExpectations(t)
			},
		},
		{
			name: "2. update by id error",
			req:  test_var.SYSTEM_PARAMETER_REQUEST,
			stubs: func() {

				systemParameterRepo.EXPECT().UpdateByID(mock.Anything, test_var.SYSTEM_PARAMETER_ID, test_var.SYSTEM_PARAMETER_REQUEST, mock.Anything).Return(nil, errors.New("Error")).Once()
			},
			assert: func(res *entity.SystemParameterRepository, err error) {
				assert.Error(t, err, errors.New("Error"))
				repo.AssertExpectations(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stubs()
			res, err := systemParameterSvc.UpdateByID(context.Background(), test_var.MandatoryRequest, tt.id, tt.req)
			tt.assert(res, err)
		})
	}
}

func TestSystemParameterService_DeleteById(t *testing.T) {
	tests := []struct {
		name   string
		id     string
		stubs  func()
		assert func(*entity.SystemParameterRepository, error)
	}{
		{
			name: "1. delete by id success",
			id:   test_var.SYSTEM_PARAMETER_ID_STR,
			stubs: func() {
				systemParameterRepo.EXPECT().DeleteByID(mock.Anything, test_var.SYSTEM_PARAMETER_ID, mock.Anything).Return(nil).Once()
			},
			assert: func(res *entity.SystemParameterRepository, err error) {
				assert.NoError(t, err)
				repo.AssertExpectations(t)
			},
		},
		{
			name: "2. delete by id error",
			stubs: func() {

				systemParameterRepo.EXPECT().DeleteByID(mock.Anything, test_var.SYSTEM_PARAMETER_ID, mock.Anything).Return(errors.New("Error")).Once()
			},
			assert: func(res *entity.SystemParameterRepository, err error) {
				assert.Error(t, err, errors.New("Error"))
				repo.AssertExpectations(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stubs()
			res, err := systemParameterSvc.DeleteById(context.Background(), test_var.MandatoryRequest, tt.id)
			tt.assert(res, err)
		})
	}
}
