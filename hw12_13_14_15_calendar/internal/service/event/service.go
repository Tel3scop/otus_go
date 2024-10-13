package event

import (
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/client/db"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/service"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/storage"
)

type serv struct {
	eventStorage storage.EventStorage
	txManager    db.TxManager
}

// NewService функция возвращает новый сервис событий.
func NewService(
	eventStorage storage.EventStorage,
	txManager db.TxManager,
) service.EventService {
	return &serv{
		eventStorage: eventStorage,
		txManager:    txManager,
	}
}
