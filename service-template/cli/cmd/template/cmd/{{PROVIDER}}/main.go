package main

import (
	"os"
	"os/signal"
	"syscall"
	"{{MODULE_NAME}}/internal/controller"
	"{{MODULE_NAME}}/internal/di"
	"{{MODULE_NAME}}/internal/inbound"
	"{{MODULE_NAME}}/internal/scheduler"
	"{{MODULE_NAME}}/internal/shared/deps"

	_ "{{MODULE_NAME}}/docs"

	"github.com/labstack/gommon/log"
)

//	@title			{{SERVICE_TITLE}}
//	@version		1.0
//	@description	Swagger tools for {{PROVIDER}} integrator

// @contact.name	@supplyTeam
func main() {
	err := di.Container.Invoke(func(deps deps.Deps, listeners inbound.Listeners, controllers controller.Controllers, sch scheduler.Scheduler) error {
		ch := make(chan os.Signal, 1)

		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

		// - listen kafka
		go listeners.Listen()

		// listen to controllers
		go controllers.Listen()

		// credential scheduler
		go func(ch chan<- os.Signal) {
			err := sch.SchedulerCredential.CredentialPredefine()
			if err != nil {
				log.Errorf("error starting credential scheduler - %s", err.Error())
				ch <- os.Interrupt
			}
		}(ch)

		// system parameter scheduler
		go func(ch chan<- os.Signal) {
			err := sch.SchedulerSystemParameter.RefreshSystemParameterPredifine()
			if err != nil {
				log.Errorf("error starting system parameter scheduler - %s", err.Error())
				ch <- os.Interrupt
			}
		}(ch)

		<-ch
		deps.GetLogger(nil).Info("shutting down application and releasing resources")
		err := deps.Close()
		if err != nil {
			log.Error("error closing dependencies " + err.Error())
		}
		deps.GetLogger(nil).Info("done!")

		return nil
	})
	if err != nil {
		log.Panicf("application failed to start %s", err)
	}
}
