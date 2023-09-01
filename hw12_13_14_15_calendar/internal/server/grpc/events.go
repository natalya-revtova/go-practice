package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/models"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/pkg/api/calendarpb"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) CreateEvent(ctx context.Context, req *calendarpb.CreateEventRequest) (*calendarpb.CreateEventResponse, error) { //nolint:lll
	log := s.log.With(slog.String("request_id", middleware.GetReqID(ctx)))

	if err := validateCreateRequest(req); err != nil {
		log.Error("Validate event", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	eventID, err := s.app.CreateEvent(ctx, toModelForCreate(req))
	if err != nil {
		log.Error("Create event", "error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &calendarpb.CreateEventResponse{Id: eventID}, nil
}

func toModelForCreate(event *calendarpb.CreateEventRequest) *models.Event {
	var description *string
	if event.GetDescription() != "" {
		tmp := event.GetDescription()
		description = &tmp
	}

	var notTime *time.Duration
	if event.GetNotificationTime().AsDuration() != 0 {
		tmp := event.GetNotificationTime().AsDuration()
		notTime = &tmp
	}

	return &models.Event{
		Title:            event.GetTitle(),
		Description:      description,
		UserID:           event.GetUserId(),
		StartDate:        event.GetStartDate().AsTime(),
		EndDate:          event.GetEndDate().AsTime(),
		NotificationTime: notTime,
	}
}

func validateCreateRequest(event *calendarpb.CreateEventRequest) error {
	if len(event.GetTitle()) == 0 {
		return errors.New("field title is empty")
	}
	if event.GetUserId() == 0 {
		return errors.New("field userId is empty")
	}
	if event.GetStartDate() == nil {
		return errors.New("field startDate is empty")
	}
	if event.GetEndDate() == nil {
		return errors.New("field endDate is empty")
	}
	return nil
}

func (s *Server) UpdateEvent(ctx context.Context, req *calendarpb.Event) (*emptypb.Empty, error) {
	log := s.log.With(slog.String("request_id", middleware.GetReqID(ctx)))

	if len(req.GetId()) == 0 {
		err := errors.New("field id is empty")
		log.Error("Validate event", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.app.UpdateEvent(ctx, toModelForUpdate(req)); err != nil {
		log.Error("Update event", "error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func toModelForUpdate(event *calendarpb.Event) *models.Event {
	var description *string
	if event.GetDescription() != "" {
		tmp := event.GetDescription()
		description = &tmp
	}

	var notTime *time.Duration
	if event.GetNotificationTime().AsDuration() != 0 {
		tmp := event.GetNotificationTime().AsDuration()
		notTime = &tmp
	}

	return &models.Event{
		ID:               event.GetId(),
		Title:            event.GetTitle(),
		Description:      description,
		UserID:           event.GetUserId(),
		StartDate:        event.GetStartDate().AsTime(),
		EndDate:          event.GetEndDate().AsTime(),
		NotificationTime: notTime,
	}
}

func (s *Server) DeleteEvent(ctx context.Context, req *calendarpb.DeleteEventRequest) (*emptypb.Empty, error) {
	log := s.log.With(slog.String("request_id", middleware.GetReqID(ctx)))

	if len(req.GetId()) == 0 {
		err := errors.New("field id is empty")
		log.Error("Validate event", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.app.DeleteEvent(ctx, req.Id); err != nil {
		log.Error("Delete event", "error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) GetEventsByDay(ctx context.Context, req *calendarpb.EventsRequestByDate) (*calendarpb.EventsResponse, error) { //nolint:lll
	log := s.log.With(slog.String("request_id", middleware.GetReqID(ctx)))

	if err := validateRequestByDate(req); err != nil {
		log.Error("Validate request by day", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	events, err := s.app.GetEventByDay(ctx, req.GetUserId(), req.GetStartDate().AsTime())
	if err != nil {
		log.Error("Can not get events for the selected day",
			"user_id", req.GetUserId(),
			"day", req.GetStartDate().AsTime(),
			"error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toProtoEvents(events), nil
}

func (s *Server) GetEventsByWeek(ctx context.Context, req *calendarpb.EventsRequestByDate) (*calendarpb.EventsResponse, error) { //nolint:lll
	log := s.log.With(slog.String("request_id", middleware.GetReqID(ctx)))

	if err := validateRequestByDate(req); err != nil {
		log.Error("Validate request by day", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	events, err := s.app.GetEventByWeek(ctx, req.GetUserId(), req.GetStartDate().AsTime())
	if err != nil {
		log.Error("Can not get events for the selected week",
			"user_id", req.GetUserId(),
			"week", req.GetStartDate().AsTime(),
			"error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toProtoEvents(events), nil
}

func (s *Server) GetEventsByMonth(ctx context.Context, req *calendarpb.EventsRequestByDate) (*calendarpb.EventsResponse, error) { //nolint:lll
	log := s.log.With(slog.String("request_id", middleware.GetReqID(ctx)))

	if err := validateRequestByDate(req); err != nil {
		log.Error("Validate request by day", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	events, err := s.app.GetEventByMonth(ctx, req.GetUserId(), req.GetStartDate().AsTime())
	if err != nil {
		log.Error("Can not get events for the selected month",
			"user_id", req.GetUserId(),
			"month", req.GetStartDate().AsTime(),
			"error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toProtoEvents(events), nil
}

func validateRequestByDate(req *calendarpb.EventsRequestByDate) error {
	if req.GetUserId() == 0 {
		return errors.New("field userId is empty")
	}
	if req.GetStartDate() == nil {
		return errors.New("field startDate is empty")
	}
	return nil
}

func toProtoEvents(events []models.Event) *calendarpb.EventsResponse {
	pbEvents := make([]*calendarpb.Event, len(events))
	for i := range events {
		var description string
		if events[i].Description != nil {
			description = *events[i].Description
		}

		var notTime time.Duration
		if events[i].NotificationTime != nil {
			notTime = *events[i].NotificationTime
		}

		pbEvents[i] = &calendarpb.Event{
			Id:               events[i].ID,
			Title:            events[i].Title,
			Description:      description,
			UserId:           events[i].UserID,
			StartDate:        timestamppb.New(events[i].StartDate),
			EndDate:          timestamppb.New(events[i].EndDate),
			NotificationTime: durationpb.New(notTime),
		}
	}

	return &calendarpb.EventsResponse{Events: pbEvents}
}
