package service

import (
	"{{MODULE_NAME}}/config"
	"{{MODULE_NAME}}/internal/outbound"
	"{{MODULE_NAME}}/internal/service/rule_supplier"
	"{{MODULE_NAME}}/internal/service/rule_supplier/{{PROVIDER}}"
	"{{MODULE_NAME}}/internal/shared/deps"

	"github.com/pkg/errors"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/errormapper"
	"go.uber.org/dig"
)

type RuleRegistryParams struct {
	dig.In
	Deps     deps.Deps
	Config   *config.Config
	Outbound outbound.OutboundItf
}

func ProvideSearchRuleRegistry(p RuleRegistryParams) *rule_supplier.SearchRuleRegistry {
	reg := rule_supplier.NewSearchRuleRegistry()
	distType := p.Config.DistributionType
	if distType == "" {
		panic("config.DistributionType is required but not set")
	}
	base := {{PROVIDER}}.Base{Deps: p.Deps, Outbound: p.Outbound}
	reg.Register(distType, {{PROVIDER}}.NewSearchRule(base))
	return reg
}

func ProvideBookingRuleRegistry(p RuleRegistryParams) *rule_supplier.BookingRuleRegistry {
	reg := rule_supplier.NewBookingRuleRegistry()
	distType := p.Config.DistributionType
	if distType == "" {
		panic("config.DistributionType is required but not set")
	}
	base := {{PROVIDER}}.Base{Deps: p.Deps, Outbound: p.Outbound}
	reg.Register(distType, {{PROVIDER}}.NewBookingRule(base))
	return reg
}

func ProvideIssuedRuleRegistry(p RuleRegistryParams) *rule_supplier.IssuedRuleRegistry {
	reg := rule_supplier.NewIssuedRuleRegistry()
	distType := p.Config.DistributionType
	if distType == "" {
		panic("config.DistributionType is required but not set")
	}
	base := {{PROVIDER}}.Base{Deps: p.Deps, Outbound: p.Outbound}
	reg.Register(distType, {{PROVIDER}}.NewIssuedRule(base))
	return reg
}

func ProvideCancelBookingRuleRegistry(p RuleRegistryParams) *rule_supplier.CancelBookingRuleRegistry {
	reg := rule_supplier.NewCancelBookingRuleRegistry()
	distType := p.Config.DistributionType
	if distType == "" {
		panic("config.DistributionType is required but not set")
	}
	base := {{PROVIDER}}.Base{Deps: p.Deps, Outbound: p.Outbound}
	reg.Register(distType, {{PROVIDER}}.NewCancelBookingRule(base))
	return reg
}

func Register(container *dig.Container) error {
	if err := container.Provide(NewHealthCheckService); err != nil {
		return errors.Wrap(err, "failed to provide health check service")
	}

	if err := container.Provide(NewKafkaPublisherService, dig.As(new(KafkaPublisherServiceItf))); err != nil {
		return errors.Wrap(err, "failed to provide kafka publisher service")
	}

	// Rule registries
	if err := container.Provide(ProvideSearchRuleRegistry); err != nil {
		return errors.Wrap(err, "failed to provide search rule registry")
	}
	if err := container.Provide(ProvideBookingRuleRegistry); err != nil {
		return errors.Wrap(err, "failed to provide booking rule registry")
	}
	if err := container.Provide(ProvideIssuedRuleRegistry); err != nil {
		return errors.Wrap(err, "failed to provide issued rule registry")
	}
	if err := container.Provide(ProvideCancelBookingRuleRegistry); err != nil {
		return errors.Wrap(err, "failed to provide cancel booking rule registry")
	}

	// Processor services
	if err := container.Provide(NewSearchProcessorService); err != nil {
		return errors.Wrap(err, "failed to provide search processor service")
	}
	if err := container.Provide(NewBookingProcessorService); err != nil {
		return errors.Wrap(err, "failed to provide booking processor service")
	}
	if err := container.Provide(NewIssuedProcessorService); err != nil {
		return errors.Wrap(err, "failed to provide issued processor service")
	}
	if err := container.Provide(NewCancelBookingProcessorService); err != nil {
		return errors.Wrap(err, "failed to provide cancel booking processor service")
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
