package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"{{MODULE_NAME}}/internal/predefine"
	"{{MODULE_NAME}}/internal/repository"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/entity"
	"{{MODULE_NAME}}/internal/shared/util"

	commonModel "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
)

type (
	SystemParameterServiceItf interface {
		FindAll(ctx context.Context, mr commonModel.MandatoryRequest) ([]entity.SystemParameterRepository, error)
		FindAllPaginated(ctx context.Context, mr commonModel.MandatoryRequest, limit int64, page int64) (SystemParameterPage, error)
		FindByID(ctx context.Context, mr commonModel.MandatoryRequest, id string) (*entity.SystemParameterRepository, error)
		FindByVariable(ctx context.Context, mr commonModel.MandatoryRequest, variable string) (*entity.SystemParameterRepository, error)
		FindByVariableByCache(ctx context.Context, mr commonModel.MandatoryRequest, variable string) (entity.SystemParameterRepository, error)
		Create(ctx context.Context, mr commonModel.MandatoryRequest, data entity.SystemParameter) (*entity.SystemParameterRepository, error)
		UpdateByID(ctx context.Context, mr commonModel.MandatoryRequest, id string, data entity.SystemParameter) (*entity.SystemParameterRepository, error)
		DeleteById(ctx context.Context, mr commonModel.MandatoryRequest, id string) (*entity.SystemParameterRepository, error)
	}

	SystemParameterImpl struct {
		deps      deps.Deps
		predefine predefine.Predefined
		repo      repository.SystemParameterInterface
	}

	SystemParameterPage struct {
		Content      []entity.SystemParameterRepository `json:"content"`
		TotalPage    int64
		TotalElement int64
	}
)

func NewSystemParameterServiceImpl(deps deps.Deps, repo repository.SystemParameterInterface, predefine predefine.Predefined) (SystemParameterServiceItf, error) {
	return &SystemParameterImpl{deps: deps, repo: repo, predefine: predefine}, nil
}

func (impl *SystemParameterImpl) FindAll(ctx context.Context, mr commonModel.MandatoryRequest) ([]entity.SystemParameterRepository, error) {
	result := make([]entity.SystemParameterRepository, 0)

	predefineMap := impl.predefine.GetAll(ctx)
	for _, value := range predefineMap {
		data, ok := value.(entity.SystemParameterRepository)
		if !ok {
			continue
		}

		result = append(result, data)
	}

	if len(result) > 0 {
		return result, nil
	}

	result, err := impl.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (impl *SystemParameterImpl) FindAllPaginated(ctx context.Context, mr commonModel.MandatoryRequest, limit int64, page int64) (SystemParameterPage, error) {
	ctxOperation := "system_parameter_find_all"

	result, err := impl.repo.FindAllPaginate(ctx, limit, page)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return SystemParameterPage{}, shared.ErrDataNotExist
	}
	count, err := impl.repo.Count(ctx)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return SystemParameterPage{}, errors.New("COUNT DOCUMENT ERROR")
	}
	totalPage := count / limit
	if count%limit > 0 {
		totalPage++
	}

	return SystemParameterPage{
		Content:      result,
		TotalPage:    totalPage,
		TotalElement: count,
	}, nil
}

func (impl *SystemParameterImpl) FindByID(ctx context.Context, mr commonModel.MandatoryRequest, id string) (*entity.SystemParameterRepository, error) {
	ctxOperation := "system_parameter_find_by_id"

	idRequest, err := util.ObjectIDFromString(id)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return nil, err
	}

	result, err := impl.repo.FindByID(ctx, idRequest)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return nil, err
	}

	return result, nil
}
func (impl *SystemParameterImpl) FindByVariable(ctx context.Context, mr commonModel.MandatoryRequest, variable string) (*entity.SystemParameterRepository, error) {
	ctxOperation := "system_parameter_find_variable"
	result, err := impl.repo.FindByVariable(ctx, variable)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return &entity.SystemParameterRepository{}, err
	}
	return result, nil
}

func (impl *SystemParameterImpl) FindByVariableByCache(ctx context.Context, mr commonModel.MandatoryRequest, variable string) (entity.SystemParameterRepository, error) {
	ctxOperation := "system_parameter_find_by_cache"
	cache, err := impl.predefine.GetDataByID(ctx, variable)
	if err != nil {
		result, err := impl.repo.FindByVariable(ctx, variable)
		if err != nil {
			err = shared.ErrDataNotExist
			impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, fmt.Sprintf("find variable %s got err: %s", variable, err.Error())))
			return entity.SystemParameterRepository{}, err
		}
		impl.predefine.SetData(ctx, variable, *result, time.Duration(impl.deps.Config.CacheConfig.DefaultTtlDuration))

		return *result, nil
	}
	data, _ := cache.(entity.SystemParameterRepository)
	return data, nil
}

func (impl *SystemParameterImpl) Create(ctx context.Context, mr commonModel.MandatoryRequest, data entity.SystemParameter) (*entity.SystemParameterRepository, error) {
	ctxOperation := "system_parameter_create"

	_, err := impl.repo.FindByVariable(ctx, data.Variable)
	if err == nil {
		err := errors.New("VARIABLE ALREADY EXISTS")
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return &entity.SystemParameterRepository{}, err
	}

	result, err := impl.repo.Create(ctx, data, mr)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return &entity.SystemParameterRepository{}, err
	}

	return result, nil
}

func (impl *SystemParameterImpl) UpdateByID(ctx context.Context, mr commonModel.MandatoryRequest, id string, data entity.SystemParameter) (*entity.SystemParameterRepository, error) {
	ctxOperation := "system_parameter_update_by_id"
	idRequest, err := util.ObjectIDFromString(id)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return &entity.SystemParameterRepository{}, err
	}
	result, err := impl.repo.UpdateByID(ctx, idRequest, data, mr)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return &entity.SystemParameterRepository{}, err
	}

	return result, nil
}
func (impl *SystemParameterImpl) DeleteById(ctx context.Context, mr commonModel.MandatoryRequest, id string) (*entity.SystemParameterRepository, error) {
	ctxOperation := "system_parameter_deleted_by_id"

	idRequest, err := util.ObjectIDFromString(id)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return nil, err
	}

	err = impl.repo.DeleteByID(ctx, idRequest, mr)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return nil, err
	}
	return &entity.SystemParameterRepository{}, nil

}
