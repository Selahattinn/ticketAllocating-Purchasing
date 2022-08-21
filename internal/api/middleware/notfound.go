package middleware

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/response"
	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/utils"
)

func NotFoundMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		errBag := utils.ErrorBag{Code: utils.NotFoundErrCode}

		return c.Status(fiber.StatusNotFound).JSON(response.NewErrorResponse(c.Context(), errBag, utils.NotFoundMsg))
	}
}
