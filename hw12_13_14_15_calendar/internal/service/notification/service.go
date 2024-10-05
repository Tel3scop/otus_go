package notification

import (
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/service"
)

type notificationService struct {
	queueService service.QueueService
	eventService service.EventService
}

// NewService функция возвращает новый сервис для работы с уведомлениями.
func NewService(
	queueService service.QueueService,
	eventService service.EventService,
) service.NotificationService {
	return &notificationService{
		queueService: queueService,
		eventService: eventService,
	}
}
