package main

import (
	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/sirupsen/logrus"

	di "github.com/Selahattinn/ticketAllocating-Purchasing"
	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api"
	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/middleware"
	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/utils"
	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/validation"
)

const (
	bodyLimit = 10 * 1024 * 1024
)

type application struct {
	logger          *logrus.Logger
	postgreInstance mysql.IMysqlInstance
}

func initApplication(a *application) *fiber.App {
	app := fiber.New(fiber.Config{
		BodyLimit:    bodyLimit,
		ErrorHandler: middleware.ErrorMiddleware(),
	})

	a.addHealthCheckRoutes(app)
	a.addCommonMiddleware(app)

	route := di.InitRoute(
		a.logger,
		a.postgreInstance,
	)
	route.SetupRoutes(&api.RouteContext{
		App: app,
	})

	app.Use(middleware.NotFoundMiddleware())

	return app
}

func (a *application) addHealthCheckRoutes(app *fiber.App) {
	healthCheckHandler := di.InitHealthCheck()
	app.Get("/liveness", healthCheckHandler.Liveness)
	app.Get("/readiness", healthCheckHandler.Readiness)
}

func (a *application) addCommonMiddleware(app *fiber.App) {
	app.Use(requestid.New())
	app.Use(middleware.LoggerMiddleware(a.logger))

	validator := validation.InitValidator()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals(utils.ValidatorKey, validator)
		return c.Next()
	})
}
