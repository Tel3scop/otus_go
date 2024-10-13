package event

import (
	"context"

	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/entity"
	eventAPI "github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/pkg/event_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateEvent handler для обновления существующего события.
func (i *Implementation) UpdateEvent(ctx context.Context, req *eventAPI.UpdateEventRequest) (
	*eventAPI.UpdateEventResponse,
	error,
) {
	result := &eventAPI.UpdateEventResponse{}
	event := req.GetEvent()

	entityEvent := entity.Event{
		ID:           req.GetEventId(),
		Title:        event.GetTitle(),
		DateTime:     event.GetDateTime().AsTime(),
		Duration:     event.GetDuration().AsDuration(),
		Description:  event.GetDescription(),
		UserID:       event.GetUserId(),
		NotifyBefore: event.GetNotificationTime().AsDuration(),
	}

	err := i.eventService.Update(ctx, entityEvent)
	if err != nil {
		return result, status.Errorf(codes.Internal, "failed to update event: %v", err)
	}
	result.Success = true

	return result, nil
}
