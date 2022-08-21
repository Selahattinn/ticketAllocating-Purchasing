package api

import (
	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/handler"
	"github.com/gofiber/fiber/v2"
)

type RouteContext struct {
	App *fiber.App
}

type IRoute interface {
	SetupRoutes(r *RouteContext)
}

type route struct {
	homeHandler handler.IHomeHandler
}

func NewRoute(
	hHandler handler.IHomeHandler,
) IRoute {
	return &route{
		homeHandler: hHandler,
	}
}

func (r *route) SetupRoutes(rc *RouteContext) {
	v1Group := rc.App.Group("/v1")
	v1Group.Get("/", r.homeHandler.Home)

}
