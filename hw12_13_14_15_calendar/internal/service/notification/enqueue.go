package notification

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Tel3scop/helpers/logger"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/entity"
	"go.uber.org/zap"
)

// Enqueue помещает уведомление в очередь.
func (s *notificationService) Enqueue(ctx context.Context, notification entity.Notification) error {
	body, err := json.Marshal(notification)
	if err != nil {
		logger.Error("failed to encode notification", zap.Error(err), zap.Any("notification", notification))
		return fmt.Errorf("failed to marshal notification: %w", err)
	}
	return s.queueService.Enqueue(ctx, body)
}
