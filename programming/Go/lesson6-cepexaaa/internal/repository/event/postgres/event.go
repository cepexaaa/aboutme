package postgres

import (
	"context"
	"errors"
	"homework/internal/domain"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrEventNotFound = errors.New("event not found")

type EventRepository struct {
	pool *pgxpool.Pool
}

func NewEventRepository(pool *pgxpool.Pool) *EventRepository {
	return &EventRepository{
		pool,
	}
}

func (r *EventRepository) SaveEvent(ctx context.Context, event *domain.Event) error {
	if event == nil {
		return errors.New("event is nil")
	}

	_, err := r.pool.Exec(ctx,
		`INSERT INTO events (timestamp, sensor_serial_number, sensor_id, payload) 
		VALUES ($1, $2, $3, $4)`,
		event.Timestamp,
		event.SensorSerialNumber,
		event.SensorID,
		event.Payload,
	)

	return err
}

func (r *EventRepository) GetLastEventBySensorID(ctx context.Context, id int64) (*domain.Event, error) {
	var e domain.Event
	err := r.pool.QueryRow(ctx,
		`SELECT timestamp, sensor_serial_number, sensor_id, payload 
		FROM events 
		WHERE sensor_id = $1 
		ORDER BY timestamp DESC 
		LIMIT 1`,
		id,
	).Scan(
		&e.Timestamp,
		&e.SensorSerialNumber,
		&e.SensorID,
		&e.Payload,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEventNotFound
		}
		return nil, err
	}

	return &e, nil
}
