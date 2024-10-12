package event

import (
	"context"

	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/entity"
	eventAPI "github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/pkg/event_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListEventsForDay handler для получения списка событий на определенный день.
func (i *Implementation) ListEventsForDay(ctx context.Context, req *eventAPI.ListEventsForDayRequest) (
	*eventAPI.ListEventsForDayResponse,
	error,
) {
	events, err := i.eventService.List(ctx, req.GetDate().AsTime(), entity.PeriodDay)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list events for day: %v", err)
	}

	return &eventAPI.ListEventsForDayResponse{Events: convertToProtoEvents(events)}, nil
}

// ListEventsForWeek handler для получения списка событий на определенную неделю.
func (i *Implementation) ListEventsForWeek(ctx context.Context, req *eventAPI.ListEventsForWeekRequest) (
	*eventAPI.ListEventsForWeekResponse,
	error,
) {
	events, err := i.eventService.List(ctx, req.GetStartDate().AsTime(), entity.PeriodWeek)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list events for week: %v", err)
	}

	return &eventAPI.ListEventsForWeekResponse{Events: convertToProtoEvents(events)}, nil
}

// ListEventsForMonth handler для получения списка событий на определенный месяц.
func (i *Implementation) ListEventsForMonth(ctx context.Context, req *eventAPI.ListEventsForMonthRequest) (
	*eventAPI.ListEventsForMonthResponse,
	error,
) {
	events, err := i.eventService.List(ctx, req.GetStartDate().AsTime(), entity.PeriodMonth)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list events for month: %v", err)
	}

	return &eventAPI.ListEventsForMonthResponse{Events: convertToProtoEvents(events)}, nil
}
