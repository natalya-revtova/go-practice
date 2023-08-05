package app

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/storage"
)

var ErrTimeBusy = errors.New("time is busy")

type Logger interface {
	With(args ...any) Logger
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}

type Storage interface {
	CreateEvent(context.Context, storage.Event) error
	UpdateEvent(context.Context, storage.Event) error
	DeleteEvent(context.Context, string) error
	GetEventByDay(context.Context, int64, time.Time) ([]storage.Event, error)
	GetEventByWeek(context.Context, int64, time.Time) ([]storage.Event, error)
	GetEventByMonth(context.Context, int64, time.Time) ([]storage.Event, error)
}

type App struct {
	storage Storage
	log     Logger
}

func New(logger Logger, storage Storage) *App {
	return &App{
		storage: storage,
		log:     logger,
	}
}

func (a *App) CreateEvent(ctx context.Context, event Event) error {
	events, err := a.storage.GetEventByDay(ctx, event.UserID, getStartDate(event.StartDate))
	if err != nil && !errors.Is(err, storage.ErrNoEventsFound) {
		return err
	}

	for i := range events {
		if busy(event.StartDate, event.EndDate, events[i].StartDate, events[i].EndDate) {
			a.log.Warn(
				"The selected time is already busy",
				"event_id", events[i].ID,
				"start", events[i].StartDate,
				"end", events[i].EndDate)
			return ErrTimeBusy
		}
	}

	event.ID = generateEventID()

	if err := a.storage.CreateEvent(ctx, event.ToStorageModel()); err != nil {
		a.log.Error("Can not add event to storage", "event_id", event.ID, "error", err)
		return err
	}
	return nil
}

func getStartDate(eventDate time.Time) time.Time {
	return time.Date(eventDate.Year(), eventDate.Month(), eventDate.Day(), 0, 0, 0, 0, time.Local)
}

func busy(newStart, newEnd, oldStart, oldEnd time.Time) bool {
	return (newStart.After(oldStart) && newStart.Before(oldEnd)) ||
		(newEnd.After(oldStart) && newEnd.Before(oldEnd)) ||
		(newStart.Before(oldStart) && newEnd.After(oldEnd))
}

func (a *App) UpdateEvent(ctx context.Context, eventID string, event Event) error {
	event.ID = eventID

	if err := a.storage.UpdateEvent(ctx, event.ToStorageModel()); err != nil {
		a.log.Error("Can not update event", "event_id", event.ID, "error", err)
		return err
	}
	return nil
}

func (a *App) DeleteEvent(ctx context.Context, eventID string) error {
	if err := a.storage.DeleteEvent(ctx, eventID); err != nil {
		a.log.Error("Can not delete event", "event_id", eventID, "error", err)
		return err
	}
	return nil
}

func (a *App) GetEventByDay(ctx context.Context, userID int64, day time.Time) ([]Event, error) {
	events, err := a.storage.GetEventByDay(ctx, userID, day)
	if err != nil {
		a.log.Error("Can not get events for the selected day", "user_id", userID, "day", day, "error", err)
		return nil, err
	}
	result := make(Events, len(events))
	return result.FromStorageModel(events), nil
}

func (a *App) GetEventByWeek(ctx context.Context, userID int64, week time.Time) ([]Event, error) {
	events, err := a.storage.GetEventByWeek(ctx, userID, week)
	if err != nil {
		a.log.Error("Can not get events for the selected week", "user_id", userID, "week", week, "error", err)
		return nil, err
	}
	result := make(Events, len(events))
	return result.FromStorageModel(events), nil
}

func (a *App) GetEventByMonth(ctx context.Context, userID int64, month time.Time) ([]Event, error) {
	events, err := a.storage.GetEventByMonth(ctx, userID, month)
	if err != nil {
		a.log.Error("Can not get events for the selected month", "user_id", userID, "month", month, "error", err)
		return nil, err
	}
	result := make(Events, len(events))
	return result.FromStorageModel(events), nil
}

func generateEventID() string {
	return uuid.New().String()
}
