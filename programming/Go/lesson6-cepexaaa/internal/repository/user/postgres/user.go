package postgres

import (
	"context"
	"errors"
	"homework/internal/domain"
	"homework/internal/usecase"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (r *UserRepository) SaveUser(ctx context.Context, user *domain.User) error {
	if user == nil {
		return errors.New("user is nil")
	}

	_, err := r.pool.Exec(ctx,
		`INSERT INTO users (name) 
		VALUES ($1)`,
		user.Name,
	)

	return err
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, name 
		FROM users 
		WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Name)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, usecase.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
