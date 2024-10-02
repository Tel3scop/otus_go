package event

import (
	"context"
	"time"

	"github.com/Tel3scop/helpers/logger"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/entity"
	"go.uber.org/zap"
)

func (s *serv) List(ctx context.Context, date time.Time, period entity.PeriodType) ([]entity.Event, error) {
	logger.Info(
		"Getting event list...",
		zap.String("date", date.String()),
		zap.String("period", string(period)),
	)

	var list []entity.Event
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		list, errTx = s.eventStorage.List(ctx, date, period)
		if errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	logger.Info(
		"Event list got",
		zap.Int("count", len(list)),
		zap.String("date", date.String()),
		zap.String("period", string(period)),
	)

	return list, nil
}
