package http

import (
	"github.com/gofiber/fiber/v2"
)

func (s *Server) initV1Routes(api fiber.Router) {
	apiV1 := api.Group("/v1")
	apiV1.Use(
		func(ctx *fiber.Ctx) error {
			return ctx.Next()
		},
		func(fctx *fiber.Ctx) error {
			fctx.Context().SetContentType(fiber.MIMEApplicationJSON)

			return fctx.Next()
		},
	)

	s.initHealthCheckRoutes(apiV1)
}

func (s *Server) initHealthCheckRoutes(apiV1 fiber.Router) {
	apiV1.Get("/health", s.healthCheck.HealthCheck)
}
