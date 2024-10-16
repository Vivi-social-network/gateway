package http

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Vivi-social-network/core/logger"
	"github.com/Vivi-social-network/gateway/internal/config"
	"github.com/Vivi-social-network/gateway/internal/server/http/handlers"
)

type Server struct {
	addr        string
	enablePprof bool
	env         config.Env

	srv *fiber.App

	logger *logger.Logger

	healthCheck *handlers.HealthCheck
}

func New(
	cfg config.HTTPServer,
	env config.Env,
	logger *logger.Logger,
	healthCheck *handlers.HealthCheck,
) (*Server, error) {
	fiberSrv := fiber.New(fiber.Config{
		UnescapePath:      cfg.UnescapePath,
		BodyLimit:         cfg.BodyLimit,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		AppName:           cfg.AppName,
		EnablePrintRoutes: cfg.EnablePrintRoutes,
		Network:           fiber.NetworkTCP,
	})

	srv := &Server{
		addr:        cfg.Address,
		env:         env,
		enablePprof: cfg.EnablePprof,

		srv:    fiberSrv,
		logger: logger,

		healthCheck: healthCheck,
	}
	srv.initRoutes()

	return srv, nil
}

func (s *Server) Listen(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		err := s.srv.Shutdown()
		if err != nil {
			s.logger.Error("shutdown https server failed", zap.Error(err))
		}
	}()

	return s.srv.Listen(s.addr)
}

func (s *Server) initRoutes() {
	s.srv.Use(
		recover.New(recover.Config{
			EnableStackTrace: true,
		}),
		requestid.New(requestid.Config{
			Generator: func() string {
				return uuid.New().String()
			},
		}),
		cors.New(),
	)

	if s.enablePprof {
		s.srv.Use(pprof.New(pprof.Config{}))
	}

	api := s.srv.Group("/api")

	s.initV1Routes(api)
}
