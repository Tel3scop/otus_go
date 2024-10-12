package notification

import (
	"context"
	"time"

	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/entity"
)

// NotifyOnEvent проверяет события и отправляет уведомления.
func (s *notificationService) NotifyOnEvent(ctx context.Context) error {
	events, err := s.eventService.List(ctx, time.Now(), entity.PeriodDay)
	if err != nil {
		return err
	}
	if len(events) == 0 {
		return nil
	}

	for _, event := range events {
		if time.Now().Add(event.NotifyBefore) != event.DateTime {
			continue
		}

		newNotification := entity.Notification{
			EventID:       event.ID,
			EventTitle:    event.Title,
			EventDateTime: event.DateTime,
			UserID:        event.UserID,
		}
		_ = s.Enqueue(ctx, newNotification)
	}
	return nil
}
