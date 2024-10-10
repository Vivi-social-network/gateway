package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

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

	logger := createLogger(cfg.Logger.Level)

	logger.Info("service starting", zap.String("env", cfg.Env), zap.Int("CPUs", runtime.NumCPU()))

	logger.Info("configure server")
	srv, err := http.New(
		cfg.Servers.HTTP,
		cfg.IsDev(),
		logger,
	)
	if err != nil {
		logger.Fatal("cannot create server", zap.Error(err))
	}

	logger.Info("start http server")
	go func() {
		if err := srv.Listen(ctx); err != nil {
			logger.Fatal("cannot run server", zap.Error(err))
		}
	}()

	logger.Info("started")
	<-ctx.Done()
}

func createLogger(logLevel int8) *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.Level(logLevel)),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
	}

	return zap.Must(zapConfig.Build())
}
