package app

import (
	"time"

	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/storage"
)

type Event struct {
	ID               string         `json:"id"`
	Title            string         `json:"title"`
	Description      *string        `json:"description"`
	UserID           int64          `json:"userId"`
	StartDate        time.Time      `json:"startDate"`
	EndDate          time.Time      `json:"endDate"`
	NotificationTime *time.Duration `json:"notificationTime"`
}

func (e Event) ToStorageModel() storage.Event {
	return storage.Event{
		ID:               e.ID,
		Title:            e.Title,
		Description:      e.Description,
		UserID:           e.UserID,
		StartDate:        e.StartDate,
		EndDate:          e.EndDate,
		NotificationTime: e.NotificationTime,
	}
}

type Events []Event

func (e Events) FromStorageModel(events []storage.Event) Events {
	for i := range events {
		e[i] = Event{
			ID:               events[i].ID,
			Title:            events[i].Title,
			Description:      events[i].Description,
			UserID:           events[i].UserID,
			StartDate:        events[i].StartDate,
			EndDate:          events[i].EndDate,
			NotificationTime: events[i].NotificationTime,
		}
	}
	return e
}
