package event

import (
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/entity"
	eventAPI "github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/pkg/event_v1"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// convertToProtoEvents преобразует список entity.Event в список protobuf Event.
func convertToProtoEvents(events []entity.Event) []*eventAPI.Event {
	protoEvents := make([]*eventAPI.Event, 0, len(events))
	for _, event := range events {
		protoEvents = append(protoEvents, &eventAPI.Event{
			Id:               event.ID,
			Title:            event.Title,
			DateTime:         timestamppb.New(event.DateTime),
			Duration:         durationpb.New(event.Duration),
			Description:      event.Description,
			UserId:           event.UserID,
			NotificationTime: durationpb.New(event.NotifyBefore),
		})
	}

	return protoEvents
}
