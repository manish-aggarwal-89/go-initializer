package controller

import (
	"fmt"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/errormapper"
	"{{MODULE_NAME}}/internal/shared/deps"

	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/credential"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/credential_master_data"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/promotion"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/system_param"
	"go.uber.org/dig"
)

type (
	Controllers struct {
		dig.In
		HealthCheckController   *HealthCheckController
		SearchController        *SearchControllerImpl
		BookingController       *BookingController
		CancelBookingController *CancelBookingController
		IssuedController        *IssuedController

		// Controllers from common lib (github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO)
		// provided by common lib Register() functions in internal/di/init.go
		SystemParamController            *system_param.Controller
		CredentialController             *credential.Controller
		CredentialMasterDataController   *credential_master_data.Controller
		PromotionController              *promotion.Controller
		IntegratorErrorMappingController *errormapper.Controller

		Deps deps.Deps
	}
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

// @basePath /{{BASE_PATH}}
func GetBasePath(distributionType string) string {
	return "{{BASE_PATH}}"
}

func (impl *Controllers) RegisterRoutes() {
	var e = impl.Deps.Echo
	base := GetBasePath(impl.Deps.Config.DistributionType)

	e.GET(base+"/application/health", impl.HealthCheckController.HealthCheck)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	primeBooking := e.Group(base)
	primeBooking.POST("/search", impl.SearchController.Search)
	primeBooking.POST("/v2/search", impl.SearchController.SearchV2)
	primeBooking.POST("/booking", impl.BookingController.Booking)
	primeBooking.POST("/issued", impl.IssuedController.Issued)
	primeBooking.POST("/cancel-booking", impl.CancelBookingController.CancelBook)

	// Controllers from common lib - use RegisterRoutesWithGroup
	credentialGroup := e.Group(base + "/credential")
	impl.CredentialController.RegisterRoutesWithGroup(credentialGroup)

	credentialMasterDataGroup := e.Group(base + "/credential-master-data")
	impl.CredentialMasterDataController.RegisterRoutesWithGroup(credentialMasterDataGroup)

	sysParamGroup := e.Group(base + "/system-parameter")
	impl.SystemParamController.RegisterRoutesWithGroup(sysParamGroup)

	promotionGroup := e.Group(base + "/promotion")
	impl.PromotionController.RegisterRoutesWithGroup(promotionGroup)

	// Integrator Error Mapping routes
	integratorErrorMapping := e.Group(base + "/integrator-error-mapping")
	impl.IntegratorErrorMappingController.RegisterRoutesWithGroup(integratorErrorMapping)
}

func Register(container *dig.Container) error {
	if err := container.Provide(NewHealthCheckController); err != nil {
		return errors.Wrap(err, "failed to provide health check controller")
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
	if err := container.Provide(errormapper.NewController); err != nil {
		return errors.Wrap(err, "failed to provide issued controller")
	}
	// CredentialController, CredentialMasterDataController, SystemParamController,IntegratorErrorMappingController and PromotionController
	// are provided by common lib Register() functions in internal/di/init.go
	return nil
}
