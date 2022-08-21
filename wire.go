//go:build wireinject
// +build wireinject

package ticketAllocating_Purchasing

import (
	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api"
	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/handler"
	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/healthcheck"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

var repositoryProviders = wire.NewSet()

var serviceProviders = wire.NewSet()

var orchestratorProviders = wire.NewSet()

var handlerProviders = wire.NewSet(
	handler.NewHomeHandler,
)

var allProviders = wire.NewSet(
	repositoryProviders,
	serviceProviders,
	orchestratorProviders,
	handlerProviders,
)

func InitHealthCheck() healthcheck.IHealthCheckHandler {
	wire.Build(healthcheck.NewHealthCheckHandler)
	return nil
}

func InitRoute(
	l *logrus.Logger,
) api.IRoute {
	wire.Build(allProviders, api.NewRoute)
	return nil
}
