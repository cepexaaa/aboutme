package inmemory

import (
	"context"
	"errors"
	"homework/internal/domain"
	"homework/internal/usecase"
	"sync"
)

type EventRepository struct {
	mu     sync.RWMutex
	events map[int64][]*domain.Event
}

func NewEventRepository() *EventRepository {
	return &EventRepository{events: map[int64][]*domain.Event{}}
}

func (r *EventRepository) SaveEvent(ctx context.Context, event *domain.Event) error {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	if event == nil {
		return errors.New("event is nil")
	}
	r.events[event.SensorID] = append(r.events[event.SensorID], event)
	return nil
}

func (r *EventRepository) GetLastEventBySensorID(ctx context.Context, id int64) (*domain.Event, error) {
	events, err := r.GetEventsBySensorID(ctx, id)
	if err != nil {
		return nil, err
	}
	maxTime := events[0].Timestamp
	lastEvent := events[0]
	for _, e := range events {
		if !maxTime.After(e.Timestamp) {
			maxTime = e.Timestamp
			lastEvent = e
		}
	}
	return lastEvent, nil
}

func (r *EventRepository) GetEventsBySensorID(ctx context.Context, sensorID int64) ([]*domain.Event, error) {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return nil, ctxErr
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	events, exists := r.events[sensorID]
	if !exists || len(events) == 0 {
		return nil, usecase.ErrEventNotFound
	}
	return events, nil
}
