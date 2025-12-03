package service

import (
	"context"
	"errors"
	"{{MODULE_NAME}}/internal/predefine"
	"{{MODULE_NAME}}/internal/repository"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/dto"
	"{{MODULE_NAME}}/internal/shared/entity"
	"{{MODULE_NAME}}/internal/shared/util"

	commonModel "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
)

type (
	CredentialServiceItf interface {
		FindByID(ctx context.Context, mr commonModel.MandatoryRequest, id string) (entity.CredentialRepository, error)
		FindWithPaginated(ctx context.Context, mr commonModel.MandatoryRequest, cf dto.CredentialFilter, page int64, size int64, sort string, sortDirect string) (CredentialPaginatedResponse, error)
		FindBySupplier(ctx context.Context, mr commonModel.MandatoryRequest, supplier string) (entity.CredentialRepository, error)
		FindBySupplierInCache(ctx context.Context, mr commonModel.MandatoryRequest, supplier string) (entity.CredentialRepository, error)
		Create(ctx context.Context, mr commonModel.MandatoryRequest, data entity.Credential) (entity.CredentialRepository, error)
		UpdateByID(ctx context.Context, mr commonModel.MandatoryRequest, id string, data entity.Credential) (entity.CredentialRepository, error)
		DeleteById(ctx context.Context, mr commonModel.MandatoryRequest, id string) (entity.CredentialRepository, error)
	}

	CredentialImpl struct {
		deps       deps.Deps
		predefine  predefine.Predefined
		repo       repository.CredentialRepositoryInterface
		serviceLog CredentialLogServiceItf
	}

	CredentialPaginatedResponse struct {
		Content       []entity.CredentialRepository `json:"content"`
		TotalPages    int64                         `json:"totalPages"`
		TotalElements int64                         `json:"totalElements"`
		Size          int64                         `json:"size"`
	}
)

func NewCredentialService(deps deps.Deps, repo repository.CredentialRepositoryInterface, predefine predefine.Predefined, serviceLog CredentialLogServiceItf) (*CredentialImpl, error) {
	return &CredentialImpl{deps: deps, repo: repo, predefine: predefine, serviceLog: serviceLog}, nil
}

// FindByID implements CredentialService.
func (impl *CredentialImpl) FindByID(ctx context.Context, mr commonModel.MandatoryRequest, id string) (entity.CredentialRepository, error) {
	ctxOperation := "credential_find_by_id"

	idRequest, err := util.ObjectIDFromString(id)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return entity.CredentialRepository{}, shared.ErrDataNotExist
	}

	result, err := impl.repo.FindByID(ctx, idRequest)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return entity.CredentialRepository{}, shared.ErrDataNotExist
	}

	return *result, nil
}

// FindWithPaginated implements CredentialService.
func (impl *CredentialImpl) FindWithPaginated(ctx context.Context, mr commonModel.MandatoryRequest, cf dto.CredentialFilter, page int64, size int64, sort string, sortDirect string) (CredentialPaginatedResponse, error) {
	ctxOperation := "credential_find_with_paginated"

	result, err := impl.repo.FindAllPaginate(ctx, cf, page, size, sort, sortDirect)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return CredentialPaginatedResponse{}, shared.ErrDataNotExist
	}
	count, err := impl.repo.Count(ctx)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return CredentialPaginatedResponse{}, shared.ErrCountDocumentError
	}

	return CredentialPaginatedResponse{
		Content:       result,
		TotalPages:    getTotalPage(count, size),
		TotalElements: count,
		Size:          size,
	}, nil
}

// FindBySupplier implements CredentialService.
func (impl *CredentialImpl) FindBySupplier(ctx context.Context, mr commonModel.MandatoryRequest, supplier string) (entity.CredentialRepository, error) {
	ctxOperation := "credential_find_by_supplier"

	result, err := impl.repo.FindBySupplierAndIsStaging(ctx, supplier, impl.deps.Config.IsStaging)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return entity.CredentialRepository{}, shared.ErrDataNotExist
	}

	return *result, nil
}

// FindBySupplierByCache implements CredentialService.
func (impl *CredentialImpl) FindBySupplierInCache(ctx context.Context, mr commonModel.MandatoryRequest, supplier string) (entity.CredentialRepository, error) {
	ctxOperation := "credential_find_by_supplier_by_cache"

	key := constructKey(supplier)
	cache, err := impl.predefine.GetDataByID(ctx, key)
	if err != nil {
		result, err := impl.repo.FindBySupplier(ctx, supplier)
		if err != nil {
			impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err.Error()))
			return entity.CredentialRepository{}, err
		}
		if err := impl.predefine.SetData(ctx, key, *result, impl.deps.Config.CacheConfig.DefaultTtlDuration); err != nil {
			impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err.Error()))
		}
		return *result, nil
	}
	data, _ := cache.(entity.CredentialRepository)
	return data, nil
}

// Create implements CredentialService.
func (impl *CredentialImpl) Create(ctx context.Context, mr commonModel.MandatoryRequest, data entity.Credential) (entity.CredentialRepository, error) {
	ctxOperation := "credential_create"

	_, err := impl.repo.FindBySupplierAndIsStaging(ctx, data.Supplier, data.IsStaging)
	if err == nil {
		err := errors.New("SUPPLIER - DUPLICATE DATA")
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return entity.CredentialRepository{}, shared.ErrDuplicateData
	}

	result, err := impl.repo.Create(ctx, data, mr)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return entity.CredentialRepository{}, err
	}

	_, errlog := impl.serviceLog.Create(ctx, mr, result)
	if errlog != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, errlog))
	}

	return *result, nil
}

// UpdateByID implements CredentialService.
func (impl *CredentialImpl) UpdateByID(ctx context.Context, mr commonModel.MandatoryRequest, id string, data entity.Credential) (entity.CredentialRepository, error) {
	ctxOperation := "credential_update_by_id"

	idRequest, err := util.ObjectIDFromString(id)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return entity.CredentialRepository{}, shared.ErrDataNotExist
	}

	existingData, err := impl.repo.FindBySupplierAndIsStaging(ctx, data.Supplier, data.IsStaging)
	if err == nil && existingData.ID != idRequest {
		err := errors.New("SUPPLIER - Duplicate Data")
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return entity.CredentialRepository{}, shared.ErrDuplicateData
	}

	result, err := impl.repo.UpdateByID(ctx, idRequest, data, mr)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return entity.CredentialRepository{}, err
	}

	_, errlog := impl.serviceLog.Create(ctx, mr, result)
	if errlog != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, errlog))
	}

	return *result, nil
}

// DeleteById implements CredentialService.
func (impl *CredentialImpl) DeleteById(ctx context.Context, mr commonModel.MandatoryRequest, id string) (entity.CredentialRepository, error) {
	ctxOperation := "credential_deleted_by_id"

	idRequest, err := util.ObjectIDFromString(id)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return entity.CredentialRepository{}, shared.ErrDataIsExist
	}

	result, err := impl.repo.DeleteByID(ctx, idRequest, mr)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return entity.CredentialRepository{}, err
	}

	_, errlog := impl.serviceLog.Create(ctx, mr, result)
	if errlog != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, errlog))
	}
	return *result, nil
}

func getTotalPage(count int64, size int64) int64 {
	if size < 1 {
		return 1
	}
	totalPage := count / size
	if count%size > 0 {
		totalPage++
	}
	return totalPage
}

func constructKey(supplier string) string {
	return shared.PREDEFINED_CREDENTIAL_CACHE_KEY_PREFIX + supplier
}
