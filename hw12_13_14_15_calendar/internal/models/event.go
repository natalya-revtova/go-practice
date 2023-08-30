package models

import (
	"time"

	"github.com/snabb/isoweek"
)

type Event struct {
	ID               string         `db:"id"`
	Title            string         `db:"title"`
	Description      *string        `db:"description"`
	UserID           int64          `db:"user_id"`
	StartDate        time.Time      `db:"start_date"`
	EndDate          time.Time      `db:"end_date"`
	NotificationTime *time.Duration `db:"notification_time"`

	Day   time.Time `db:"day"`
	Week  time.Time `db:"week"`
	Month time.Time `db:"month"`
}

func (e *Event) UpdateForStorage() *Event {
	if !e.StartDate.IsZero() {
		e.Day = getDay(e.StartDate)
		e.Week = getWeek(e.StartDate)
		e.Month = getMonth(e.StartDate)
	}
	return e
}

func getDay(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
}

func getWeek(date time.Time) time.Time {
	year, week := date.ISOWeek()
	return isoweek.StartTime(year, week, time.UTC)
}

func getMonth(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
}
