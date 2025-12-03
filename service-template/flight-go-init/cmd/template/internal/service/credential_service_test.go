package service

import (
	"context"
	"errors"
	"testing"
	"time"
	"{{MODULE_NAME}}/config"
	"{{MODULE_NAME}}/internal/predefine"
	"{{MODULE_NAME}}/internal/repository"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/dto"
	"{{MODULE_NAME}}/internal/shared/entity"
	"{{MODULE_NAME}}/internal/shared/test_var"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/log/logrus"
)

var (
	credentialSvc        CredentialServiceItf
	mockCredentialLogSvc = new(MockCredentialLogServiceItf)
	mockPredfine         = new(predefine.MockPredefined)
	credentialRepo       = new(repository.MockCredentialRepositoryInterface)
)

func init() {
	logger, err := logrus.DefaultLog()
	if err != nil {
		panic(err)
	}
	credentialSvc, err = NewCredentialService(
		deps.Deps{
			Logger: logger,
			Config: &config.Config{
				CacheConfig: config.CacheConfig{
					DefaultTtlDuration:             time.Hour,
					DefaultCleanupIntervalDuration: time.Hour,
				},
			},
		},
		credentialRepo,
		mockPredfine,
		mockCredentialLogSvc,
	)
	if err != nil {
		panic(err)
	}
}

func TestCredentialService_Create(t *testing.T) {
	tests := []struct {
		name   string
		req    entity.Credential
		stubs  func()
		assert func(entity.CredentialRepository, error)
	}{
		{
			name: "1. create success",
			req:  test_var.CREDENTIAL,
			stubs: func() {
				credentialRepo.EXPECT().FindBySupplierAndIsStaging(mock.Anything, test_var.CREDENTIAL.Supplier, mock.Anything).Return(nil, shared.ErrDataNotExist).Once()
				credentialRepo.EXPECT().Create(mock.Anything, test_var.CREDENTIAL, test_var.MandatoryRequest).Return(&test_var.CredentialRepo, nil).Once()
				mockCredentialLogSvc.EXPECT().Create(mock.Anything, test_var.MandatoryRequest, &test_var.CredentialRepo).Return(nil, nil).Once()
			},
			assert: func(res entity.CredentialRepository, err error) {
				assert.Equal(t, test_var.CredentialRepo, res)
				repo.AssertExpectations(t)
			},
		},
		{
			name: "2. key is exist then duplicate data",
			req:  test_var.CREDENTIAL,
			stubs: func() {
				credentialRepo.EXPECT().FindBySupplierAndIsStaging(mock.Anything, test_var.CREDENTIAL.Supplier, mock.Anything).Return(&entity.CredentialRepository{}, nil).Once()
			},
			assert: func(res entity.CredentialRepository, err error) {
				assert.Equal(t, err, shared.ErrDuplicateData)
				repo.AssertExpectations(t)
			},
		},
		{
			name: "3. create failed",
			req:  test_var.CREDENTIAL,
			stubs: func() {
				credentialRepo.EXPECT().FindBySupplierAndIsStaging(mock.Anything, test_var.CREDENTIAL.Supplier, mock.Anything).Return(nil, shared.ErrDataIsExist).Once()
				credentialRepo.EXPECT().Create(mock.Anything, test_var.CREDENTIAL, test_var.MandatoryRequest).Return(nil, errors.New("ERROR")).Once()
			},
			assert: func(res entity.CredentialRepository, err error) {
				assert.Equal(t, err, errors.New("ERROR"))
				repo.AssertExpectations(t)
			},
		},
		{
			name: "4. create success but create log failed",
			req:  test_var.CREDENTIAL,
			stubs: func() {
				credentialRepo.EXPECT().FindBySupplierAndIsStaging(mock.Anything, test_var.CREDENTIAL.Supplier, mock.Anything).Return(nil, shared.ErrDataIsExist).Once()
				credentialRepo.EXPECT().Create(mock.Anything, test_var.CREDENTIAL, test_var.MandatoryRequest).Return(&test_var.CredentialRepo, nil).Once()
				mockCredentialLogSvc.EXPECT().Create(mock.Anything, test_var.MandatoryRequest, &test_var.CredentialRepo).Return(nil, errors.New("ERROR")).Once()
			},
			assert: func(res entity.CredentialRepository, err error) {
				assert.Equal(t, test_var.CredentialRepo, res)
				repo.AssertExpectations(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stubs()
			res, err := credentialSvc.Create(context.Background(), test_var.MandatoryRequest, tt.req)
			tt.assert(res, err)
		})
	}
}

func TestCredentialService_Update(t *testing.T) {
	tests := []struct {
		name   string
		req    entity.Credential
		reqId  string
		stubs  func()
		assert func(entity.CredentialRepository, error)
	}{
		{
			name:  "1. update success",
			req:   test_var.CREDENTIAL,
			reqId: test_var.CredentialIdString,
			stubs: func() {
				credentialRepo.EXPECT().FindBySupplierAndIsStaging(mock.Anything, test_var.CREDENTIAL.Supplier, mock.Anything).Return(nil, shared.ErrDataIsExist).Once()
				credentialRepo.EXPECT().UpdateByID(mock.Anything, test_var.CredentialId, test_var.CREDENTIAL, test_var.MandatoryRequest).Return(&test_var.CredentialRepo, nil).Once()
				mockCredentialLogSvc.EXPECT().Create(mock.Anything, test_var.MandatoryRequest, &test_var.CredentialRepo).Return(nil, nil).Once()
			},
			assert: func(res entity.CredentialRepository, err error) {
				assert.Equal(t, test_var.CredentialRepo, res)
				repo.AssertExpectations(t)
			},
		},
		{
			name:  "2. key is exist",
			req:   test_var.CREDENTIAL,
			reqId: test_var.CredentialIdString,
			stubs: func() {
				credentialRepo.EXPECT().FindBySupplierAndIsStaging(mock.Anything, test_var.CREDENTIAL.Supplier, mock.Anything).Return(&entity.CredentialRepository{}, nil).Once()
			},
			assert: func(res entity.CredentialRepository, err error) {
				assert.Equal(t, shared.ErrDuplicateData, err)
				repo.AssertExpectations(t)
			},
		},
		{
			name:  "3. update failed",
			req:   test_var.CREDENTIAL,
			reqId: test_var.CredentialIdString,
			stubs: func() {
				credentialRepo.EXPECT().FindBySupplierAndIsStaging(mock.Anything, test_var.CREDENTIAL.Supplier, mock.Anything).Return(nil, shared.ErrDataIsExist).Once()
				credentialRepo.EXPECT().UpdateByID(mock.Anything, test_var.CredentialId, test_var.CREDENTIAL, test_var.MandatoryRequest).Return(nil, errors.New("ERROR")).Once()
			},
			assert: func(res entity.CredentialRepository, err error) {
				assert.Equal(t, errors.New("ERROR"), err)
				repo.AssertExpectations(t)
			},
		},
		{
			name:  "4. update success but create log failed",
			req:   test_var.CREDENTIAL,
			reqId: test_var.CredentialIdString,
			stubs: func() {
				credentialRepo.EXPECT().FindBySupplierAndIsStaging(mock.Anything, test_var.CREDENTIAL.Supplier, mock.Anything).Return(nil, shared.ErrDataIsExist).Once()
				credentialRepo.EXPECT().UpdateByID(mock.Anything, test_var.CredentialId, test_var.CREDENTIAL, test_var.MandatoryRequest).Return(&test_var.CredentialRepo, nil).Once()
				mockCredentialLogSvc.EXPECT().Create(mock.Anything, test_var.MandatoryRequest, &test_var.CredentialRepo).Return(nil, errors.New("ERROR")).Once()
			},

			assert: func(res entity.CredentialRepository, err error) {
				assert.Equal(t, test_var.CredentialRepo, res)
				assert.Nil(t, err)
				repo.AssertExpectations(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stubs()
			res, err := credentialSvc.UpdateByID(context.Background(), test_var.MandatoryRequest, tt.reqId, tt.req)
			tt.assert(res, err)
		})
	}
}

func TestCredentialService_Delete(t *testing.T) {
	tests := []struct {
		name   string
		req    string
		stubs  func()
		assert func(entity.CredentialRepository, error)
	}{
		{
			name: "1. delete success",
			req:  test_var.CredentialIdString,
			stubs: func() {
				credentialRepo.EXPECT().DeleteByID(mock.Anything, test_var.CredentialId, test_var.MandatoryRequest).Return(&test_var.CredentialRepo, nil).Once()
				mockCredentialLogSvc.EXPECT().Create(mock.Anything, test_var.MandatoryRequest, &test_var.CredentialRepo).Return(nil, nil).Once()
			},
			assert: func(res entity.CredentialRepository, err error) {
				assert.Nil(t, err)
				repo.AssertExpectations(t)
			},
		},
		{
			name: "2. delete failed",
			req:  test_var.CredentialIdString,
			stubs: func() {
				credentialRepo.EXPECT().DeleteByID(mock.Anything, test_var.CredentialId, test_var.MandatoryRequest).Return(nil, errors.New("ERROR")).Once()
			},
			assert: func(res entity.CredentialRepository, err error) {
				assert.Equal(t, errors.New("ERROR"), err)
				repo.AssertExpectations(t)
			},
		},
		{
			name: "3. delete success but create log failed",
			req:  test_var.CredentialIdString,
			stubs: func() {
				credentialRepo.EXPECT().DeleteByID(mock.Anything, test_var.CredentialId, test_var.MandatoryRequest).Return(&test_var.CredentialRepo, nil).Once()
				mockCredentialLogSvc.EXPECT().Create(mock.Anything, test_var.MandatoryRequest, &test_var.CredentialRepo).Return(nil, errors.New("ERROR")).Once()
			},

			assert: func(res entity.CredentialRepository, err error) {
				assert.Nil(t, err)
				repo.AssertExpectations(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stubs()
			res, err := credentialSvc.DeleteById(context.Background(), test_var.MandatoryRequest, tt.req)
			tt.assert(res, err)
		})
	}
}

func TestCredentialService_FindById(t *testing.T) {
	tests := []struct {
		name   string
		req    string
		stubs  func()
		assert func(entity.CredentialRepository, error)
	}{
		{
			name: "1. find by id success",
			req:  test_var.CredentialIdString,
			stubs: func() {
				credentialRepo.EXPECT().FindByID(mock.Anything, test_var.CredentialId).Return(&test_var.CredentialRepo, nil).Once()
			},
			assert: func(res entity.CredentialRepository, err error) {
				assert.Equal(t, test_var.CredentialRepo, res)
				assert.Nil(t, err)
				repo.AssertExpectations(t)
			},
		},
		{
			name: "2. find by id failed",
			req:  test_var.CredentialIdString,
			stubs: func() {
				credentialRepo.EXPECT().FindByID(mock.Anything, test_var.CredentialId).Return(nil, errors.New("ERROR")).Once()
			},
			assert: func(res entity.CredentialRepository, err error) {
				assert.Equal(t, shared.ErrDataNotExist, err)
				repo.AssertExpectations(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stubs()
			res, err := credentialSvc.FindByID(context.Background(), test_var.MandatoryRequest, tt.req)
			tt.assert(res, err)
		})
	}
}

func TestCredentialService_FindBySupplier(t *testing.T) {
	tests := []struct {
		name   string
		req    string
		stubs  func()
		assert func(entity.CredentialRepository, error)
	}{
		{
			name: "1. find by supplier success",
			req:  test_var.CredentialSupplierString,
			stubs: func() {
				credentialRepo.EXPECT().FindBySupplierAndIsStaging(mock.Anything,
					test_var.CredentialSupplierString, false).
					Return(&test_var.CredentialRepo, nil).Once()
			},
			assert: func(res entity.CredentialRepository, err error) {
				assert.Nil(t, err)
				assert.Equal(t, test_var.CredentialRepo, res)
				repo.AssertExpectations(t)
			},
		},
		{
			name: "2. find by supplier failed",
			req:  test_var.CredentialSupplierString,
			stubs: func() {
				credentialRepo.EXPECT().FindBySupplierAndIsStaging(mock.Anything,
					test_var.CredentialSupplierString, false).
					Return(nil, errors.New("ERROR")).Once()
			},

			assert: func(res entity.CredentialRepository, err error) {
				assert.Equal(t, shared.ErrDataNotExist, err)
				repo.AssertExpectations(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stubs()
			res, err := credentialSvc.FindBySupplier(context.Background(), test_var.MandatoryRequest, tt.req)
			tt.assert(res, err)
		})
	}
}

func TestCredentialService_FindBySupplierByCache(t *testing.T) {
	tests := []struct {
		name   string
		req    string
		stubs  func()
		assert func(entity.CredentialRepository, error)
	}{
		{
			name: "1. find by supplier from cache success",
			req:  test_var.CredentialSupplierString,
			stubs: func() {
				mockPredfine.EXPECT().GetDataByID(mock.Anything, test_var.CredentialPredefineKey).Return(test_var.CredentialRepo, nil).Once()
			},
			assert: func(res entity.CredentialRepository, err error) {
				assert.Equal(t, test_var.CredentialRepo, res)
				repo.AssertExpectations(t)
			},
		},
		{
			name: "2. find by supplier from db success",
			req:  test_var.CredentialSupplierString,
			stubs: func() {
				mockPredfine.EXPECT().GetDataByID(mock.Anything, test_var.CredentialPredefineKey).Return(nil, shared.ErrDataIsExist).Once()
				credentialRepo.EXPECT().FindBySupplier(mock.Anything, test_var.CredentialSupplierString).Return(&test_var.CredentialRepo, nil).Once()
				mockPredfine.EXPECT().SetData(mock.Anything, test_var.CredentialPredefineKey, test_var.CredentialRepo, mock.Anything).Return(nil).Once()
			},
			assert: func(res entity.CredentialRepository, err error) {
				assert.Equal(t, test_var.CredentialRepo, res)
				repo.AssertExpectations(t)
			},
		},
		{
			name: "3. find by supplier from db but set cache failed",
			req:  test_var.CredentialSupplierString,
			stubs: func() {
				mockPredfine.EXPECT().GetDataByID(mock.Anything, test_var.CredentialPredefineKey).Return(nil, shared.ErrDataIsExist).Once()
				credentialRepo.EXPECT().FindBySupplier(mock.Anything, test_var.CredentialSupplierString).Return(&test_var.CredentialRepo, nil).Once()
				mockPredfine.EXPECT().SetData(mock.Anything, test_var.CredentialPredefineKey, test_var.CredentialRepo, mock.Anything).Return(errors.New("ERROR")).Once()
			},
			assert: func(res entity.CredentialRepository, err error) {
				assert.Equal(t, test_var.CredentialRepo, res)
				repo.AssertExpectations(t)
			},
		},
		{
			name: "4. find by supplier from db failed",
			req:  test_var.CredentialSupplierString,
			stubs: func() {
				mockPredfine.EXPECT().GetDataByID(mock.Anything, test_var.CredentialPredefineKey).Return(nil, shared.ErrDataIsExist).Once()
				credentialRepo.EXPECT().FindBySupplier(mock.Anything, test_var.CredentialSupplierString).Return(nil, errors.New("ERROR")).Once()
			},
			assert: func(res entity.CredentialRepository, err error) {
				assert.Equal(t, errors.New("ERROR"), err)
				repo.AssertExpectations(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stubs()
			res, err := credentialSvc.FindBySupplierInCache(context.Background(), test_var.MandatoryRequest, tt.req)
			tt.assert(res, err)
		})
	}
}

func TestCredentialService_FindWithPaginated(t *testing.T) {
	tests := []struct {
		name   string
		req    dto.CredentialFilter
		stubs  func()
		assert func(CredentialPaginatedResponse, error)
	}{
		{
			name: "1. find with paginated success",
			req:  test_var.CredentialFilter,
			stubs: func() {
				credentialRepo.EXPECT().FindAllPaginate(mock.Anything, test_var.CredentialFilter, int64(0), int64(10), string(dto.ASC), string(dto.ID)).Return(test_var.CredentialRepos, nil).Once()
				credentialRepo.EXPECT().Count(mock.Anything).Return(int64(1), nil).Once()
			},
			assert: func(res CredentialPaginatedResponse, err error) {
				assert.Equal(t, test_var.CredentialRepos, res.Content)
				assert.Nil(t, err)
				repo.AssertExpectations(t)
			},
		},
		{
			name: "2. find with paginated failed",
			req:  test_var.CredentialFilter,
			stubs: func() {
				credentialRepo.EXPECT().FindAllPaginate(mock.Anything, test_var.CredentialFilter, int64(0), int64(10), string(dto.ASC), string(dto.ID)).Return(nil, errors.New("ERROR")).Once()
			},
			assert: func(res CredentialPaginatedResponse, err error) {
				assert.Equal(t, shared.ErrDataNotExist, err)
				repo.AssertExpectations(t)
			},
		},
		{
			name: "3. find with paginated, count document failed",
			req:  test_var.CredentialFilter,
			stubs: func() {
				credentialRepo.EXPECT().FindAllPaginate(mock.Anything, test_var.CredentialFilter, int64(0), int64(10), string(dto.ASC), string(dto.ID)).Return(test_var.CredentialRepos, nil).Once()
				credentialRepo.EXPECT().Count(mock.Anything).Return(int64(0), errors.New("ERROR")).Once()
			},

			assert: func(res CredentialPaginatedResponse, err error) {
				assert.Equal(t, shared.ErrCountDocumentError, err)
				repo.AssertExpectations(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stubs()
			res, err := credentialSvc.FindWithPaginated(context.Background(), test_var.MandatoryRequest, tt.req, int64(0), int64(10), string(dto.ASC), string(dto.ID))
			tt.assert(res, err)
		})
	}
}
