package utils

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

func GetContextFromFiber(c *fiber.Ctx) context.Context {
	return c.Locals(ContextKey).(context.Context)
}
