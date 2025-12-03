package service

import (
	"context"
	"{{MODULE_NAME}}/internal/repository"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/util"

	commonModel "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
)

type (
	CredentialMasterDataServiceItf interface {
		FindAllMap(ctx context.Context, mr commonModel.MandatoryRequest) (map[string]repository.CredentialMasterData, error)
		FindAllList(ctx context.Context, mr commonModel.MandatoryRequest) ([]repository.CredentialMasterData, error)
		Create(ctx context.Context, mr commonModel.MandatoryRequest, data repository.CredentialMasterData) (*repository.CredentialMasterDataRepository, error)
		UpdateByID(ctx context.Context, mr commonModel.MandatoryRequest, id string, data repository.CredentialMasterData) (*repository.CredentialMasterDataRepository, error)
		DeleteById(ctx context.Context, mr commonModel.MandatoryRequest, id string) (*repository.CredentialMasterDataRepository, error)
	}

	CredentialMasterDataServiceImpl struct {
		deps deps.Deps
		repo repository.CredentialMasterDataInterface
	}
)

func NewCredentialMasterDataService(deps deps.Deps, repo repository.CredentialMasterDataInterface) (CredentialMasterDataServiceItf, error) {
	return &CredentialMasterDataServiceImpl{deps: deps, repo: repo}, nil
}

func (impl *CredentialMasterDataServiceImpl) FindAllMap(ctx context.Context, mr commonModel.MandatoryRequest) (response map[string]repository.CredentialMasterData, err error) {
	ctxOperation := "find_all_map_credential_master_data"
	result, err := impl.repo.FindAll(ctx)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return map[string]repository.CredentialMasterData{}, err
	}

	resultMap := make(map[string]repository.CredentialMasterData)

	for _, val := range result {
		resultMap[val.Key] = repository.CredentialMasterData{
			ID:             val.ID.Hex(),
			Key:            val.Key,
			Label:          val.Label,
			IsReadOnly:     val.IsReadOnly,
			IsRequired:     val.IsRequired,
			InputFieldType: val.InputFieldType,
			DefaultValues:  val.DefaultValues,
		}
	}

	return resultMap, nil
}
func (impl *CredentialMasterDataServiceImpl) FindAllList(ctx context.Context, mr commonModel.MandatoryRequest) (response []repository.CredentialMasterData, err error) {
	ctxOperation := "find_all_credential_master_data"
	result, err := impl.repo.FindAll(ctx)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return nil, shared.ErrDataNotExist
	}
	resultList := make([]repository.CredentialMasterData, 0)
	for _, val := range result {
		resultList = append(resultList, repository.CredentialMasterData{
			ID:             val.ID.Hex(),
			Key:            val.Key,
			Label:          val.Label,
			IsReadOnly:     val.IsReadOnly,
			IsRequired:     val.IsRequired,
			InputFieldType: val.InputFieldType,
			DefaultValues:  val.DefaultValues,
		})
	}
	return resultList, nil
}

func (impl *CredentialMasterDataServiceImpl) Create(ctx context.Context, mr commonModel.MandatoryRequest, data repository.CredentialMasterData) (response *repository.CredentialMasterDataRepository, err error) {
	ctxOperation := "create_credential_master_data"
	_, err = impl.repo.FindByKey(ctx, data.Key)
	if err == nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return nil, shared.ErrDuplicateData
	}

	result, err := impl.repo.Create(ctx, data, mr)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (impl *CredentialMasterDataServiceImpl) UpdateByID(ctx context.Context, mr commonModel.MandatoryRequest, id string, data repository.CredentialMasterData) (response *repository.CredentialMasterDataRepository, err error) {
	ctxOperation := "credential_master_data_update_by_id"
	idRequest, err := util.ObjectIDFromString(id)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return nil, err
	}

	existingData, err := impl.repo.FindByID(ctx, idRequest)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return nil, shared.ErrDataNotExist
	}

	if existingData.Key != data.Key {
		_, err := impl.repo.FindByKey(ctx, data.Key)
		if err == nil {
			impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
			return nil, shared.ErrDuplicateData
		}
	}

	result, err := impl.repo.UpdateByID(ctx, idRequest, data, mr)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return nil, err
	}

	return result, nil
}

func (impl *CredentialMasterDataServiceImpl) DeleteById(ctx context.Context, mr commonModel.MandatoryRequest, id string) (response *repository.CredentialMasterDataRepository, err error) {
	ctxOperation := "credential_master_data_delete_by_id"
	idRequest, err := util.ObjectIDFromString(id)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return nil, shared.ErrDataNotExist
	}

	err = impl.repo.DeleteByID(ctx, idRequest)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return nil, err
	}

	return nil, nil
}
