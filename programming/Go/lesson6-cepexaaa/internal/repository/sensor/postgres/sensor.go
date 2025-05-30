package postgres

import (
	"context"
	"errors"
	"homework/internal/domain"
	"homework/internal/usecase"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SensorRepository struct {
	pool *pgxpool.Pool
}

func NewSensorRepository(pool *pgxpool.Pool) *SensorRepository {
	return &SensorRepository{
		pool: pool,
	}
}

var validSensorTypes = map[domain.SensorType]bool{
	domain.SensorTypeADC:            true,
	domain.SensorTypeContactClosure: true,
}

func (r *SensorRepository) SaveSensor(ctx context.Context, sensor *domain.Sensor) error {
	if sensor == nil {
		return errors.New("sensor is nil")
	}
	if !validSensorTypes[sensor.Type] {
		return usecase.ErrWrongSensorType
	}

	_, err := r.pool.Exec(ctx,
		`INSERT INTO sensors (serial_number, type, current_state, description, is_active, registered_at, last_activity)
		   VALUES ($1, $2, $3, $4, $5, $6, $7)
		   ON CONFLICT (serial_number) DO UPDATE
		   SET type = $2,
			   current_state = $3,
			   description = $4,
			   is_active = $5,
			   last_activity = $7`,
		sensor.SerialNumber,
		sensor.Type,
		sensor.CurrentState,
		sensor.Description,
		sensor.IsActive,
		time.Now(),
		sensor.LastActivity,
	)
	return err
}

func (r *SensorRepository) GetSensors(ctx context.Context) ([]domain.Sensor, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, serial_number, type, current_state, description, is_active, registered_at, last_activity 
		FROM sensors`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sensors []domain.Sensor
	for rows.Next() {
		var sensor domain.Sensor
		var sensorType string
		if err := rows.Scan(
			&sensor.ID,
			&sensor.SerialNumber,
			&sensorType,
			&sensor.CurrentState,
			&sensor.Description,
			&sensor.IsActive,
			&sensor.RegisteredAt,
			&sensor.LastActivity,
		); err != nil {
			return nil, err
		}
		sensor.Type = domain.SensorType(sensorType)
		sensors = append(sensors, sensor)
	}
	return sensors, nil
}

func (r *SensorRepository) GetSensorByID(ctx context.Context, id int64) (*domain.Sensor, error) {
	var sensor domain.Sensor
	var sensorType string
	err := r.pool.QueryRow(ctx,
		`SELECT id, serial_number, type, current_state, description, is_active, registered_at, last_activity 
		FROM sensors 
		WHERE id = $1`,
		id,
	).Scan(
		&sensor.ID,
		&sensor.SerialNumber,
		&sensorType,
		&sensor.CurrentState,
		&sensor.Description,
		&sensor.IsActive,
		&sensor.RegisteredAt,
		&sensor.LastActivity,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, usecase.ErrSensorNotFound
	}
	if err != nil {
		return nil, err
	}
	sensor.Type = domain.SensorType(sensorType)
	return &sensor, nil
}

func (r *SensorRepository) GetSensorBySerialNumber(ctx context.Context, sn string) (*domain.Sensor, error) {
	var sensor domain.Sensor
	var sensorType string
	err := r.pool.QueryRow(ctx,
		`SELECT id, serial_number, type, current_state, description, is_active, registered_at, last_activity 
		FROM sensors 
		WHERE serial_number = $1`,
		sn,
	).Scan(
		&sensor.ID,
		&sensor.SerialNumber,
		&sensorType,
		&sensor.CurrentState,
		&sensor.Description,
		&sensor.IsActive,
		&sensor.RegisteredAt,
		&sensor.LastActivity,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, usecase.ErrSensorNotFound
	}
	if err != nil {
		return nil, err
	}
	sensor.Type = domain.SensorType(sensorType)
	return &sensor, nil
}
