package internalhttp

import (
	"context"
)

type Server struct {
	log Logger
	app Application
}

type Logger interface {
	Info(string)
	Warn(string)
	Error(string)
	Debug(string)
}

type Application interface {
	CreateEvent(ctx context.Context, id, title string) error
}

func NewServer(logger Logger, app Application) *Server {
	return &Server{
		log: logger,
		app: app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.log.Info("Server started")
	<-ctx.Done()
	s.log.Info("Server stopped")
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return nil
}

// TODO
