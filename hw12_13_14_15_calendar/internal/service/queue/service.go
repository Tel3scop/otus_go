package queue

import (
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/service"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/storage"
)

type queueServ struct {
	queueStorage storage.QueueStorage
	eventService service.EventService
}

// NewService функция возвращает новый сервис для работы с очередями.
func NewService(
	queueStorage storage.QueueStorage,
	eventService service.EventService,
) service.QueueService {
	return &queueServ{
		queueStorage: queueStorage,
		eventService: eventService,
	}
}
