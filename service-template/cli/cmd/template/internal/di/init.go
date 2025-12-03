package di

import (
	"fmt"
	"{{MODULE_NAME}}/config"
	"{{MODULE_NAME}}/internal/controller"
	"{{MODULE_NAME}}/internal/inbound"
	"{{MODULE_NAME}}/internal/outbound"
	"{{MODULE_NAME}}/internal/predefine"
	"{{MODULE_NAME}}/internal/repository"
	"{{MODULE_NAME}}/internal/scheduler"
	"{{MODULE_NAME}}/internal/service"
	"{{MODULE_NAME}}/internal/shared/deps"

	"go.uber.org/dig"
)

var Container = dig.New()

func init() {
	// - config
	if err := Container.Provide(config.New); err != nil {
		panic(fmt.Sprintf("failed to provide config %s", err))
	}

	// - metrics
	if err := Container.Provide(NewMetric); err != nil {
		panic(fmt.Sprintf("failed to provide metric %s", err))
	}

	// - logger
	if err := Container.Provide(NewLogger); err != nil {
		panic(fmt.Sprintf("failed to provide logger %s", err))
	}
	// - logger using logrus
	if err := Container.Provide(NewLoggerV2); err != nil {
		panic(fmt.Sprintf("failed to provide logger v2 %s", err))
	}

	// - mongoDB
	if err := Container.Provide(NewMongoDB); err != nil {
		panic(fmt.Sprintf("failed to provide mongoDB %s", err))
	}

	// - echo
	if err := Container.Provide(NewEcho); err != nil {
		panic(fmt.Sprintf("failed to provide echo %s", err))
	}

	// - deps
	if err := deps.Register(Container); err != nil {
		panic(fmt.Sprintf("failed to register deps %s", err))
	}

	// - inbound
	if err := inbound.Register(Container); err != nil {
		panic(fmt.Sprintf("failed to register inbound %s", err))
	}

	// - controllers
	if err := controller.Register(Container); err != nil {
		panic(fmt.Sprintf("failed to register controller %s", err))
	}

	// - service
	if err := service.Register(Container); err != nil {
		panic(fmt.Sprintf("failed to register service %s", err))
	}

	// - repository
	if err := repository.Register(Container); err != nil {
		panic(fmt.Sprintf("failed to register repository %s", err))
	}

	// - scheduler
	if err := scheduler.Register(Container); err != nil {
		panic(fmt.Sprintf("failed to register scheduler %s", err))
	}

	// - predefine
	if err := predefine.Register(Container); err != nil {
		panic(fmt.Sprintf("failed to register predefine %s", err))
	}

	// - outbound
	if err := outbound.Register(Container); err != nil {
		panic(fmt.Sprintf("failed to register predefine %s", err))
	}
}
