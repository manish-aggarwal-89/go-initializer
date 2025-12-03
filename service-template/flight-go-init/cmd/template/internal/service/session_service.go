package service

import (
	"context"
	"errors"
	"fmt"
	"{{MODULE_NAME}}/internal/outbound"
	"{{MODULE_NAME}}/internal/predefine"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/dto"
	"{{MODULE_NAME}}/internal/shared/entity"

	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
)

type SessionService interface {
	GetToken(ctx context.Context, mr common.MandatoryRequest, supplierCode string) (dto.SessionSignatureDto, error)
	SetFreshStateToContext(ctx context.Context, mr common.MandatoryRequest, supplierCode string) (resCtx context.Context, err error)
	SetStateToContextFromCache(ctx context.Context, mr common.MandatoryRequest, supplierCode string) (resCtx context.Context, err error)
	SetStateToContextFromCacheWithFallback(ctx context.Context, mr common.MandatoryRequest, supplierCode string) (resCtx context.Context, err error)
	SetFreshStateToContextWithCredential(ctx context.Context, mr common.MandatoryRequest, credential entity.CredentialRepository) (resCtx context.Context, err error)
}

func NewSessionServiceImpl(deps deps.Deps, outbound outbound.OutboundItf, predefine predefine.Predefined, credentialSvc CredentialServiceItf) SessionService {
	return &SessionServiceImpl{
		deps:          deps,
		outbound:      outbound,
		predefine:     predefine,
		credentialSvc: credentialSvc,
	}
}

type SessionServiceImpl struct {
	deps          deps.Deps
	outbound      outbound.OutboundItf
	predefine     predefine.Predefined
	credentialSvc CredentialServiceItf
}

func (impl *SessionServiceImpl) GetToken(ctx context.Context, mr common.MandatoryRequest, supplierCode string) (dto.SessionSignatureDto, error) {
	var (
		token string
		err   error
	)
	//credential, err := impl.credentialSvc.FindBySupplierInCache(ctx, mr, supplierCode)
	//if err != nil {
	//	return dto.SessionSignatureDto{}, err
	//}

	// Obtain token by calling outbound

	return dto.SessionSignatureDto{
		AuthToken: token,
	}, err
}

func (impl *SessionServiceImpl) SetFreshStateToContext(ctx context.Context, mr common.MandatoryRequest, supplierCode string) (resCtx context.Context, err error) {
	opCtx := "SessionServiceImpl.SetFreshStateToContext"
	var token string

	impl.deps.GetLogger(ctx).Infof("[%s] %s: exec", mr.RequestId, opCtx)
	defer func() {
		if err != nil {
			impl.deps.GetLogger(ctx).Errorf("[%s] %s, error: %+v", mr.RequestId, opCtx, err)
		}
	}()

	//repoCredential, err := impl.credentialSvc.FindBySupplierInCache(ctx, mr, supplierCode)
	//if err != nil {
	//	return nil, err
	//}

	// Obtain token by calling outbound

	newCtx := context.WithValue(ctx, shared.AuthorizationTokenKey, token)
	return newCtx, nil
}

func (impl *SessionServiceImpl) SetStateToContextFromCache(ctx context.Context, mr common.MandatoryRequest, supplierCode string) (resCtx context.Context, err error) {
	opCtx := "SessionServiceImpl.SetStateToContextFromCache"

	impl.deps.GetLogger(ctx).Infof("[%s] %s: exec", mr.RequestId, opCtx)
	defer func() {
		if err != nil {
			impl.deps.GetLogger(ctx).Errorf("[%s] %s, error: %+v", mr.RequestId, opCtx, err)
		}
	}()

	predefinedData, err := impl.predefine.GetDataByID(ctx, fmt.Sprintf(shared.PREDEFINED_SESSION_CACHE_KEY_TEMPLATE, supplierCode))
	if err != nil {
		return nil, err
	}

	cachedSessionData, ok := predefinedData.(dto.SessionSignatureDto)
	if !ok {
		return nil, errors.New("failed to cast cached session data")
	}

	newCtx := context.WithValue(ctx, shared.AuthorizationTokenKey, cachedSessionData.AuthToken)
	return newCtx, nil
}

func (impl *SessionServiceImpl) SetStateToContextFromCacheWithFallback(ctx context.Context, mr common.MandatoryRequest, supplierCode string) (resCtx context.Context, err error) {
	opCtx := "SessionServiceImpl.SetStateToContextFromCacheWithFallback"

	impl.deps.GetLogger(ctx).Infof("[%s] %s: exec", mr.RequestId, opCtx)
	defer func() {
		if err != nil {
			impl.deps.GetLogger(ctx).Errorf("[%s] %s, error: %+v", mr.RequestId, opCtx, err)
		}
	}()

	resCtx, err = impl.SetStateToContextFromCache(ctx, mr, supplierCode)
	if err != nil {
		impl.deps.GetLogger(ctx).Infof("[%s] %s: failed setting state from cache, fallback to SetFreshStateToContext", mr.RequestId, opCtx)

		resCtx, err = impl.SetFreshStateToContext(ctx, mr, supplierCode)
		if err != nil {
			return nil, err
		}
	}

	return resCtx, nil
}

func (impl *SessionServiceImpl) SetFreshStateToContextWithCredential(ctx context.Context, mr common.MandatoryRequest, credential entity.CredentialRepository) (resCtx context.Context, err error) {
	opCtx := "SessionServiceImpl.SetFreshStateToContextWithCredential"

	var token string
	impl.deps.GetLogger(ctx).Infof("[%s] %s: exec", mr.RequestId, opCtx)
	defer func() {
		if err != nil {
			impl.deps.GetLogger(ctx).Errorf("[%s] %s, error: %+v", mr.RequestId, opCtx, err)
		}
	}()

	//call outbound to get token

	newCtx := context.WithValue(ctx, shared.AuthorizationTokenKey, token)

	return newCtx, nil
}
