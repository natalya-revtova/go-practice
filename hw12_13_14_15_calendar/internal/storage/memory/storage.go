package memorystorage

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/storage"
)

type (
	ID     = string
	Events map[ID]storage.Event
	Dates  map[time.Time]map[ID]struct{}
)

type Storage struct {
	events Events
	dates  Dates
	mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{
		events: make(Events),
		dates:  make(Dates),
	}
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[event.ID] = event
	s.saveEventDate(event.StartDate, event.EndDate, event.ID)

	return nil
}

func (s *Storage) saveEventDate(start, end time.Time, eventID string) {
	startDate := getDate(start)
	endDate := getDate(end)

	startDate = startDate.AddDate(0, 0, -1) // if event starts and ends in the same day
	for startDate != endDate {
		startDate = startDate.AddDate(0, 0, 1)
		if _, ok := s.dates[startDate]; !ok {
			s.dates[startDate] = make(map[ID]struct{})
		}
		s.dates[startDate][eventID] = struct{}{}
	}
}

func getDate(eventDate time.Time) time.Time {
	return time.Date(eventDate.Year(), eventDate.Month(), eventDate.Day(), 0, 0, 0, 0, time.Local)
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.RLock()
	updated, ok := s.events[event.ID]
	if !ok {
		return storage.ErrEventNotExist
	}
	s.mu.RUnlock()

	updateEventFields(&updated, event)

	s.mu.Lock()
	s.events[event.ID] = updated
	s.mu.Unlock()

	return nil
}

func updateEventFields(updated *storage.Event, event storage.Event) {
	if event.Description != nil {
		updated.Description = event.Description
	}
	if len(event.Title) != 0 {
		updated.Title = event.Title
	}
	if !event.StartDate.IsZero() {
		updated.StartDate = event.StartDate
	}
	if !event.EndDate.IsZero() {
		updated.EndDate = event.EndDate
	}
	if event.NotificationTime != nil {
		updated.NotificationTime = event.NotificationTime
	}
}

func (s *Storage) DeleteEvent(ctx context.Context, eventID string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.RLock()
	deleted, ok := s.events[eventID]
	if !ok {
		return storage.ErrEventNotExist
	}
	s.mu.RUnlock()

	s.mu.Lock()
	s.deleteEventDate(deleted.StartDate, deleted.EndDate, eventID)
	delete(s.events, eventID)
	s.mu.Unlock()

	return nil
}

func (s *Storage) deleteEventDate(start, end time.Time, eventID string) {
	startDate := getDate(start)
	endDate := getDate(end)

	startDate = startDate.AddDate(0, 0, -1) // if event starts and ends in the same day
	for startDate != endDate {
		startDate = startDate.AddDate(0, 0, 1)
		delete(s.dates[startDate], eventID)
	}
	if len(s.dates[startDate]) == 0 {
		delete(s.dates, startDate)
	}
}

func (s *Storage) GetEventByDay(ctx context.Context, userID int64, date time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getEventByPeriod(ctx, userID, date, date)
}

func (s *Storage) GetEventByWeek(ctx context.Context, userID int64, date time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getEventByPeriod(ctx, userID, date, date.AddDate(0, 0, 7))
}

func (s *Storage) GetEventByMonth(ctx context.Context, userID int64, date time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getEventByPeriod(ctx, userID, date, date.AddDate(0, 1, 0))
}

func (s *Storage) getEventByPeriod(ctx context.Context, userID int64, start, end time.Time) ([]storage.Event, error) {
	dates := make(map[ID]struct{}, 0)

	start = start.AddDate(0, 0, -1) // for period == one day
	for start != end {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		start = start.AddDate(0, 0, 1)
		dayEvents, ok := s.dates[start]
		if !ok {
			continue
		}

		for id := range dayEvents {
			if _, ok := dates[id]; !ok {
				if s.events[id].UserID == userID {
					dates[id] = struct{}{}
				}
			}
		}
	}

	if len(dates) == 0 {
		return nil, nil
	}

	events := make([]storage.Event, 0, len(dates))
	for id := range dates {
		events = append(events, s.events[id])
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].StartDate.Before(events[j].StartDate)
	})
	return events, nil
}
