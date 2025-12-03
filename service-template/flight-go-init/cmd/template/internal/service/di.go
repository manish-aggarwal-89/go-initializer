package service

import (
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

type (
	Services struct {
		dig.In
		HealthCheckService     HealthCheckServiceItf
		SystemParameterService SystemParameterServiceItf
		SearchService          SearchServiceItf
	}
)

func Register(container *dig.Container) error {
	if err := container.Provide(NewHealthCheckService); err != nil {
		return errors.Wrap(err, "failed to provide health check service")
	}
	if err := container.Provide(NewCredentialService, dig.As(new(CredentialServiceItf))); err != nil {
		return errors.Wrap(err, "failed to provide credential service")
	}
	if err := container.Provide(NewCredentialMasterDataService); err != nil {
		return errors.Wrap(err, "failed to provide credential master data service")
	}
	if err := container.Provide(NewCredentialLogService); err != nil {
		return errors.Wrap(err, "failed to provide credential log service")
	}
	if err := container.Provide(NewSystemParameterServiceImpl); err != nil {
		return errors.Wrap(err, "failed to provide system parameter service")
	}
	if err := container.Provide(NewKafkaPublisherService, dig.As(new(KafkaPublisherServiceItf))); err != nil {
		return errors.Wrap(err, "failed to provide kafka publisher service")
	}
	if err := container.Provide(NewSearchServiceImpl, dig.As(new(SearchServiceItf))); err != nil {
		return errors.Wrap(err, "failed to provide search service")
	}
	if err := container.Provide(NewBookingServiceImpl, dig.As(new(BookingServiceItf))); err != nil {
		return errors.Wrap(err, "failed to provide booking service")
	}
	if err := container.Provide(NewCancelBookingServiceImpl, dig.As(new(CancelBookingServiceItf))); err != nil {
		return errors.Wrap(err, "failed to provide cancel-booking service")
	}
	if err := container.Provide(NewIssuedServiceImpl, dig.As(new(IssuedServiceItf))); err != nil {
		return errors.Wrap(err, "failed to provide issued service")
	}
	if err := container.Provide(NewBAUSearchService); err != nil {
		return errors.Wrap(err, "failed to provide BAU search service")
	}
	if err := container.Provide(NewBAUBookingService); err != nil {
		return errors.Wrap(err, "failed to provide BAU booking service")
	}
	if err := container.Provide(NewBAUIssuedService); err != nil {
		return errors.Wrap(err, "failed to provide BAU issued service")
	}
	if err := container.Provide(NewSessionServiceImpl, dig.As(new(SessionService))); err != nil {
		return errors.Wrap(err, "failed to provide Session service")
	}

	return nil
}
