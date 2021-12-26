package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/app"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/config"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/logger"
	httpServer "github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/server/http"
	memoryStorage "github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/storage/memory"
	"go.uber.org/zap"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "config.yaml", "config file")
}

func main() {
	flag.Parse()
	conf, err := config.NewConfig(configFile)
	if err != nil {
		fmt.Printf("failed to init conf: %s", err)
		os.Exit(1)
	}

	zapLogger, err := logger.NewZapLogger(conf.Logger)
	storage := memoryStorage.New()
	calendar := app.New(zapLogger, storage)
	server := httpServer.NewServer(conf.Server, zapLogger, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			zapLogger.Error("failed to stop http server: " + err.Error())
		}
	}()

	if err := server.Start(ctx); err != nil {
		zapLogger.Error("failed to start http server", zap.Error(err))
	}
}
