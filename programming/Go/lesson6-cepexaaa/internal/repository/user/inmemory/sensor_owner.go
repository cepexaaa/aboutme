package inmemory

import (
	"context"
	"homework/internal/domain"
	"sync"
)

type SensorOwnerRepository struct {
	mu           sync.RWMutex
	id           int64
	sensorOwners map[int64][]domain.SensorOwner
}

func NewSensorOwnerRepository() *SensorOwnerRepository {
	return &SensorOwnerRepository{id: 1, sensorOwners: make(map[int64][]domain.SensorOwner)}
}

func (r *SensorOwnerRepository) SaveSensorOwner(ctx context.Context, sensorOwner domain.SensorOwner) error {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	r.sensorOwners[sensorOwner.UserID] = append(r.sensorOwners[sensorOwner.UserID], sensorOwner)
	return nil
}

func (r *SensorOwnerRepository) GetSensorsByUserID(ctx context.Context, userID int64) ([]domain.SensorOwner, error) {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return nil, ctxErr
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	sensorOwners, exists := r.sensorOwners[userID]
	if !exists {
		return nil, nil
	}
	return sensorOwners, nil
}
