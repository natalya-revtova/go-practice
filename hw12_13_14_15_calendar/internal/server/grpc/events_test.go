package grpc

import (
	"context"
	"errors"
	"log"
	"net"
	"testing"
	"time"

	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/logger"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/models"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/server/mocks"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/storage"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/pkg/api/calendarpb"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func startServer(t *testing.T) (*mocks.Calendar, calendarpb.CalendarClient, func()) {
	t.Helper()

	buffer := 1024 * 1024
	lis := bufconn.Listen(buffer)

	appMock := mocks.NewCalendar(t)

	calendarSrv := &Server{
		app: appMock,
		log: logger.NewMock(),
	}

	calendarSrv.srv = grpc.NewServer()
	calendarpb.RegisterCalendarServer(calendarSrv.srv, calendarSrv)

	go func() {
		if err := calendarSrv.srv.Serve(lis); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()

	conn, err := grpc.DialContext(context.Background(), "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	closeFn := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		calendarSrv.srv.Stop()
	}

	return appMock, calendarpb.NewCalendarClient(conn), closeFn
}

func TestCreateEvent(t *testing.T) {
	appMock, client, closeConn := startServer(t)
	defer closeConn()

	cases := []struct {
		name          string
		event         *calendarpb.CreateEventRequest
		validateError error
		mockError     error
		code          codes.Code
	}{
		{
			name: "success",
			event: &calendarpb.CreateEventRequest{
				Title:            "test",
				Description:      "test",
				UserId:           1,
				StartDate:        timestamppb.New(time.Date(2023, 8, 16, 12, 0, 0, 0, time.UTC)),
				EndDate:          timestamppb.New(time.Date(2023, 8, 16, 13, 0, 0, 0, time.UTC)),
				NotificationTime: durationpb.New(5 * time.Second),
			},
		},
		{
			name: "empty title",
			event: &calendarpb.CreateEventRequest{
				Description:      "test",
				UserId:           1,
				StartDate:        timestamppb.New(time.Date(2023, 8, 16, 12, 0, 0, 0, time.UTC)),
				EndDate:          timestamppb.New(time.Date(2023, 8, 16, 13, 0, 0, 0, time.UTC)),
				NotificationTime: durationpb.New(5 * time.Second),
			},
			validateError: errors.New("field title is empty"),
			code:          codes.InvalidArgument,
		},
		{
			name: "empty userId",
			event: &calendarpb.CreateEventRequest{
				Title:            "test",
				Description:      "test",
				StartDate:        timestamppb.New(time.Date(2023, 8, 16, 12, 0, 0, 0, time.UTC)),
				EndDate:          timestamppb.New(time.Date(2023, 8, 16, 13, 0, 0, 0, time.UTC)),
				NotificationTime: durationpb.New(5 * time.Second),
			},
			validateError: errors.New("field userId is empty"),
			code:          codes.InvalidArgument,
		},
		{
			name: "empty start_ate",
			event: &calendarpb.CreateEventRequest{
				Title:            "test",
				Description:      "test",
				UserId:           1,
				EndDate:          timestamppb.New(time.Date(2023, 8, 16, 13, 0, 0, 0, time.UTC)),
				NotificationTime: durationpb.New(5 * time.Second),
			},
			validateError: errors.New("field startDate is empty"),
			code:          codes.InvalidArgument,
		},
		{
			name: "empty endDate",
			event: &calendarpb.CreateEventRequest{
				Title:            "test",
				Description:      "test",
				UserId:           1,
				StartDate:        timestamppb.New(time.Date(2023, 8, 16, 12, 0, 0, 0, time.UTC)),
				NotificationTime: durationpb.New(5 * time.Second),
			},
			validateError: errors.New("field endDate is empty"),
			code:          codes.InvalidArgument,
		},
		{
			name: "create event error",
			event: &calendarpb.CreateEventRequest{
				Title:            "test",
				Description:      "test",
				UserId:           1,
				StartDate:        timestamppb.New(time.Date(2023, 8, 16, 12, 0, 0, 0, time.UTC)),
				EndDate:          timestamppb.New(time.Date(2023, 8, 16, 13, 0, 0, 0, time.UTC)),
				NotificationTime: durationpb.New(5 * time.Second),
			},
			mockError: errors.New("unexpected error"),
			code:      codes.Internal,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.validateError == nil || tc.mockError != nil {
				appMock.On("CreateEvent", mock.Anything, toModelForCreate(tc.event)).
					Return(tc.mockError).
					Once()
			}

			_, err := client.CreateEvent(context.Background(), tc.event)

			switch {
			case tc.mockError != nil:
				require.Equal(t, status.Error(tc.code, tc.mockError.Error()), err)

			case tc.validateError != nil:
				require.Equal(t, status.Error(tc.code, tc.validateError.Error()), err)

			default:
				require.NoError(t, err)
			}
		})
	}
}

func TestUpdateEvent(t *testing.T) {
	appMock, client, closeConn := startServer(t)
	defer closeConn()

	cases := []struct {
		name          string
		event         *calendarpb.Event
		validateError error
		mockError     error
		code          codes.Code
	}{
		{
			name: "success",
			event: &calendarpb.Event{
				Id:               "id-1",
				NotificationTime: durationpb.New(5 * time.Second),
			},
		},
		{
			name: "empty id",
			event: &calendarpb.Event{
				NotificationTime: durationpb.New(5 * time.Second),
			},
			validateError: errors.New("field id is empty"),
			code:          codes.InvalidArgument,
		},
		{
			name: "update event error",
			event: &calendarpb.Event{
				Id:               "id-1",
				NotificationTime: durationpb.New(5 * time.Second),
			},
			mockError: storage.ErrEventNotExist,
			code:      codes.Internal,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.validateError == nil || tc.mockError != nil {
				appMock.On("UpdateEvent", mock.Anything, toModelForUpdate(tc.event)).
					Return(tc.mockError).
					Once()
			}

			_, err := client.UpdateEvent(context.Background(), tc.event)

			switch {
			case tc.mockError != nil:
				require.Equal(t, status.Error(tc.code, tc.mockError.Error()), err)

			case tc.validateError != nil:
				require.Equal(t, status.Error(tc.code, tc.validateError.Error()), err)

			default:
				require.NoError(t, err)
			}
		})
	}
}

func TestDeleteEvent(t *testing.T) {
	appMock, client, closeConn := startServer(t)
	defer closeConn()

	cases := []struct {
		name          string
		request       *calendarpb.DeleteEventRequest
		validateError error
		mockError     error
		code          codes.Code
	}{
		{
			name: "success",
			request: &calendarpb.DeleteEventRequest{
				Id: "id-1",
			},
		},
		{
			name:          "empty id",
			request:       &calendarpb.DeleteEventRequest{},
			validateError: errors.New("field id is empty"),
			code:          codes.InvalidArgument,
		},
		{
			name: "delete event error",
			request: &calendarpb.DeleteEventRequest{
				Id: "id-1",
			},
			mockError: storage.ErrEventNotExist,
			code:      codes.Internal,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.validateError == nil || tc.mockError != nil {
				appMock.On("DeleteEvent", mock.Anything, tc.request.GetId()).
					Return(tc.mockError).
					Once()
			}

			_, err := client.DeleteEvent(context.Background(), tc.request)

			switch {
			case tc.mockError != nil:
				require.Equal(t, status.Error(tc.code, tc.mockError.Error()), err)

			case tc.validateError != nil:
				require.Equal(t, status.Error(tc.code, tc.validateError.Error()), err)

			default:
				require.NoError(t, err)
			}
		})
	}
}

func TestGetEventsByDay(t *testing.T) {
	appMock, client, closeConn := startServer(t)
	defer closeConn()

	cases := []struct {
		name          string
		request       *calendarpb.EventsRequestByDate
		date          time.Time
		events        []models.Event
		validateError error
		mockError     error
		code          codes.Code
	}{
		{
			name: "success",
			request: &calendarpb.EventsRequestByDate{
				UserId:    1,
				StartDate: timestamppb.New(time.Date(2023, 8, 16, 0, 0, 0, 0, time.UTC)),
			},
			date: time.Date(2023, 8, 16, 0, 0, 0, 0, time.UTC),
			events: []models.Event{
				{
					Title:            "test",
					Description:      nil,
					UserID:           1,
					StartDate:        time.Date(2023, 8, 16, 12, 0, 0, 0, time.UTC),
					EndDate:          time.Date(2023, 8, 16, 13, 0, 0, 0, time.UTC),
					NotificationTime: nil,
					Day:              time.Date(2023, 8, 16, 0, 0, 0, 0, time.UTC),
					Week:             time.Date(2023, 8, 14, 0, 0, 0, 0, time.UTC),
					Month:            time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "empty userId",
			request: &calendarpb.EventsRequestByDate{
				StartDate: timestamppb.New(time.Date(2023, 8, 16, 0, 0, 0, 0, time.UTC)),
			},
			validateError: errors.New("field userId is empty"),
			code:          codes.InvalidArgument,
		},
		{
			name: "empty startDate",
			request: &calendarpb.EventsRequestByDate{
				UserId: 1,
			},
			validateError: errors.New("field startDate is empty"),
			code:          codes.InvalidArgument,
		},
		{
			name: "get events error",
			request: &calendarpb.EventsRequestByDate{
				UserId:    1,
				StartDate: timestamppb.New(time.Date(2023, 8, 16, 0, 0, 0, 0, time.UTC)),
			},
			date:      time.Date(2023, 8, 16, 0, 0, 0, 0, time.UTC),
			mockError: errors.New("unexpected error"),
			code:      codes.Internal,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.validateError == nil || tc.mockError != nil {
				appMock.On("GetEventByDay", mock.Anything, tc.request.UserId, tc.date).
					Return(tc.events, tc.mockError).
					Once()
			}

			resp, err := client.GetEventsByDay(context.Background(), tc.request)

			switch {
			case tc.mockError != nil:
				require.Equal(t, status.Error(tc.code, tc.mockError.Error()), err)

			case tc.validateError != nil:
				require.Equal(t, status.Error(tc.code, tc.validateError.Error()), err)

			default:
				require.NoError(t, err)
				require.Equal(t, toProtoEvents(tc.events).GetEvents(), resp.GetEvents())
			}
		})
	}
}

func TestGetEventsByWeek(t *testing.T) {
	appMock, client, closeConn := startServer(t)
	defer closeConn()

	cases := []struct {
		name          string
		request       *calendarpb.EventsRequestByDate
		date          time.Time
		events        []models.Event
		validateError error
		mockError     error
		code          codes.Code
	}{
		{
			name: "success",
			request: &calendarpb.EventsRequestByDate{
				UserId:    1,
				StartDate: timestamppb.New(time.Date(2023, 8, 14, 0, 0, 0, 0, time.UTC)),
			},
			date: time.Date(2023, 8, 14, 0, 0, 0, 0, time.UTC),
			events: []models.Event{
				{
					Title:            "test",
					Description:      nil,
					UserID:           1,
					StartDate:        time.Date(2023, 8, 16, 12, 0, 0, 0, time.UTC),
					EndDate:          time.Date(2023, 8, 16, 13, 0, 0, 0, time.UTC),
					NotificationTime: nil,
					Day:              time.Date(2023, 8, 16, 0, 0, 0, 0, time.UTC),
					Week:             time.Date(2023, 8, 14, 0, 0, 0, 0, time.UTC),
					Month:            time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "empty userId",
			request: &calendarpb.EventsRequestByDate{
				StartDate: timestamppb.New(time.Date(2023, 8, 14, 0, 0, 0, 0, time.UTC)),
			},
			validateError: errors.New("field userId is empty"),
			code:          codes.InvalidArgument,
		},
		{
			name: "empty startDate",
			request: &calendarpb.EventsRequestByDate{
				UserId: 1,
			},
			validateError: errors.New("field startDate is empty"),
			code:          codes.InvalidArgument,
		},
		{
			name: "get events error",
			request: &calendarpb.EventsRequestByDate{
				UserId:    1,
				StartDate: timestamppb.New(time.Date(2023, 8, 14, 0, 0, 0, 0, time.UTC)),
			},
			date:      time.Date(2023, 8, 14, 0, 0, 0, 0, time.UTC),
			mockError: errors.New("unexpected error"),
			code:      codes.Internal,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.validateError == nil || tc.mockError != nil {
				appMock.On("GetEventByWeek", mock.Anything, tc.request.UserId, tc.date).
					Return(tc.events, tc.mockError).
					Once()
			}

			resp, err := client.GetEventsByWeek(context.Background(), tc.request)

			switch {
			case tc.mockError != nil:
				require.Equal(t, status.Error(tc.code, tc.mockError.Error()), err)

			case tc.validateError != nil:
				require.Equal(t, status.Error(tc.code, tc.validateError.Error()), err)

			default:
				require.NoError(t, err)
				require.Equal(t, toProtoEvents(tc.events).GetEvents(), resp.GetEvents())
			}
		})
	}
}

func TestGetEventsByMonth(t *testing.T) {
	appMock, client, closeConn := startServer(t)
	defer closeConn()

	cases := []struct {
		name          string
		request       *calendarpb.EventsRequestByDate
		date          time.Time
		events        []models.Event
		validateError error
		mockError     error
		code          codes.Code
	}{
		{
			name: "success",
			request: &calendarpb.EventsRequestByDate{
				UserId:    1,
				StartDate: timestamppb.New(time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC)),
			},
			date: time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
			events: []models.Event{
				{
					Title:            "test",
					Description:      nil,
					UserID:           1,
					StartDate:        time.Date(2023, 8, 16, 12, 0, 0, 0, time.UTC),
					EndDate:          time.Date(2023, 8, 16, 13, 0, 0, 0, time.UTC),
					NotificationTime: nil,
					Day:              time.Date(2023, 8, 16, 0, 0, 0, 0, time.UTC),
					Week:             time.Date(2023, 8, 14, 0, 0, 0, 0, time.UTC),
					Month:            time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "empty userId",
			request: &calendarpb.EventsRequestByDate{
				StartDate: timestamppb.New(time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC)),
			},
			validateError: errors.New("field userId is empty"),
			code:          codes.InvalidArgument,
		},
		{
			name: "empty startDate",
			request: &calendarpb.EventsRequestByDate{
				UserId: 1,
			},
			validateError: errors.New("field startDate is empty"),
			code:          codes.InvalidArgument,
		},
		{
			name: "get events error",
			request: &calendarpb.EventsRequestByDate{
				UserId:    1,
				StartDate: timestamppb.New(time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC)),
			},
			date:      time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
			mockError: errors.New("unexpected error"),
			code:      codes.Internal,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.validateError == nil || tc.mockError != nil {
				appMock.On("GetEventByMonth", mock.Anything, tc.request.UserId, tc.date).
					Return(tc.events, tc.mockError).
					Once()
			}

			resp, err := client.GetEventsByMonth(context.Background(), tc.request)

			switch {
			case tc.mockError != nil:
				require.Equal(t, status.Error(tc.code, tc.mockError.Error()), err)

			case tc.validateError != nil:
				require.Equal(t, status.Error(tc.code, tc.validateError.Error()), err)

			default:
				require.NoError(t, err)
				require.Equal(t, toProtoEvents(tc.events).GetEvents(), resp.GetEvents())
			}
		})
	}
}
