package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/Vivi-social-network/core/logger"
	"github.com/Vivi-social-network/gateway/internal/config"
	"github.com/Vivi-social-network/gateway/internal/server/http"
	"github.com/Vivi-social-network/gateway/internal/server/http/handlers"
)

var (
	configPath = flag.String("cfg", "configs/dev.yaml", "Config file path")
)

func main() {
	ctx, cancelCtx := context.WithCancel(context.Background())
	syscallNotifyChannel := make(chan os.Signal, 1)
	handleSysCalls(syscallNotifyChannel)

	go func() {
		for range syscallNotifyChannel {
			cancelCtx()
		}
	}()

	flag.Parse()
	if len(*configPath) == 0 {
		panic("config path must not be empty")
	}

	cfg, err := config.Parse(*configPath)
	if err != nil {
		panic(fmt.Sprintf("cannot parse config: %v", err))
	}

	log := logger.New(cfg.Logger)

	log.Info("service starting", "env", cfg.Env, "CPUs", runtime.NumCPU())

	log.Info("configure handlers")
	healthCheck := handlers.NewHealthCheck()

	log.Info("configure server")
	srv, err := http.New(
		cfg.Servers.HTTP,
		cfg.Env,
		log,
		healthCheck,
	)
	if err != nil {
		log.Error("cannot create server", err)
		return
	}

	log.Info("start http server")
	errChan := make(chan error, 1)
	go func(errChan chan error) {
		if err := srv.Listen(ctx); err != nil {
			errChan <- fmt.Errorf("cannot run server: %w", err)
		}
	}(errChan)

	log.Info("started")
	for {
		select {
		case <-ctx.Done():
			log.Info("shutting down")
			return
		case err := <-errChan:
			log.Error("error", err)
			return
		}
	}
}
