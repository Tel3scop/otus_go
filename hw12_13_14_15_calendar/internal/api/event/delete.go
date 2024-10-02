package event

import (
	"context"

	eventAPI "github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/pkg/event_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteEvent handler для удаления события по ID.
func (i *Implementation) DeleteEvent(ctx context.Context, req *eventAPI.DeleteEventRequest) (
	*eventAPI.DeleteEventResponse,
	error,
) {
	result := &eventAPI.DeleteEventResponse{}
	err := i.eventService.Delete(ctx, req.GetEventId())
	if err != nil {
		return result, status.Errorf(codes.Internal, "failed to delete event: %v", err)
	}
	result.Success = true

	return result, nil
}
