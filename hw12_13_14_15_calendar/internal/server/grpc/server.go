package grpc

import (
	"context"
	"net"

	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/config"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/proto/pb"
	"google.golang.org/grpc"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type Handler interface {
	CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error)
	DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error)
	UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error)
	GetAllEventsOfUser(ctx context.Context, req *pb.GetAllEventsOfUserRequest) (*pb.GetAllEventsOfUserResponse, error)
}

func NewServer(conf config.GrpcServerConf, log Logger, handler Handler) Server {
	return Server{
		conf:    conf,
		log:     log,
		handler: handler,
		grpc:    grpc.NewServer(),
	}
}

type Server struct {
	conf    config.GrpcServerConf
	log     Logger
	handler Handler
	grpc    *grpc.Server
	pb.UnimplementedCalendarServer
}

func (s *Server) Start(ctx context.Context) error {
	lsn, err := net.Listen("tcp", s.conf.Addr)
	if err != nil {
		return err
	}
	pb.RegisterCalendarServer(s.grpc, s)
	go func() {
		<-ctx.Done()
		s.Stop()
	}()

	s.log.Info("starting grpc server on %s", lsn.Addr().String())

	return s.grpc.Serve(lsn)
}

func (s *Server) Stop() {
	s.grpc.GracefulStop()
}

func (s *Server) CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	return s.handler.CreateEvent(ctx, req)
}

func (s *Server) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	return s.handler.DeleteEvent(ctx, req)
}

func (s *Server) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	return s.handler.UpdateEvent(ctx, req)
}

func (s *Server) GetAllEventsOfUser(ctx context.Context, req *pb.GetAllEventsOfUserRequest) (*pb.GetAllEventsOfUserResponse, error) {
	return s.handler.GetAllEventsOfUser(ctx, req)
}
