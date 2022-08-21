package middleware

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/response"
	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/utils"
)

func ErrorMiddleware() func(c *fiber.Ctx, err error) error {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		errBag := utils.ErrorBag{Code: utils.UnexpectedErrCode}

		return c.Status(code).JSON(response.NewErrorResponse(c.Context(), errBag, utils.UnexpectedMsg))
	}
}
