package usecase

import (
	"context"
	"homework/internal/domain"
)

type User struct {
	ur  UserRepository
	sor SensorOwnerRepository
	sr  SensorRepository
}

func NewUser(ur UserRepository, sor SensorOwnerRepository, sr SensorRepository) *User {
	return &User{
		ur:  ur,
		sor: sor,
		sr:  sr,
	}
}

func (u *User) RegisterUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return nil, ctxErr
	}
	if user == nil || user.Name == "" {
		return nil, ErrInvalidUserName
	}
	err := u.ur.SaveUser(ctx, user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (u *User) AttachSensorToUser(ctx context.Context, userID, sensorID int64) error {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}
	_, err := u.ur.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	_, err = u.sr.GetSensorByID(ctx, sensorID)
	if err != nil {
		return err
	}

	sensorOwner := domain.SensorOwner{
		UserID:   userID,
		SensorID: sensorID,
	}

	return u.sor.SaveSensorOwner(ctx, sensorOwner)
}

func (u *User) GetUserSensors(ctx context.Context, userID int64) ([]domain.Sensor, error) {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return nil, ctxErr
	}
	user, err := u.ur.GetUserByID(ctx, userID)
	if user == nil || err != nil {
		return nil, ErrUserNotFound
	}
	sensorOwners, err := u.sor.GetSensorsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	sensors := make([]domain.Sensor, 0, len(sensorOwners))
	for _, so := range sensorOwners {
		sensor, err := u.sr.GetSensorByID(ctx, so.SensorID)
		if err != nil {
			return nil, err
		}
		sensors = append(sensors, *sensor)
	}

	return sensors, nil
}
