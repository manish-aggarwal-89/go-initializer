**{{SERVICE_TITLE}}**

# Table of content

- [Table Of Content](#table-of-content)
    - [Dependency Injection (di.go)](#dependency-injection-digo)
    - [Controller](#controller)

## Dependency Injection (di.go)

This is the code to create HTTP controller in Go. On this layer you must have di.go to init the controller layer. Base
code on dig.go is :

```go
package controller

import (
  "github.com/labstack/echo/v4/middleware"
  "go.uber.org/dig"
  "{{MODULE_NAME}}/internal/shared/deps"
)

type (
  Controllers struct {
    dig.In
    // add struct of controller like that
    HealthCheckController *HealthCheckController
    Deps                  deps.Deps
  }
)

const (
  // this is base path, for example http://localhost:8080/{{SERVICE_NAME}}
  base = "{{SERVICE_NAME}}"
)

func (impl *Controllers) Listen() {
  var e = impl.Deps.Echo

  // - register middleware
  e.Use(middleware.Gzip())
  e.Use(middleware.Recover())

  // - register routes
  impl.RegisterRoutes()

  if err := e.Start(fmt.Sprintf(":%d", impl.Deps.Config.HttpServerPort)); err != nil {
    impl.Deps.GetLogger(nil).Fatalf("failed to start http server %s", err)
  }
}

// register all http routes here
func (impl *Controllers) RegisterRoutes() {
  var (
    e = impl.Deps.Echo
  )

  // health check
  e.GET("/application/health", impl.HealthCheckController.HealthCheck)
  e.GET("/swagger/*", echoSwagger.WrapHandler)
}

// register your controllers here
func Register(container *dig.Container) error {
  //sample register healthcheck controller
  if err := container.Provide(NewHealthCheckController); err != nil {
    return errors.Wrap(err, "failed to provide health check controller")
  }
  return nil
}

```

## Controller

Sample code for controller

```go
package controller

type (
	HealthCheckController struct {
		deps    deps.Deps
        // adding data interface need to call in here
		service service.HealthCheckService
	}
)

// @Summary      Endpoint to check health service
// @Description  Endpoint to check health service
//
//	@Param			storeId		header		string	true	"Authentication header" default(TIKETCOM)
//	@Param			channelId	header		string	true	"Authentication header" default(WEB)
//	@Param			requestId	header		string	true	"Authentication header" default(321321321)
//	@Param			username	header		string	true	"Authentication header" default(guest)
//	@Param			serviceId	header		string	true	"Authentication header" default(GATEWAY)
//	@Param			resellerId	header		string	false	"Authentication header"
//	@Param			identity	header		string	false	"Authentication header"
//
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      500
// @Router       /application/health [get]
func (impl *HealthCheckController) HealthCheck(c echo.Context) error {
	// context operation tracing
	ctx_o := "health_check"
	var (
		ctx = c.Request().Context()
	)

	mr := common.MandatoryRequest{}
	mr = mr.BindMandatory(c)

	if err := impl.service.HealthCheck(ctx, mr); err != nil {
		return c.String(400, err.Error())
	}

	impl.deps.GetLogger(ctx).Infof("%s %s %s", shared.REST_IMPL, ctx_o, util.ObjToJson(mr))

	return c.String(200, "ok")
}


// constructor of the process
func NewHealthCheckController(deps deps.Deps, service service.HealthCheckService) (*HealthCheckController, error) {
	return &HealthCheckController{deps: deps, service: service}, nil
}

```

For generating swagger using swag init we need add comment for declare the generation of swagger. You can refer this
document too [swaggo](https://github.com/swaggo/swag) to check posibility comment for your business

```
// @Summary      Endpoint to check health service
// @Description  Endpoint to check health service
//
//	@Param			storeId		header		string	true	"Authentication header" default(TIKETCOM)
//	@Param			channelId	header		string	true	"Authentication header" default(WEB)
//	@Param			requestId	header		string	true	"Authentication header" default(321321321)
//	@Param			username	header		string	true	"Authentication header" default(guest)
//	@Param			serviceId	header		string	true	"Authentication header" default(GATEWAY)
//	@Param			resellerId	header		string	false	"Authentication header"
//	@Param			identity	header		string	false	"Authentication header"
//
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      500
// @Router       /application/health [get]
```
