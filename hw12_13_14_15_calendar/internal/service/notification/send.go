package notification

import (
	"context"

	"github.com/Tel3scop/helpers/logger"
	"go.uber.org/zap"
)

// Send отправляет полученные сообщения из очереди.
func (s *notificationService) Send(ctx context.Context) error {
	notification, err := s.Dequeue(ctx)
	if err != nil {
		return err
	}

	logger.Info("Получено новое сообщение", zap.Any("notification", notification))

	return nil
}
