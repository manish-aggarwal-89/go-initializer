package service

import (
	"github.com/pkg/errors"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/errormapper"
	"go.uber.org/dig"
)

type (
	Services struct {
		dig.In
		HealthCheckService HealthCheckServiceItf
		SearchService      SearchServiceItf
	}
)

func Register(container *dig.Container) error {
	if err := container.Provide(NewHealthCheckService); err != nil {
		return errors.Wrap(err, "failed to provide health check service")
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

	// Integrator error mapping (common lib) + BAU analytics (common lib)
	if err := container.Provide(errormapper.NewIntegratorErrorMappingImpl, dig.As(new(errormapper.IntegratorErrorMappingService))); err != nil {
		return errors.Wrap(err, "failed to provide integrator error mapping service")
	}

	if err := container.Provide(NewBauSearchService); err != nil {
		return errors.Wrap(err, "failed to provide BAU search service")
	}
	if err := container.Provide(NewBauBookingService); err != nil {
		return errors.Wrap(err, "failed to provide BAU booking service")
	}
	if err := container.Provide(NewBauIssuedService); err != nil {
		return errors.Wrap(err, "failed to provide BAU issued service")
	}
	return nil
}
