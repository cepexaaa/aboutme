package inmemory

import (
	"context"
	"errors"
	"homework/internal/domain"
	"homework/internal/usecase"
	"sync"
	"time"
)

type SensorRepository struct {
	mu          sync.RWMutex
	id          int64
	sensorDB    map[int64]*domain.Sensor
	serialIndex map[string]*domain.Sensor
}

func NewSensorRepository() *SensorRepository {
	return &SensorRepository{
		id:          1,
		sensorDB:    make(map[int64]*domain.Sensor),
		serialIndex: make(map[string]*domain.Sensor),
	}
}

var validSensorTypes = map[domain.SensorType]bool{
	domain.SensorTypeADC:            true,
	domain.SensorTypeContactClosure: true,
}

func (r *SensorRepository) SaveSensor(ctx context.Context, sensor *domain.Sensor) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}

	if sensor == nil {
		return errors.New("sensor is nil")
	}
	if !validSensorTypes[sensor.Type] {
		return usecase.ErrWrongSensorType
	}
	sensor.RegisteredAt = time.Now()
	r.sensorDB[r.id] = sensor
	r.serialIndex[sensor.SerialNumber] = sensor
	sensor.ID = r.id
	r.id++
	return nil
}

func (r *SensorRepository) GetSensors(ctx context.Context) ([]domain.Sensor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if ctxErr := ctx.Err(); ctxErr != nil {
		return nil, ctxErr
	}

	sensors := make([]domain.Sensor, 0, len(r.sensorDB))
	for _, sensor := range r.sensorDB {
		sensors = append(sensors, *sensor)
	}
	return sensors, nil
}

func (r *SensorRepository) GetSensorByID(ctx context.Context, id int64) (*domain.Sensor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if ctxErr := ctx.Err(); ctxErr != nil {
		return nil, ctxErr
	}

	sensor, exists := r.sensorDB[id]
	if !exists {
		return nil, usecase.ErrSensorNotFound
	}
	return sensor, nil
}

func (r *SensorRepository) GetSensorBySerialNumber(ctx context.Context, sn string) (*domain.Sensor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if ctxErr := ctx.Err(); ctxErr != nil {
		return nil, ctxErr
	}

	sensor, exists := r.serialIndex[sn]
	if !exists {
		return nil, usecase.ErrSensorNotFound
	}
	return sensor, nil
}
