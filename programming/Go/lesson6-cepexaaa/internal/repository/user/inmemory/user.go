package inmemory

import (
	"context"
	"errors"
	"homework/internal/domain"
	"homework/internal/usecase"
	"sync"
)

type UserRepository struct {
	id     int64
	mu     sync.RWMutex
	userDB map[int64]*domain.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		id:     1,
		userDB: make(map[int64]*domain.User),
	}
}

func (r *UserRepository) SaveUser(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}

	if user == nil {
		return errors.New("user is nil")
	}
	r.userDB[r.id] = user
	r.id++
	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return nil, ctxErr
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.userDB[id]
	if !exists {
		return nil, usecase.ErrUserNotFound
	}
	return user, nil
}
