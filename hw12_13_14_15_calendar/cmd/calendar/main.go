package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/app"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/config"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/handler"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/logger"
	grpcServer "github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/server/grpc"
	httpServer "github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/server/http"
	memoryStorage "github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/storage/memory"
	sqlStorage "github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "config.yaml", "config file")
}

func main() {
	flag.Parse()
	conf, err := config.NewConfig(configFile)
	if err != nil {
		fmt.Printf("failed to init conf: %s\n", err)
		os.Exit(1)
	}
	storage, err := getStorage(conf.App.Storage)
	if err != nil {
		fmt.Printf("failed to get storage: %s\n", err)
		os.Exit(1)
	}
	log, err := logger.NewZapLogger(conf.Logger)
	if err != nil {
		fmt.Printf("failed to get logger: %s\n", err)
		os.Exit(1)
	}

	if err := storage.Connect(context.Background()); err != nil {
		fmt.Printf("failed to connect to storage: %s\n", err)
		os.Exit(1)
	}
	defer storage.Close()

	calendar := app.New(log, storage)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	grpc := grpcServer.NewServer(conf.Server.Grpc, log, handler.NewHandler(calendar, log))
	gateway := httpServer.NewServer(conf.Server, log)

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		grpc.Stop()
		if err := gateway.Stop(ctx); err != nil {
			log.Error("failed to stop http server: " + err.Error())
		}
	}()

	go func() {
		err := grpc.Start(ctx)
		if err != nil {
			cancel()
		}
	}()

	if err := gateway.Start(ctx); err != nil {
		log.Error("failed to start http server: ", err.Error())
		cancel()
	}

}

func getStorage(config config.Storage) (app.Storage, error) {
	if config.Name == "sql" {
		return sqlStorage.New(config.SQL), nil
	}
	if config.Name == "memory" {
		return memoryStorage.New(config.Memory), nil
	}
	return nil, errors.New("unknown storage: " + config.Name)
}
