package internalhttp

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/app"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/config"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	conf       config.Server
	log        Logger
	httpServer *http.Server
}

func NewServer(config config.Server, logger Logger) *Server {
	return &Server{
		conf: config,
		log:  logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	conn, err := grpc.DialContext(
		context.Background(),
		s.conf.Grpc.Addr,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	err = pb.RegisterCalendarHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    s.conf.Http.Addr,
		Handler: gwmux,
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

	s.log.Info(fmt.Sprintf("listening: %s", s.conf.Http.Addr))

	return gwServer.ListenAndServe()
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
