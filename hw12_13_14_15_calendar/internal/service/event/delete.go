package event

import (
	"context"
)

func (s *serv) Delete(ctx context.Context, requestID string) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		if errTx := s.eventStorage.Delete(ctx, requestID); errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
