package service

import (
	"context"
	"errors"
	"testing"
	"{{MODULE_NAME}}/internal/repository"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/dto"
	"{{MODULE_NAME}}/internal/shared/entity"
	"{{MODULE_NAME}}/internal/shared/test_var"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/log/logrus"
)

var (
	credentialLogSvc  CredentialLogServiceItf
	credentialLogRepo = new(repository.MockCredentialLogRepositoryInterface)
)

func init() {
	logger, err := logrus.DefaultLog()
	if err != nil {
		panic(err)
	}
	credentialLogSvc, err = NewCredentialLogService(
		deps.Deps{
			Logger: logger,
		},
		credentialLogRepo,
	)
	if err != nil {
		panic(err)
	}
}

func TestCredentialLogService_Create(t *testing.T) {
	tests := []struct {
		name   string
		req    entity.CredentialRepository
		stubs  func()
		assert func(*entity.CredentialLogRepository)
	}{
		{
			name: "1. create success",
			req:  test_var.CredentialRepo,
			stubs: func() {
				credentialLogRepo.EXPECT().Create(mock.Anything, test_var.CredentialLog, test_var.MandatoryRequest).Return(&test_var.CredentialLogRepo, nil).Once()
			},
			assert: func(res *entity.CredentialLogRepository) {
				assert.Equal(t, &test_var.CredentialLogRepo, res)
				credentialLogRepo.AssertExpectations(t)
			},
		},
		{
			name: "2. create failed",
			req:  test_var.CredentialRepo,
			stubs: func() {
				credentialLogRepo.EXPECT().Create(mock.Anything, test_var.CredentialLog, test_var.MandatoryRequest).Return(nil, errors.New("ERROR")).Once()
			},
			assert: func(res *entity.CredentialLogRepository) {
				assert.Empty(t, res)
				credentialLogRepo.AssertExpectations(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stubs()
			res, _ := credentialLogSvc.Create(context.Background(), test_var.MandatoryRequest, &tt.req)
			tt.assert(res)
		})
	}
}

func TestCredentialLogService_FindByID(t *testing.T) {
	tests := []struct {
		name   string
		req    string
		stubs  func()
		assert func([]dto.CredentialLogResponse, error)
	}{
		{
			name: "1. find by id success",
			req:  test_var.CredentialIdString,
			stubs: func() {
				credentialLogRepo.EXPECT().FindByID(mock.Anything, test_var.CredentialId).Return(test_var.CredentialLogRepos, nil).Once()
			},
			assert: func(res []dto.CredentialLogResponse, err error) {
				assert.Nil(t, err)
				assert.Equal(t, test_var.CredentialLogResps, res)
				credentialLogRepo.AssertExpectations(t)
			},
		},
		{
			name: "2. find by id  failed",
			req:  test_var.CredentialIdString,
			stubs: func() {
				credentialLogRepo.EXPECT().FindByID(mock.Anything, test_var.CredentialId).Return(nil, errors.New("ERROR")).Once()
			},
			assert: func(res []dto.CredentialLogResponse, err error) {
				assert.Equal(t, errors.New("ERROR"), err)
				assert.Equal(t, []dto.CredentialLogResponse{}, res)
				credentialLogRepo.AssertExpectations(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stubs()
			res, err := credentialLogSvc.FindByID(context.Background(), test_var.MandatoryRequest, tt.req)
			tt.assert(res, err)
		})
	}
}
