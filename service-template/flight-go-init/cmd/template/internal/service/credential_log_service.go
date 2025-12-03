package service

import (
	"context"
	"{{MODULE_NAME}}/internal/repository"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/dto"
	"{{MODULE_NAME}}/internal/shared/entity"
	"{{MODULE_NAME}}/internal/shared/util"

	commonModel "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
)

type (
	CredentialLogServiceItf interface {
		FindByID(ctx context.Context, mr commonModel.MandatoryRequest, id string) ([]dto.CredentialLogResponse, error)
		Create(ctx context.Context, mr commonModel.MandatoryRequest, data *entity.CredentialRepository) (*entity.CredentialLogRepository, error)
	}

	CredentialLogImpl struct {
		deps deps.Deps
		repo repository.CredentialLogRepositoryInterface
	}
)

func NewCredentialLogService(deps deps.Deps, repo repository.CredentialLogRepositoryInterface) (CredentialLogServiceItf, error) {
	return &CredentialLogImpl{deps: deps, repo: repo}, nil
}

// Create implements CredentialLogService.
func (impl *CredentialLogImpl) Create(ctx context.Context, mr commonModel.MandatoryRequest, data *entity.CredentialRepository) (*entity.CredentialLogRepository, error) {
	ctxOperation := "credential_log_create"

	dataLog := entity.CredentialLog{
		Value: data,
	}

	result, err := impl.repo.Create(ctx, dataLog, mr)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return nil, err
	}

	return result, nil
}

// FindByID implements CredentialLogService.
func (impl *CredentialLogImpl) FindByID(ctx context.Context, mr commonModel.MandatoryRequest, id string) ([]dto.CredentialLogResponse, error) {
	ctxOperation := "credential_log_find_by_id"

	idRequest, err := util.ObjectIDFromString(id)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return []dto.CredentialLogResponse{}, err
	}

	result, err := impl.repo.FindByID(ctx, idRequest)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogService(mr, ctxOperation, shared.PROCESS, err))
		return []dto.CredentialLogResponse{}, err
	}
	return toCredentialLogResponse(result), nil

}

func toCredentialLogResponse(logs []entity.CredentialLogRepository) []dto.CredentialLogResponse {
	result := make([]dto.CredentialLogResponse, 0)
	for _, log := range logs {
		lastUpdateDate := util.FormatTime(util.DateTimeFormatWithSecond, log.CreateDate)
		logRes := dto.CredentialLogResponse{
			LastUpdatedDate: lastUpdateDate,
			LastUpdatedBy:   log.CreateBy,
			LogData:         util.ObjToJson(log.Value),
		}
		result = append(result, logRes)
	}
	return result
}
