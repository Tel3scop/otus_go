package event

import (
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/service"
	eventPkg "github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/pkg/event_v1"
)

// Implementation структура для работы с хэндлерами событий.
type Implementation struct {
	eventPkg.UnimplementedEventServiceServer
	eventService service.EventService
}

// NewImplementation новый экземпляр структуры Implementation.
func NewImplementation(eventService service.EventService) *Implementation {
	return &Implementation{
		eventService: eventService,
	}
}
