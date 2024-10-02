package event

import (
	"context"
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
