package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/app"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/config"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type Application interface {
	CreateEvent(ctx context.Context, event app.Event) error
	DeleteEvent(ctx context.Context, eventID uuid.UUID) error
	UpdateEvent(ctx context.Context, event app.Event) error
	GetAllEventsOfUser(ctx context.Context, userID uuid.UUID) ([]app.Event, error)
}

type Server struct {
	conf       config.ServerConf
	log        Logger
	app        Application
	httpServer *http.Server
}

func NewServer(config config.ServerConf, logger Logger, app Application) *Server {
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

	s.log.Info(fmt.Sprintf("listening: %s", s.conf.Addr))

	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func handler(logger Logger) http.Handler {
	router := mux.NewRouter()
	router.Use(newLoggingMiddleware(logger))
	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(200)
		rw.Write([]byte("Hello user"))
	})

	return router
}
