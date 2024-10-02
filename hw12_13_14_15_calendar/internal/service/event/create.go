package event

import (
	"context"

	"github.com/Tel3scop/helpers/logger"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/entity"
	"go.uber.org/zap"
)

func (s *serv) Create(ctx context.Context, dto entity.Event) (string, error) {
	logger.Info("Creating event...", zap.String("title", dto.Title), zap.String("desc", dto.Description))
	var id string

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.eventStorage.Create(ctx, dto)
		if errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	logger.Info("Event created", zap.String("uuid", id), zap.String("title", dto.Title))
	return id, nil
}
