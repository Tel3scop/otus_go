package notification

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Tel3scop/helpers/logger"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/entity"
	"go.uber.org/zap"
)

// Dequeue читает уведомление из очереди.
func (s *notificationService) Dequeue(ctx context.Context) (*entity.Notification, error) {
	message, err := s.queueService.Dequeue(ctx)
	if err != nil {
		return nil, err
	}

	var notification entity.Notification
	err = json.Unmarshal(message, &notification)
	if err != nil {
		logger.Error("failed to unmarshal notification", zap.Error(err), zap.Any("notification", notification))
		return nil, fmt.Errorf("failed to unmarshal notification: %w", err)
	}

	return &notification, nil
}
