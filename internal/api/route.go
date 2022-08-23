package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/Selahattinn/ticketAllocating-Purchasing/configs"
	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/handler"
	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/swagger"
)

type RouteContext struct {
	App *fiber.App
}

type IRoute interface {
	SetupRoutes(r *RouteContext)
}

type route struct {
	homeHandler   handler.IHomeHandler
	ticketHandler handler.ITicketHandler
}

func NewRoute(
	hHandler handler.IHomeHandler,
	tHandler handler.ITicketHandler,
) IRoute {
	return &route{
		homeHandler:   hHandler,
		ticketHandler: tHandler,
	}
}

func (r *route) SetupRoutes(rc *RouteContext) {
	if !configs.TicketApp.Web.IsProductionEnv() {
		r.docRoutes(rc.App)
	}

	v1Group := rc.App.Group("/v1")
	v1Group.Get("/", r.homeHandler.Home)

	r.TicketRoutes(v1Group)
	r.TicketOptionsRoutes(v1Group)
}

func (r *route) docRoutes(fr fiber.Router) {
	fr.Get("/swagger/*", swagger.New(swagger.Config{
		URL:         fmt.Sprintf("/swagger/%s", swagger.DefaultDocURL),
		DeepLinking: true,
	}))
}

func (r *route) TicketRoutes(fr fiber.Router) {
	ticketsGroup := fr.Group("/ticket")

	ticketsGroup.Get("/:id", r.ticketHandler.Get)
}

func (r *route) TicketOptionsRoutes(fr fiber.Router) {
	ticketOptionsGroup := fr.Group("/ticket_options")

	ticketOptionsGroup.Post("/", r.ticketHandler.Create)
	ticketOptionsGroup.Post("/:id/purchase", r.ticketHandler.Purchase)
}
