package usecase

import (
	"context"
	"errors"
	"homework/internal/domain"
)

type Sensor struct {
	sr SensorRepository
}

func NewSensor(sr SensorRepository) *Sensor {
	return &Sensor{sr: sr}
}

var validSensorTypes = map[domain.SensorType]bool{
	domain.SensorTypeADC:            true,
	domain.SensorTypeContactClosure: true,
}

func (s *Sensor) RegisterSensor(ctx context.Context, sensor *domain.Sensor) (*domain.Sensor, error) {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return nil, ctxErr
	}
	if sensor == nil {
		return nil, errors.New("sensor is nil")
	}
	if len(sensor.SerialNumber) != 10 {
		return nil, ErrWrongSensorSerialNumber
	}

	if !validSensorTypes[sensor.Type] {
		return sensor, ErrWrongSensorType
	}
	existingSensor, err := s.sr.GetSensorBySerialNumber(ctx, sensor.SerialNumber)

	if err != nil {
		if !errors.Is(err, ErrSensorNotFound) {
			return nil, err
		}
	} else {
		return existingSensor, nil
	}

	if err := s.sr.SaveSensor(ctx, sensor); err != nil {
		return nil, err
	}

	return sensor, nil
}

func (s *Sensor) GetSensors(ctx context.Context) ([]domain.Sensor, error) {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return nil, ctxErr
	}
	return s.sr.GetSensors(ctx)
}

func (s *Sensor) GetSensorByID(ctx context.Context, id int64) (*domain.Sensor, error) {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return nil, ctxErr
	}
	return s.sr.GetSensorByID(ctx, id)
}
