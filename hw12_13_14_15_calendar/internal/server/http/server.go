package internalhttp

import (
	"context"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"

	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/app"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/server/http/middleware"
)

type Server struct {
	log        *zap.Logger
	app        Application
	httpServer *http.Server
}

type Application interface {
	CreateEvent(ctx context.Context, event app.Event) error
	GetAllEvents(ctx context.Context, userID string) ([]app.Event, error)
}

func NewServer(logger *zap.Logger, app Application) *Server {
	return &Server{
		log: logger,
		app: app,
	}
}

func (s *Server) Start(addr string, ctx context.Context) error {
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: handler(s.log),
	}

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := s.Stop(ctx)
		if err != nil {
			s.log.Error("failed to stop server")
		}
	}()

	s.log.Info("listening", zap.String("addr", addr))

	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func handler(logger *zap.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "hello\n")
	})

	return middleware.NewLoggerMiddleware(logger, mux)
}
