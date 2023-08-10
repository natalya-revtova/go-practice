package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/storage"
	"github.com/snabb/isoweek"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateEvent(t *testing.T) {
	memoryStorage := New()

	newEvents := generateEvents(time.Date(2010, 1, 1, 13, 0, 0, 0, time.UTC),
		time.Date(2010, 1, 1, 15, 0, 0, 0, time.UTC), 1)
	wantEvents := getEvents(newEvents)
	wantDays, wantWeeks, wantMonths := getDates(newEvents)

	err := memoryStorage.CreateEvent(context.Background(), newEvents[0])
	require.NoError(t, err)

	assert.Equal(t, wantEvents, memoryStorage.events)
	assert.Equal(t, wantDays, memoryStorage.days)
	assert.Equal(t, wantWeeks, memoryStorage.weeks)
	assert.Equal(t, wantMonths, memoryStorage.months)
}

func TestDeleteEvent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		memoryStorage := New()

		newEvents := generateEvents(time.Date(2010, 1, 1, 13, 0, 0, 0, time.UTC),
			time.Date(2010, 1, 1, 15, 0, 0, 0, time.UTC), 1)
		wantEvents := getEvents(nil)
		wantDays, wantWeeks, wantMonths := getDates(nil)

		err := memoryStorage.CreateEvent(context.Background(), newEvents[0])
		require.NoError(t, err)

		err = memoryStorage.DeleteEvent(context.Background(), newEvents[0].ID)
		require.NoError(t, err)

		assert.Equal(t, wantEvents, memoryStorage.events)
		assert.Equal(t, wantDays, memoryStorage.days)
		assert.Equal(t, wantWeeks, memoryStorage.weeks)
		assert.Equal(t, wantMonths, memoryStorage.months)
	})

	t.Run("event does not exist", func(t *testing.T) {
		memoryStorage := New()

		err := memoryStorage.DeleteEvent(context.Background(), "id")
		require.ErrorIs(t, err, storage.ErrEventNotExist)
	})
}

func TestUpdateEvent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		memoryStorage := New()

		newEvents := generateEvents(time.Date(2010, 1, 1, 13, 0, 0, 0, time.UTC),
			time.Date(2010, 1, 1, 15, 0, 0, 0, time.UTC), 1)

		err := memoryStorage.CreateEvent(context.Background(), newEvents[0])
		require.NoError(t, err)
		require.Equal(t, newEvents[0].Title, memoryStorage.events[newEvents[0].ID].Title)

		updatedEvent := storage.Event{
			ID:    newEvents[0].ID,
			Title: "new title",
		}

		err = memoryStorage.UpdateEvent(context.Background(), updatedEvent)
		require.NoError(t, err)
		require.Equal(t, updatedEvent.Title, memoryStorage.events[newEvents[0].ID].Title)
	})

	t.Run("event does not exist", func(t *testing.T) {
		memoryStorage := New()

		err := memoryStorage.UpdateEvent(context.Background(), storage.Event{ID: "id"})
		require.ErrorIs(t, err, storage.ErrEventNotExist)
	})
}

func TestGetEventByDay(t *testing.T) {
	memoryStorage := New()

	newEvents := generateEvents(time.Date(2010, 1, 1, 13, 0, 0, 0, time.UTC),
		time.Date(2010, 1, 1, 15, 0, 0, 0, time.UTC), 2)

	err := memoryStorage.CreateEvent(context.Background(), newEvents[0])
	require.NoError(t, err)

	err = memoryStorage.CreateEvent(context.Background(), newEvents[1])
	require.NoError(t, err)

	got, err := memoryStorage.GetEventByDay(context.Background(), newEvents[0].UserID,
		time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	assert.Equal(t, []storage.Event{newEvents[0]}, got)
}

func TestGetEventByWeek(t *testing.T) {
	memoryStorage := New()

	newEvents := generateEvents(time.Date(2010, 1, 1, 13, 0, 0, 0, time.UTC),
		time.Date(2010, 1, 1, 15, 0, 0, 0, time.UTC), 2)

	err := memoryStorage.CreateEvent(context.Background(), newEvents[0])
	require.NoError(t, err)

	err = memoryStorage.CreateEvent(context.Background(), newEvents[1])
	require.NoError(t, err)

	got, err := memoryStorage.GetEventByWeek(context.Background(), newEvents[0].UserID,
		time.Date(2009, 12, 28, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)

	require.Equal(t, newEvents, got)
}

func TestGetEventByMonth(t *testing.T) {
	memoryStorage := New()

	newEvents := generateEvents(time.Date(2010, 1, 1, 13, 0, 0, 0, time.UTC),
		time.Date(2010, 1, 1, 15, 0, 0, 0, time.UTC), 2)

	err := memoryStorage.CreateEvent(context.Background(), newEvents[0])
	require.NoError(t, err)

	err = memoryStorage.CreateEvent(context.Background(), newEvents[1])
	require.NoError(t, err)

	got, err := memoryStorage.GetEventByMonth(context.Background(), newEvents[0].UserID,
		time.Date(2010, 1, 0, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)

	require.Equal(t, newEvents, got)
}

func getDates(events []storage.Event) (dates, dates, dates) {
	days := make(dates)
	weeks := make(dates)
	months := make(dates)

	for i := range events {
		days[events[i].Day] = make(map[id]struct{})
		days[events[i].Day][events[i].ID] = struct{}{}

		weeks[events[i].Week] = make(map[id]struct{})
		weeks[events[i].Week][events[i].ID] = struct{}{}

		months[events[i].Month] = make(map[id]struct{})
		months[events[i].Month][events[i].ID] = struct{}{}
	}

	return days, weeks, months
}

func getEvents(eventsIn []storage.Event) events {
	result := make(events)
	for i := range eventsIn {
		result[eventsIn[i].ID] = eventsIn[i]
	}

	return result
}

func generateEvents(start, end time.Time, count int) []storage.Event {
	events := make([]storage.Event, 0, count)

	year, week := start.ISOWeek()
	for i := 0; i < count; i++ {
		events = append(events, storage.Event{
			ID:               uuid.New().String(),
			Title:            "some title",
			Description:      nil,
			UserID:           1,
			StartDate:        start,
			EndDate:          end,
			NotificationTime: nil,
			Day:              time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC),
			Week:             isoweek.StartTime(year, week, time.UTC),
			Month:            time.Date(start.Year(), start.Month(), 0, 0, 0, 0, 0, time.UTC),
		})

		start = start.AddDate(0, 0, 1)
		end = end.AddDate(0, 0, 1)
	}

	return events
}
