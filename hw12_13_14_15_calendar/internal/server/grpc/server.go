package grpc

import (
	"context"
	"net"

	"github.com/google/uuid"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/app"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/config"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"google.golang.org/grpc"
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
	grpc *grpc.Server
	conf config.GrpcServerConf
	log  Logger
}

func (s *Server) Start() error {
	lsn, err := net.Listen("tcp", s.conf.Addr)
	if err != nil {
		return err
	}

	s.log.Info("starting server on %s", lsn.Addr().String())

	return s.grpc.Serve(lsn)
}

func (s *Server) Stop() {
	s.grpc.Stop()
}

type Service struct {
	pb.UnimplementedCalendarServer
	app Application
}

func (s *Service) CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	userID, err := uuid.Parse(req.Event.UserID)
	if err != nil {
		return nil, err
	}

	event := app.Event{
		ID:     uuid.New(),
		UserID: userID,
		Title:  req.Event.Title,
		Desc:   req.Event.Desc,
		Begin:  req.Event.Begin.AsTime(),
		End:    req.Event.End.AsTime(),
	}

	err = s.app.CreateEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	return &pb.CreateEventResponse{
		Ok:      true,
		EventId: event.ID.String(),
	}, nil
}

func NewServer(conf config.GrpcServerConf, log Logger, app Application) Server {
	service := new(Service)
	service.app = app
	server := grpc.NewServer()
	pb.RegisterCalendarServer(server, service)
	return Server{
		grpc: server,
		conf: conf,
		log:  log,
	}
}
