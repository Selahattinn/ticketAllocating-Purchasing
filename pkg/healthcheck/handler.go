package healthcheck

import (
	"github.com/gofiber/fiber/v2"
)

type IHealthCheckHandler interface {
	Liveness(c *fiber.Ctx) error
	Readiness(c *fiber.Ctx) error
}

type healthCheckHandler struct{}

func NewHealthCheckHandler() IHealthCheckHandler {
	return &healthCheckHandler{}
}

func (h *healthCheckHandler) Liveness(c *fiber.Ctx) error {
	if !Liveness() {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": LivenessStatusShutdown})
	}

	return c.JSON(fiber.Map{"status": LivenessStatusOk})
}

func (h *healthCheckHandler) Readiness(c *fiber.Ctx) error {
	readiness := Readiness()
	if !IsConnectionSuccessful(readiness) {
		return c.Status(fiber.StatusInternalServerError).JSON(readiness)
	}

	return c.JSON(fiber.Map{"status": ReadinessStatusOk})
}
