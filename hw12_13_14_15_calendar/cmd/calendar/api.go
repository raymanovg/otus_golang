package main

import (
	"context"
	"fmt"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/logger"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/app"
	internalhttp "github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/storage/memory"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var configFile string

func init() {
	api.Flags().StringVar(&configFile, "config", "./calendar-config.yaml", "Path to configuration file")
}

var api = &cobra.Command{
	Use:   "api",
	Short: "calendar api",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := NewConfig(configFile)
		if err != nil {
			panic(fmt.Sprintf("failed to init conf: %s", err))
		}

		log.Fatal(config.Logger.Output)

		logger, err := logger.NewZapLogger(config.Logger.Output, config.Logger.Level, config.Logger.DevMode)

		storage := memorystorage.New()
		calendar := app.New(logger, storage)
		server := internalhttp.NewServer(logger, calendar)

		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		defer cancel()

		go func() {
			<-ctx.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()

			if err := server.Stop(ctx); err != nil {
				logger.Error("failed to stop http server: " + err.Error())
			}
		}()

		if err := server.Start(config.Server.Addr, ctx); err != nil {
			logger.Error("failed to start http server", zap.Error(err))
		}
	},
}
