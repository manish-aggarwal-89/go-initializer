package controller

import (
	"fmt"
	"{{MODULE_NAME}}/internal/shared/deps"

	"github.com/labstack/echo/v4/middleware"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/pkg/errors"
	"go.uber.org/dig"
)

type (
	Controllers struct {
		dig.In
		HealthCheckController          *HealthCheckController
		SearchController               *SearchControllerImpl
		BookingController              *BookingController
		CancelBookingController        *CancelBookingController
		CredentialController           *CredentialController
		CredentialMasterDataController *CredentialMasterDataController
		SystemParameterController      *SystemParameterController
		IssuedController               *IssuedController
		Deps                           deps.Deps
	}
)

const (
	base = "{{BASE_PATH}}"
)

func (impl *Controllers) Listen() {
	var e = impl.Deps.Echo

	e.Use(middleware.Recover())
	e.Use(HawkEyeAPIInbound(impl.Deps))
	e.Use(MandatoryRequestMiddleware(impl.Deps))
	impl.RegisterRoutes()
	if err := e.Start(fmt.Sprintf(":%d", impl.Deps.Config.HttpServerPort)); err != nil {
		impl.Deps.GetLogger(nil).Fatalf("failed to start http server %s", err)
	}
}

func (impl *Controllers) RegisterRoutes() {
	var e = impl.Deps.Echo

	e.GET("/application/health", impl.HealthCheckController.HealthCheck)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	primeBooking := e.Group(base)
	primeBooking.POST("/search", impl.SearchController.Search)
	primeBooking.POST("/v2/search", impl.SearchController.SearchV2)
	primeBooking.POST("/booking", impl.BookingController.Booking)
	primeBooking.POST("/issued", impl.IssuedController.Issued)
	primeBooking.POST("/cancel-booking", impl.CancelBookingController.CancelBook)

	credentialMasterData := e.Group(base + "/credential-master-data")
	credentialMasterData.POST("", impl.CredentialMasterDataController.Create)
	credentialMasterData.GET("", impl.CredentialMasterDataController.FindAllMap)
	credentialMasterData.GET("/list", impl.CredentialMasterDataController.FindAllList)
	credentialMasterData.PUT("/:id", impl.CredentialMasterDataController.UpdateByID)
	credentialMasterData.DELETE("/:id", impl.CredentialMasterDataController.DeleteByID)

	credential := e.Group(base + "/credential")
	credential.GET("/alert", impl.CredentialController.SendAlert)
	credential.GET("/:id", impl.CredentialController.FindByID)
	credential.GET("", impl.CredentialController.FindWithPaginated)
	credential.POST("", impl.CredentialController.Create)
	credential.PUT("/:id", impl.CredentialController.UpdateByID)
	credential.DELETE("/:id", impl.CredentialController.DeleteByID)
	credential.GET("/log/:id", impl.CredentialController.FindLogByID)

	sys_param := e.Group(base + "/system-parameter")
	sys_param.GET("", impl.SystemParameterController.FindAll)
	sys_param.GET("/:id", impl.SystemParameterController.FindByID)
	sys_param.GET("/variable/:variable", impl.SystemParameterController.FindByVariable)
	sys_param.POST("", impl.SystemParameterController.Create)
	sys_param.PUT("/:id", impl.SystemParameterController.UpdateByID)
	sys_param.DELETE("/:id", impl.SystemParameterController.DeleteByID)
}

func Register(container *dig.Container) error {
	if err := container.Provide(NewHealthCheckController); err != nil {
		return errors.Wrap(err, "failed to provide health check controller")
	}
	if err := container.Provide(NewCredentialController); err != nil {
		return errors.Wrap(err, "failed to provide credential controller")
	}
	if err := container.Provide(NewCredentialMasterDataController); err != nil {
		return errors.Wrap(err, "failed to provide credential master data controller")
	}
	if err := container.Provide(NewSystemParameterController); err != nil {
		return errors.Wrap(err, "failed to provide system parameter controller")
	}
	if err := container.Provide(NewSearchController); err != nil {
		return errors.Wrap(err, "failed to provide search controller")
	}
	if err := container.Provide(NewBookingController); err != nil {
		return errors.Wrap(err, "failed to provide search controller")
	}
	if err := container.Provide(NewCancelBookingController); err != nil {
		return errors.Wrap(err, "failed to provide cancel booking controller")
	}
	if err := container.Provide(NewIssuedController); err != nil {
		return errors.Wrap(err, "failed to provide issued controller")
	}

	return nil
}
