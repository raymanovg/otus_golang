package grpc

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	userID, err := uuid.Parse(req.UserID)
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
		EventID: event.ID.String(),
	}, nil
}

func (s *Service) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	eventID, err := uuid.Parse(req.EventID)
	if err != nil {
		return nil, err
	}
	if err = s.app.DeleteEvent(ctx, eventID); err != nil {
		return nil, err
	}
	return &pb.DeleteEventResponse{}, nil
}

func (s *Service) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	eventID, err := uuid.Parse(req.EventID)
	if err != nil {
		return nil, err
	}
	event := app.Event{
		ID:    eventID,
		Title: req.Event.Title,
		Desc:  req.Event.Desc,
		Begin: req.Event.Begin.AsTime(),
		End:   req.Event.End.AsTime(),
	}
	if err = s.app.UpdateEvent(ctx, event); err != nil {
		return nil, err
	}
	return &pb.UpdateEventResponse{EventID: req.EventID}, nil
}

func (s *Service) GetAllEventsOfUser(ctx context.Context, req *pb.GetAllEventsOfUserRequest) (*pb.GetAllEventsOfUserResponse, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, err
	}

	events, err := s.app.GetAllEventsOfUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	respEvents := make([]*pb.FullEvent, 0, len(events))
	for _, event := range events {
		respEvents = append(respEvents, &pb.FullEvent{
			Id:     event.ID.String(),
			UserID: event.ID.String(),
			Title:  event.Title,
			Desc:   event.Desc,
			Begin:  &timestamppb.Timestamp{Seconds: int64(event.Begin.Second())},
			End:    &timestamppb.Timestamp{Seconds: int64(event.End.Second())},
		})
	}

	return &pb.GetAllEventsOfUserResponse{Event: respEvents}, nil
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
