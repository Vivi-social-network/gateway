package handlers

import (
	"github.com/gofiber/fiber/v2"
)

type HealthCheck struct{}

func NewHealthCheck() *HealthCheck {
	return &HealthCheck{}
}

func (h *HealthCheck) HealthCheck(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{"status": "UP"})
}
