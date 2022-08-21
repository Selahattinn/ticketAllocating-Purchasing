package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/Selahattinn/ticketAllocating-Purchasing/configs"
	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/dto/resource"
	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/response"
)

type IHomeHandler interface {
	Home(c *fiber.Ctx) error
}

type homeHandler struct{}

func NewHomeHandler() IHomeHandler {
	return &homeHandler{}
}

func (h *homeHandler) Home(c *fiber.Ctx) error {
	return c.JSON(response.NewSuccessResponse(&resource.HomeResource{
		App:  configs.TicketApp.Web.AppName,
		Env:  configs.TicketApp.Web.Env,
		Time: time.Now(),
	}))
}
