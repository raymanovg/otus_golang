package internalhttp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/app"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/server/http/middleware"
)

type Server struct {
	addr string
	log  Logger
	app  Application
}

type Logger interface {
	Info(string)
	Warn(string)
	Error(string)
	Debug(string)
	Write(string)
}

type Application interface {
	CreateEvent(ctx context.Context, event app.Event) error
	GetAllEvents(ctx context.Context, userID string) ([]app.Event, error)
}

func NewServer(addr string, logger Logger, app Application) *Server {
	return &Server{
		addr: addr,
		log:  logger,
		app:  app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "hello\n")
		time.Sleep(time.Second)
	})

	server := &http.Server{
		Addr:    s.addr,
		Handler: middleware.NewLoggerMiddleware(s.log, mux),
	}

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}()

	s.log.Info(fmt.Sprintf("Listening %s", s.addr))

	return server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return nil
}
