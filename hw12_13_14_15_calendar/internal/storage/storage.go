package storage

import (
	"errors"
	"time"

	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/models"
	"github.com/snabb/isoweek"
)

const (
	InMemory = "in-memory"
	SQL      = "sql"
)

var (
	ErrNoEventsFound = errors.New("no events found")
	ErrEventNotExist = errors.New("event does not exist")
)

func FillDates(event *models.Event) {
	if !event.StartDate.IsZero() {
		event.Day = getDay(event.StartDate)
		event.Week = getWeek(event.StartDate)
		event.Month = getMonth(event.StartDate)
	}
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
