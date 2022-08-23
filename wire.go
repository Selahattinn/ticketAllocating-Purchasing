//go:build wireinject
// +build wireinject

package ticketAllocating_Purchasing

import (
	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api"
	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/handler"
	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/orchestration"
	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/ticket"
	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/healthcheck"
	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/mysql"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

var repositoryProviders = wire.NewSet(
	ticket.NewTicketRepository,
)

var serviceProviders = wire.NewSet(
	ticket.NewTicketService,
)

var orchestratorProviders = wire.NewSet(
	orchestration.NewTicketOrchestrator,
)

var handlerProviders = wire.NewSet(
	handler.NewHomeHandler,
	handler.NewTicketHandler,
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
	pi mysql.IMysqlInstance,
) api.IRoute {
	wire.Build(allProviders, api.NewRoute)
	return nil
}
