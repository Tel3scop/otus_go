package event

import (
	"context"

	"github.com/Tel3scop/helpers/logger"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/entity"
	"go.uber.org/zap"
)

func (s *serv) Update(ctx context.Context, dto entity.Event) error {
	logger.Info("Updating event...", zap.String("uuid", dto.ID), zap.String("title", dto.Title))

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		return s.eventStorage.Update(ctx, dto.ID, dto)
	})
	if err != nil {
		return err
	}

	logger.Info("Event updated", zap.String("uuid", dto.ID), zap.String("title", dto.Title))
	return nil
}
