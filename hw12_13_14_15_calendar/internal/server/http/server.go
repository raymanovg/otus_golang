package internalhttp

import (
	"context"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/config"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"

	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/app"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/server/http/middleware"
)

type Server struct {
	conf       config.ServerConf
	log        *zap.Logger
	app        Application
	httpServer *http.Server
}

type Application interface {
	CreateEvent(ctx context.Context, event app.Event) error
	DeleteEvent(ctx context.Context, eventID int64) error
	UpdateEvent(ctx context.Context, event app.Event) error
	GetAllEvents(ctx context.Context) ([]app.Event, error)
	GetAllEventsOfUser(ctx context.Context, userID int64) ([]app.Event, error)
}

func NewServer(config config.ServerConf, logger *zap.Logger, app Application) *Server {
	return &Server{
		conf: config,
		log:  logger,
		app:  app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.httpServer = &http.Server{
		Addr:    s.conf.Addr,
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

	s.log.Info("listening", zap.String("addr", s.conf.Addr))

	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func handler(logger *zap.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "hello\n")
	})

	return middleware.NewLoggerMiddleware(logger, mux)
}
