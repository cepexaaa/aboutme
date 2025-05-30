package usecase

import (
	"context"
	"errors"
	"homework/internal/domain"
	"time"
)

type Event struct {
	er EventRepository
	sr SensorRepository
}

func NewEvent(er EventRepository, sr SensorRepository) *Event {
	return &Event{er: er, sr: sr}
}

func (e *Event) ReceiveEvent(ctx context.Context, event *domain.Event) error {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}
	if event == nil {
		return errors.New("event is nil")
	}
	if event.Timestamp.IsZero() {
		return ErrInvalidEventTimestamp
	}

	sensor, err := e.sr.GetSensorBySerialNumber(ctx, event.SensorSerialNumber)
	if err != nil {
		return err
	}

	event.SensorID = sensor.ID

	sensor.CurrentState = event.Payload

	sensor.LastActivity = time.Now()

	err = e.er.SaveEvent(ctx, event)
	if err != nil {
		return err
	}

	err = e.sr.SaveSensor(ctx, sensor)
	if err != nil {
		return err
	}

	return nil
}

func (e *Event) GetLastEventBySensorID(ctx context.Context, id int64) (*domain.Event, error) {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return nil, ctxErr
	}
	return e.er.GetLastEventBySensorID(ctx, id)
}
