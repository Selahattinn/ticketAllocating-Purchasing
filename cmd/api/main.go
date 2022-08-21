package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Selahattinn/ticketAllocating-Purchasing/configs"
	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/healthcheck"
)

func main() {
	if configErr := initConfig(); configErr != nil {
		log.Fatalf("initialization config: %v", configErr)
	}

	logger := initLogger()
	app, appErr := boot(logger)

	if appErr != nil {
		logger.Fatalf("initialization: %v", appErr)
	}
	defer shutdown(app)

	server := initApplication(app)

	go func() {
		healthcheck.InitHealthCheck()

		if serveErr := server.Listen(fmt.Sprintf(":%s", configs.TicketApp.Web.Port)); serveErr != nil {
			logger.Fatalf("connection: web server %v", serveErr)
		}
	}()

	logger.Info("application successfully configured")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	<-c

	healthcheck.ServerShutdown()
	if shutdownErr := server.Shutdown(); shutdownErr != nil {
		logger.Errorf("shutdown: server %v", shutdownErr)
	}
}

func shutdown(app *application) {
	app.logger.Info("shutdown: completed")
}
