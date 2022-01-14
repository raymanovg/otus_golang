package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/app"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/proto/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
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

type Handler struct {
	app Application
	log Logger
}

func NewHandler(app Application, log Logger) Handler {
	return Handler{
		app: app,
		log: log,
	}
}

func (h Handler) CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
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

	err = h.app.CreateEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	return &pb.CreateEventResponse{
		EventID: event.ID.String(),
	}, nil
}

func (h Handler) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	eventID, err := uuid.Parse(req.EventID)
	if err != nil {
		return nil, err
	}
	if err = h.app.DeleteEvent(ctx, eventID); err != nil {
		return nil, err
	}
	return &pb.DeleteEventResponse{}, nil
}

func (h Handler) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
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
	if err = h.app.UpdateEvent(ctx, event); err != nil {
		return nil, err
	}
	return &pb.UpdateEventResponse{EventID: req.EventID}, nil
}

func (h Handler) GetAllEventsOfUser(ctx context.Context, req *pb.GetAllEventsOfUserRequest) (*pb.GetAllEventsOfUserResponse, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, err
	}

	events, err := h.app.GetAllEventsOfUser(ctx, userID)
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
