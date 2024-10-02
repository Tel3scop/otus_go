package event

import (
	"context"

	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/entity"
	eventAPI "github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/pkg/event_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateEvent handler для создания нового события.
func (i *Implementation) CreateEvent(ctx context.Context, req *eventAPI.CreateEventRequest) (
	*eventAPI.CreateEventResponse,
	error,
) {
	event := req.GetEvent()

	entityEvent := entity.Event{
		Title:        event.GetTitle(),
		DateTime:     event.GetDateTime().AsTime(),
		Duration:     event.GetDuration().AsDuration(),
		Description:  event.GetDescription(),
		UserID:       event.GetUserId(),
		NotifyBefore: event.GetNotificationTime().AsDuration(),
	}

	eventID, err := i.eventService.Create(ctx, entityEvent)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create event: %v", err)
	}

	return &eventAPI.CreateEventResponse{EventId: eventID}, nil
}
