package repository

import (
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

type (
	Repositories struct {
		dig.In
		SystemParameterRepository SystemParameterInterface
	}
)

func Register(container *dig.Container) error {
	if err := container.Provide(NewCredentialRepository); err != nil {
		return errors.Wrap(err, "failed to provide credential repository")
	}
	if err := container.Provide(NewCredentialLogRepository); err != nil {
		return errors.Wrap(err, "failed to provide credentia log repository")
	}
	if err := container.Provide(NewCredentialMasterDataRepository); err != nil {
		return errors.Wrap(err, "failed to provide credential master data repository")
	}
	if err := container.Provide(NewSystemParameterRepository); err != nil {
		return errors.Wrap(err, "failed to provide SystemParameterRepository")
	}
	if err := container.Provide(NewCacheRepository, dig.As(new(CacheRepository))); err != nil {
		return errors.Wrap(err, "failed to provide cache repository")
	}
	if err := container.Provide(NewCurrencyRepository, dig.As(new(CurrencyRepositoryInterface))); err != nil {
		return errors.Wrap(err, "failed to provide currency repository")
	}

	return nil
}
