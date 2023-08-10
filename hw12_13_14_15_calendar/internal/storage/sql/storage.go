package sqlstorage

import (
	"context"
	"database/sql"
	"time"

	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/storage"
)

type DB interface {
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type Storage struct {
	db DB
}

func New(db DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	query := `
	INSERT INTO events(id, title, description, user_id, start_date, end_date, day, week, month, notification_time)
	VALUES (:id, :title, :description, :user_id, :start_date, :end_date, :day, :week, :month, :notification_time)`

	_, err := s.db.NamedExecContext(ctx, query, event)
	return err
}

func (s *Storage) DeleteEvent(ctx context.Context, eventID string) error {
	query := `
	DELETE FROM events
	WHERE id = $1`

	_, err := s.db.ExecContext(ctx, query, eventID)
	return err
}

func (s *Storage) GetEventByDay(ctx context.Context, userID int64, day time.Time) ([]storage.Event, error) {
	query := `
	SELECT id, title, description, user_id, start_date, end_date, day, week, month, notification_time 
	FROM events
	WHERE user_id = $1 AND day = $2`

	var events []storage.Event
	return events, s.db.SelectContext(ctx, &events, query, userID, day)
}

func (s *Storage) GetEventByWeek(ctx context.Context, userID int64, week time.Time) ([]storage.Event, error) {
	query := `
	SELECT id, title, description, user_id, start_date, end_date, day, week, month, notification_time 
	FROM events
	WHERE user_id = $1 AND week = $2`

	var events []storage.Event
	return events, s.db.SelectContext(ctx, &events, query, userID, week)
}

func (s *Storage) GetEventByMonth(ctx context.Context, userID int64, month time.Time) ([]storage.Event, error) {
	query := `
	SELECT id, title, description, user_id, start_date, end_date, day, week, month, notification_time 
	FROM events
	WHERE user_id = $1 AND month = $2`

	var events []storage.Event
	return events, s.db.SelectContext(ctx, &events, query, userID, month)
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	query := buildUpdateQuery(event)
	_, err := s.db.NamedExecContext(ctx, query, event)
	return err
}

func buildUpdateQuery(event storage.Event) string {
	qb := NewUpdateQueryBuilder("events")

	qb.SetIf(event.Title != "", "title = :title")
	qb.SetIf(event.Description != nil, "description = :description")
	qb.SetIf(!event.StartDate.IsZero(), "start_date = :start_date")
	qb.SetIf(!event.EndDate.IsZero(), "end_date = :end_date")
	qb.SetIf(!event.Day.IsZero(), "day = :day")
	qb.SetIf(!event.Week.IsZero(), "week = :week")
	qb.SetIf(!event.Month.IsZero(), "month = :month")
	qb.SetIf(event.NotificationTime != nil, "notification_time = :notification_time")

	qb.Where("id = :id")
	return qb.Build()
}
