package service

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"{{MODULE_NAME}}/internal/repository"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/test_var"
	"{{MODULE_NAME}}/internal/shared/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/log/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	CredentialMasterDataIdString     = "65097c3297a071ee3a5fe554"
	CredentialMasterDataId, _        = primitive.ObjectIDFromHex(CredentialMasterDataIdString)
	credentialMasterDataService      CredentialMasterDataServiceItf
	repo                             = new(repository.MockCredentialMasterDataInterface)
	CreateCredentialMasterDataString = repository.CredentialMasterData{
		DefaultValues:  []string{"string"},
		InputFieldType: "string",
		Key:            "string",
		Label:          "string",
		IsReadOnly:     true,
		IsRequired:     true,
	}
	UpdateCredentialMasterDataString = repository.CredentialMasterData{
		DefaultValues:  []string{"string"},
		InputFieldType: "string",
		Key:            "string",
		Label:          "edit",
		IsReadOnly:     true,
		IsRequired:     true,
	}
	UpdateCredentialMasterDataStringKey = repository.CredentialMasterData{
		DefaultValues:  []string{"string"},
		InputFieldType: "string",
		Key:            "edit",
		Label:          "edit",
		IsReadOnly:     true,
		IsRequired:     true,
	}
	CredentialMasterDataString = repository.CredentialMasterData{
		ID:             CredentialMasterDataId.Hex(),
		DefaultValues:  []string{"string"},
		InputFieldType: "string",
		Key:            "string",
		Label:          "string",
		IsReadOnly:     true,
		IsRequired:     true,
	}
	CredentialMasterDataRepositoryString = repository.CredentialMasterDataRepository{
		ID:             CredentialMasterDataId,
		DefaultValues:  []string{"string"},
		InputFieldType: "string",
		Key:            "string",
		Label:          "string",
		IsReadOnly:     true,
		IsRequired:     true,
	}
	CredentialMasterDataRepositoryStringEdit = repository.CredentialMasterDataRepository{
		ID:             CredentialMasterDataId,
		DefaultValues:  []string{"string"},
		InputFieldType: "string",
		Key:            "string",
		Label:          "edit",
		IsReadOnly:     true,
		IsRequired:     true,
	}
	CredentialMasterDataRepositoryStringEditKey = repository.CredentialMasterDataRepository{
		ID:             CredentialMasterDataId,
		DefaultValues:  []string{"string"},
		InputFieldType: "string",
		Key:            "edit",
		Label:          "edit",
		IsReadOnly:     true,
		IsRequired:     true,
	}
	CredentialMasterDataRepositoryMap = map[string]repository.CredentialMasterData{
		CredentialMasterDataString.Key: CredentialMasterDataString,
	}
	CredentialMasterDataRepositoryList = []repository.CredentialMasterData{
		CredentialMasterDataString,
	}
)

func init() {
	logger, err := logrus.DefaultLog()
	if err != nil {
		panic(err)
	}
	credentialMasterDataService, err = NewCredentialMasterDataService(
		deps.Deps{
			Logger: logger,
		},
		repo,
	)
	if err != nil {
		panic(err)
	}
}

func TestCredentialMasterDataService_Create(t *testing.T) {
	tests := []struct {
		name   string
		req    repository.CredentialMasterData
		stubs  func()
		assert func(*repository.CredentialMasterDataRepository, error)
	}{
		{
			name: "1. create success",
			req:  CreateCredentialMasterDataString,
			stubs: func() {
				repo.EXPECT().FindByKey(mock.Anything, CreateCredentialMasterDataString.Key).Return(nil, shared.ErrDataIsExist).Once()
				repo.EXPECT().Create(mock.Anything, CreateCredentialMasterDataString, test_var.MandatoryRequest).Return(&CredentialMasterDataRepositoryString, nil).Once()
			},
			assert: func(res *repository.CredentialMasterDataRepository, err error) {
				assert.Equal(t, &CredentialMasterDataRepositoryString, res)
				assert.Nil(t, err)
				repo.AssertExpectations(t)
			},
		}, {
			name: "2. key is exist success",
			req:  CreateCredentialMasterDataString,
			stubs: func() {
				repo.EXPECT().FindByKey(mock.Anything, CreateCredentialMasterDataString.Key).Return(nil, nil).Once()
			},
			assert: func(res *repository.CredentialMasterDataRepository, err error) {
				assert.Equal(t, shared.ErrDuplicateData, err)
				assert.Nil(t, res)
				repo.AssertExpectations(t)
			},
		}, {
			name: "3. create failed",
			req:  CreateCredentialMasterDataString,
			stubs: func() {
				repo.EXPECT().FindByKey(mock.Anything, CreateCredentialMasterDataString.Key).Return(nil, shared.ErrDataIsExist).Once()
				repo.EXPECT().Create(mock.Anything, CreateCredentialMasterDataString, test_var.MandatoryRequest).Return(nil, errors.New("ERROR")).Once()
			},
			assert: func(res *repository.CredentialMasterDataRepository, err error) {
				assert.Equal(t, errors.New("ERROR"), err)
				assert.Nil(t, res)
				repo.AssertExpectations(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stubs()
			res, err := credentialMasterDataService.Create(context.Background(), test_var.MandatoryRequest, tt.req)
			tt.assert(res, err)
		})
	}
}

func TestCredentialMasterDataService_FindAllMap(t *testing.T) {
	tests := []struct {
		name   string
		stubs  func()
		assert func(map[string]repository.CredentialMasterData, error)
	}{
		{
			name: "1. find all map success",
			stubs: func() {
				repo.EXPECT().FindAll(mock.Anything).Return([]repository.CredentialMasterDataRepository{
					CredentialMasterDataRepositoryString,
				}, nil).Once()
			},
			assert: func(res map[string]repository.CredentialMasterData, err error) {
				assert.Nil(t, err)
				if !assert.True(t, reflect.DeepEqual(CredentialMasterDataRepositoryMap, res)) {
					t.Errorf("error with actual = %v, want %v", util.ObjToJson(res), util.ObjToJson(CredentialMasterDataRepositoryMap))
				}
				repo.AssertExpectations(t)
			},
		}, {
			name: "2. find all map empty data",
			stubs: func() {
				repo.EXPECT().FindAll(mock.Anything).Return([]repository.CredentialMasterDataRepository{}, shared.ErrDataNotExist).Once()
			},
			assert: func(res map[string]repository.CredentialMasterData, err error) {
				assert.Equal(t, shared.ErrDataNotExist, err)
				repo.AssertExpectations(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stubs()
			res, err := credentialMasterDataService.FindAllMap(context.Background(), test_var.MandatoryRequest)
			tt.assert(res, err)
		})
	}
}

func TestCredentialMasterDataService_FindAllList(t *testing.T) {
	tests := []struct {
		name   string
		stubs  func()
		assert func([]repository.CredentialMasterData, error)
	}{
		{
			name: "1. find all map success",
			stubs: func() {
				repo.EXPECT().FindAll(mock.Anything).Return([]repository.CredentialMasterDataRepository{
					CredentialMasterDataRepositoryString,
				}, nil).Once()
			},
			assert: func(res []repository.CredentialMasterData, err error) {
				assert.Nil(t, err)

				if !assert.True(t, reflect.DeepEqual(CredentialMasterDataRepositoryList, res)) {
					t.Errorf("error with actual = %v, want %v", util.ObjToJson(res), util.ObjToJson(CredentialMasterDataRepositoryList))
				}
				repo.AssertExpectations(t)
			},
		}, {
			name: "2. find all map empty data",
			stubs: func() {
				repo.EXPECT().FindAll(mock.Anything).Return([]repository.CredentialMasterDataRepository{}, shared.ErrDataNotExist).Once()
			},
			assert: func(res []repository.CredentialMasterData, err error) {
				assert.Nil(t, res)
				assert.Equal(t, shared.ErrDataNotExist, err)
				repo.AssertExpectations(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stubs()
			res, err := credentialMasterDataService.FindAllList(context.Background(), test_var.MandatoryRequest)
			tt.assert(res, err)
		})
	}
}

func TestCredentialMasterDataService_Update(t *testing.T) {
	tests := []struct {
		name   string
		id     string
		req    repository.CredentialMasterData
		stubs  func()
		assert func(*repository.CredentialMasterDataRepository, error)
	}{
		{
			name: "1. update success without key change",
			id:   CredentialMasterDataIdString,
			req:  UpdateCredentialMasterDataString,
			stubs: func() {
				repo.EXPECT().FindByID(mock.Anything, CredentialMasterDataId).Return(&CredentialMasterDataRepositoryString, nil).Once()
				repo.EXPECT().UpdateByID(mock.Anything, CredentialMasterDataId, UpdateCredentialMasterDataString, test_var.MandatoryRequest).Return(&CredentialMasterDataRepositoryStringEdit, nil).Once()
			},
			assert: func(res *repository.CredentialMasterDataRepository, err error) {
				assert.Nil(t, err)
				assert.Equal(t, &CredentialMasterDataRepositoryStringEdit, res)
				repo.AssertExpectations(t)
			},
		}, {
			name: "2. update success with key change",
			id:   CredentialMasterDataIdString,
			req:  UpdateCredentialMasterDataStringKey,
			stubs: func() {
				repo.EXPECT().FindByID(mock.Anything, CredentialMasterDataId).Return(&CredentialMasterDataRepositoryString, nil).Once()
				repo.EXPECT().FindByKey(mock.Anything, UpdateCredentialMasterDataStringKey.Key).Return(nil, shared.ErrDataIsExist).Once()
				repo.EXPECT().UpdateByID(mock.Anything, CredentialMasterDataId, UpdateCredentialMasterDataStringKey, test_var.MandatoryRequest).Return(&CredentialMasterDataRepositoryStringEditKey, nil).Once()
			},
			assert: func(res *repository.CredentialMasterDataRepository, err error) {
				assert.Nil(t, err)
				assert.Equal(t, &CredentialMasterDataRepositoryStringEditKey, res)
				repo.AssertExpectations(t)
			},
		}, {
			name: "3. update fail id fault",
			id:   "asd",
			req:  UpdateCredentialMasterDataStringKey,
			stubs: func() {
			},
			assert: func(res *repository.CredentialMasterDataRepository, err error) {
				assert.Nil(t, res)
				assert.Equal(t, errors.New("INVALID ID"), err)
				repo.AssertExpectations(t)
			},
		}, {
			name: "4. update fail data not exist",
			id:   CredentialMasterDataIdString,
			req:  UpdateCredentialMasterDataStringKey,
			stubs: func() {
				repo.EXPECT().FindByID(mock.Anything, CredentialMasterDataId).Return(nil, shared.ErrDataIsExist).Once()
			},
			assert: func(res *repository.CredentialMasterDataRepository, err error) {
				assert.Nil(t, res)
				assert.Equal(t, shared.ErrDataNotExist, err)
				repo.AssertExpectations(t)
			},
		}, {
			name: "5. update fail duplicate key",
			id:   CredentialMasterDataIdString,
			req:  UpdateCredentialMasterDataStringKey,
			stubs: func() {
				repo.EXPECT().FindByID(mock.Anything, CredentialMasterDataId).Return(&CredentialMasterDataRepositoryString, nil).Once()
				repo.EXPECT().FindByKey(mock.Anything, UpdateCredentialMasterDataStringKey.Key).Return(nil, nil).Once()
			},
			assert: func(res *repository.CredentialMasterDataRepository, err error) {
				assert.Nil(t, res)
				assert.Equal(t, shared.ErrDuplicateData, err)
				repo.AssertExpectations(t)
			},
		}, {
			name: "6. update fail",
			id:   CredentialMasterDataIdString,
			req:  UpdateCredentialMasterDataStringKey,
			stubs: func() {
				repo.EXPECT().FindByID(mock.Anything, CredentialMasterDataId).Return(&CredentialMasterDataRepositoryString, nil).Once()
				repo.EXPECT().FindByKey(mock.Anything, UpdateCredentialMasterDataStringKey.Key).Return(nil, shared.ErrDataNotExist).Once()
				repo.EXPECT().UpdateByID(mock.Anything, CredentialMasterDataId, UpdateCredentialMasterDataStringKey, test_var.MandatoryRequest).Return(nil, errors.New("ERROR")).Once()
			},
			assert: func(res *repository.CredentialMasterDataRepository, err error) {
				assert.Nil(t, res)
				assert.Equal(t, errors.New("ERROR"), err)
				repo.AssertExpectations(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stubs()
			res, err := credentialMasterDataService.UpdateByID(context.Background(), test_var.MandatoryRequest, tt.id, tt.req)
			tt.assert(res, err)
		})
	}
}

func TestCredentialMasterDataService_Delete(t *testing.T) {
	tests := []struct {
		name   string
		id     string
		stubs  func()
		assert func(*repository.CredentialMasterDataRepository, error)
	}{
		{
			name: "1. delete success",
			id:   CredentialMasterDataIdString,
			stubs: func() {
				repo.EXPECT().DeleteByID(mock.Anything, CredentialMasterDataId).Return(nil).Once()
			},
			assert: func(res *repository.CredentialMasterDataRepository, err error) {
				assert.Nil(t, res)
				assert.Nil(t, err)
				repo.AssertExpectations(t)
			},
		}, {
			name: "2. delete fail id fault",
			id:   "asd",
			stubs: func() {
			},
			assert: func(res *repository.CredentialMasterDataRepository, err error) {
				assert.Nil(t, res)
				assert.Equal(t, shared.ErrDataNotExist, err)
				repo.AssertExpectations(t)
			},
		}, {
			name: "3. delete fail",
			id:   CredentialMasterDataIdString,
			stubs: func() {
				repo.EXPECT().DeleteByID(mock.Anything, CredentialMasterDataId).Return(errors.New("ERROR")).Once()
			},
			assert: func(res *repository.CredentialMasterDataRepository, err error) {
				assert.Nil(t, res)
				assert.Equal(t, errors.New("ERROR"), err)
				repo.AssertExpectations(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stubs()
			res, err := credentialMasterDataService.DeleteById(context.Background(), test_var.MandatoryRequest, tt.id)
			tt.assert(res, err)
		})
	}
}
