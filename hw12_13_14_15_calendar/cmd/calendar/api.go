package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/app"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/storage/memory"

	"github.com/spf13/cobra"
)

var configFile string

func init() {
	api.Flags().StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

var api = &cobra.Command{
	Use:   "api",
	Short: "calendar api",
	Run: func(cmd *cobra.Command, args []string) {
		config := NewConfig()
		logg := logger.New(os.Stdout, config.Logger.Level)

		storage := memorystorage.New()
		calendar := app.New(logg, storage)

		server := internalhttp.NewServer(logg, calendar)

		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		defer cancel()

		go func() {
			<-ctx.Done()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()

			if err := server.Stop(ctx); err != nil {
				logg.Error("failed to stop http server: " + err.Error())
			}
		}()

		logg.Info("calendar api is running...")

		if err := server.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
			os.Exit(1) //nolint:gocritic
		}
	},
}
