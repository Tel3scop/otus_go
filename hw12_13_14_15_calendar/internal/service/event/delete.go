package event

import (
	"context"
	"time"
)

func (s *serv) Delete(ctx context.Context, requestID string) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		return s.eventStorage.Delete(ctx, requestID)
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *serv) DeleteByDate(ctx context.Context, date time.Time) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		return s.eventStorage.DeleteByDate(ctx, date)
	})
	if err != nil {
		return err
	}

	return nil
}
