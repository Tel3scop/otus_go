//go:build integration

package integration

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/pkg/event_v1"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EventSuite struct {
	suite.Suite
	ctx         context.Context
	eventClient event_v1.EventServiceClient
}

func (s *EventSuite) SetupSuite() {
	const Localhost = "localhost:10051"
	host := os.Getenv("GRPC_ADDRESS")
	if host == "" {
		host = Localhost
	}

	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)

	s.ctx = context.Background()
	s.eventClient = event_v1.NewEventServiceClient(conn)
}

func (s *EventSuite) SetupTest() {
	seed := time.Now().UnixNano()
	rand.New(rand.NewSource(seed))

	s.T().Log("seed:", seed)
}

func TestEventSuite(t *testing.T) {
	suite.Run(t, new(EventSuite))
}

func (s *EventSuite) TestCreateEvent() {
	req := &event_v1.CreateEventRequest{
		Event: &event_v1.Event{
			Title:       "Test Event",
			Description: "Test Description",
			DateTime:    timestamppb.Now(),
			Duration:    &durationpb.Duration{Seconds: 3600},
			UserId:      uuid.New().String(),
		},
	}
	resp, err := s.eventClient.CreateEvent(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().NotEmpty(resp.EventId)
}

func (s *EventSuite) TestUpdateEvent() {
	// Создаем событие
	createReq := &event_v1.CreateEventRequest{
		Event: &event_v1.Event{
			Title:       "Test Event",
			Description: "Test Description",
			DateTime:    timestamppb.Now(),
			Duration:    &durationpb.Duration{Seconds: 3600},
			UserId:      uuid.New().String(),
		},
	}
	createResp, err := s.eventClient.CreateEvent(s.ctx, createReq)
	s.Require().NoError(err)
	s.Require().NotNil(createResp)
	s.Require().NotEmpty(createResp.EventId)

	// Обновляем событие
	updateReq := &event_v1.UpdateEventRequest{
		EventId: createResp.EventId,
		Event: &event_v1.Event{
			Title:       "Updated Event",
			Description: "Updated Description",
			DateTime:    timestamppb.Now(),
			Duration:    &durationpb.Duration{Seconds: 7200},
			UserId:      uuid.New().String(),
		},
	}
	updateResp, err := s.eventClient.UpdateEvent(s.ctx, updateReq)
	s.Require().NoError(err)
	s.Require().NotNil(updateResp)
	s.Require().True(updateResp.Success)
}

func (s *EventSuite) TestDeleteEvent() {
	// Создаем событие
	createReq := &event_v1.CreateEventRequest{
		Event: &event_v1.Event{
			Title:       "Test Event",
			Description: "Test Description",
			DateTime:    timestamppb.Now(),
			Duration:    &durationpb.Duration{Seconds: 3600},
			UserId:      uuid.New().String(),
		},
	}
	createResp, err := s.eventClient.CreateEvent(s.ctx, createReq)
	s.Require().NoError(err)
	s.Require().NotNil(createResp)
	s.Require().NotEmpty(createResp.EventId)

	// Удаляем событие
	deleteReq := &event_v1.DeleteEventRequest{
		EventId: createResp.EventId,
	}
	deleteResp, err := s.eventClient.DeleteEvent(s.ctx, deleteReq)
	s.Require().NoError(err)
	s.Require().NotNil(deleteResp)
	s.Require().True(deleteResp.Success)
}

func (s *EventSuite) TestListEventsForDay() {
	// Создаем событие
	createReq := &event_v1.CreateEventRequest{
		Event: &event_v1.Event{
			Title:       "Test Event",
			Description: "Test Description",
			DateTime:    timestamppb.Now(),
			Duration:    &durationpb.Duration{Seconds: 3600},
			UserId:      uuid.New().String(),
		},
	}
	_, err := s.eventClient.CreateEvent(s.ctx, createReq)
	s.Require().NoError(err)

	// Получаем список событий на день
	listReq := &event_v1.ListEventsForDayRequest{
		Date: timestamppb.Now(),
	}
	listResp, err := s.eventClient.ListEventsForDay(s.ctx, listReq)
	s.Require().NoError(err)
	s.Require().NotNil(listResp)
	s.Require().NotEmpty(listResp.Events)
}

func (s *EventSuite) TestListEventsForWeek() {
	// Создаем событие
	createReq := &event_v1.CreateEventRequest{
		Event: &event_v1.Event{
			Title:       "Test Event",
			Description: "Test Description",
			DateTime:    timestamppb.Now(),
			Duration:    &durationpb.Duration{Seconds: 3600},
			UserId:      uuid.New().String(),
		},
	}
	_, err := s.eventClient.CreateEvent(s.ctx, createReq)
	s.Require().NoError(err)

	// Получаем список событий на неделю
	listReq := &event_v1.ListEventsForWeekRequest{
		StartDate: timestamppb.Now(),
	}
	listResp, err := s.eventClient.ListEventsForWeek(s.ctx, listReq)
	s.Require().NoError(err)
	s.Require().NotNil(listResp)
	s.Require().NotEmpty(listResp.Events)
}

func (s *EventSuite) TestListEventsForMonth() {
	// Создаем событие
	createReq := &event_v1.CreateEventRequest{
		Event: &event_v1.Event{
			Title:       "Test Event",
			Description: "Test Description",
			DateTime:    timestamppb.Now(),
			Duration:    &durationpb.Duration{Seconds: 3600},
			UserId:      uuid.New().String(),
		},
	}
	_, err := s.eventClient.CreateEvent(s.ctx, createReq)
	s.Require().NoError(err)

	// Получаем список событий на месяц
	listReq := &event_v1.ListEventsForMonthRequest{
		StartDate: timestamppb.Now(),
	}
	listResp, err := s.eventClient.ListEventsForMonth(s.ctx, listReq)
	s.Require().NoError(err)
	s.Require().NotNil(listResp)
	s.Require().NotEmpty(listResp.Events)
}
