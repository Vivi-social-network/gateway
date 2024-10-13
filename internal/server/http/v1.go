package http

import (
	"github.com/gofiber/fiber/v2"
)

func (s *Server) initV1Routes(api fiber.Router) {
	apiV1 := api.Group("/v1")

	s.initHealthCheckRoutes(apiV1)
}

func (s *Server) initHealthCheckRoutes(apiV1 fiber.Router) {
	apiV1.Get("/health", s.healthCheck.HealthCheck)
}
