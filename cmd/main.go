package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"

	"go.uber.org/zap"

	"github.com/Vivi-social-network/core/logger"
	"github.com/Vivi-social-network/gateway/internal/config"
	"github.com/Vivi-social-network/gateway/internal/server/http"
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

	log.Info("service starting", zap.String("env", cfg.Env), zap.Int("CPUs", runtime.NumCPU()))

	log.Info("configure server")
	srv, err := http.New(
		cfg.Servers.HTTP,
		cfg.IsDev(),
		log,
	)
	if err != nil {
		log.Fatal("cannot create server", zap.Error(err))
	}

	log.Info("start http server")
	go func() {
		if err := srv.Listen(ctx); err != nil {
			log.Fatal("cannot run server", zap.Error(err))
		}
	}()

	log.Info("started")
	<-ctx.Done()
}
