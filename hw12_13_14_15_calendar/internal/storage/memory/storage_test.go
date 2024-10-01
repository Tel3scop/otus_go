package memorystorage

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/entity"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryEventStorage_Create(t *testing.T) {
	repo := NewInMemoryEventStorage()
	event := entity.Event{
		Title:    "Test Event",
		DateTime: time.Now(),
	}

	id, err := repo.Create(context.Background(), event)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)

	savedEvent, exists := repo.events[id]
	assert.True(t, exists)
	assert.Equal(t, event.Title, savedEvent.Title)
}

func TestInMemoryEventStorage_Update(t *testing.T) {
	repo := NewInMemoryEventStorage()
	event := entity.Event{
		Title:    "Test Event",
		DateTime: time.Now(),
	}

	id, _ := repo.Create(context.Background(), event)

	event.Title = "Updated Event"
	err := repo.Update(context.Background(), id, event)
	assert.NoError(t, err)

	savedEvent, exists := repo.events[id]
	assert.True(t, exists)
	assert.Equal(t, event.Title, savedEvent.Title)

	err = repo.Update(context.Background(), "nonexistent-id", event)
	assert.Error(t, err)
	assert.Equal(t, storage.ErrEventNotFound, err)
}

func TestInMemoryEventStorage_Delete(t *testing.T) {
	repo := NewInMemoryEventStorage()
	event := entity.Event{
		Title:    "Test Event",
		DateTime: time.Now(),
	}

	id, _ := repo.Create(context.Background(), event)

	err := repo.Delete(context.Background(), id)
	assert.NoError(t, err)

	_, exists := repo.events[id]
	assert.False(t, exists)

	err = repo.Delete(context.Background(), "nonexistent-id")
	assert.Error(t, err)
	assert.Equal(t, storage.ErrEventNotFound, err)
}

func TestInMemoryEventStorage_List(t *testing.T) {
	repo := NewInMemoryEventStorage()
	now := time.Now()

	events := []entity.Event{
		{Title: "Event 1", DateTime: now},
		{Title: "Event 2", DateTime: now.AddDate(0, 0, 1)},
		{Title: "Event 3", DateTime: now.AddDate(0, 0, 7)},
	}

	for _, event := range events {
		repo.Create(context.Background(), event)
	}

	list, err := repo.List(context.Background(), now, "day")
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, events[0].Title, list[0].Title)

	list, err = repo.List(context.Background(), now, "week")
	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, events[0].Title, list[0].Title)
	assert.Equal(t, events[1].Title, list[1].Title)

	list, err = repo.List(context.Background(), now, "month")
	assert.NoError(t, err)
	assert.Len(t, list, 3)
	assert.Equal(t, events[0].Title, list[0].Title)
	assert.Equal(t, events[1].Title, list[1].Title)
	assert.Equal(t, events[2].Title, list[2].Title)

	_, err = repo.List(context.Background(), now, "invalid")
	assert.Error(t, err)
	assert.Equal(t, storage.ErrInvalidPeriod, err)
}

func TestInMemoryEventStorage_Concurrency(t *testing.T) {
	repo := NewInMemoryEventStorage()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			event := entity.Event{
				Title:    "Event " + uuid.New().String(),
				DateTime: time.Now(),
			}
			_, err := repo.Create(context.Background(), event)
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()

	// Проверяем, что все события были добавлены
	events, err := repo.List(context.Background(), time.Now(), "day")
	assert.NoError(t, err)
	assert.Len(t, events, 100)

	for _, event := range events {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			event := entity.Event{
				Title:    "Updated Event " + uuid.New().String(),
				DateTime: time.Now(),
			}
			err := repo.Update(context.Background(), id, event)
			assert.NoError(t, err)
		}(event.ID)
	}

	wg.Wait()

	for _, event := range events {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			err := repo.Delete(context.Background(), id)
			assert.NoError(t, err)
		}(event.ID)
	}

	wg.Wait()

	events, err = repo.List(context.Background(), time.Now(), "day")
	assert.NoError(t, err)
	assert.Len(t, events, 0)
}
