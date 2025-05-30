package postgres

import (
	"context"
	"homework/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SensorOwnerRepository struct {
	pool *pgxpool.Pool
}

func NewSensorOwnerRepository(pool *pgxpool.Pool) *SensorOwnerRepository {
	return &SensorOwnerRepository{
		pool,
	}
}

func (r *SensorOwnerRepository) SaveSensorOwner(ctx context.Context, sensorOwner domain.SensorOwner) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO sensors_users (user_id, sensor_id) 
		VALUES ($1, $2)`,
		sensorOwner.UserID,
		sensorOwner.SensorID,
	)
	return err
}

func (r *SensorOwnerRepository) GetSensorsByUserID(ctx context.Context, userID int64) ([]domain.SensorOwner, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT user_id, sensor_id 
		FROM sensors_users 
		WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sensorOwners []domain.SensorOwner
	for rows.Next() {
		var so domain.SensorOwner
		if err := rows.Scan(&so.UserID, &so.SensorID); err != nil {
			return nil, err
		}
		sensorOwners = append(sensorOwners, so)
	}
	return sensorOwners, nil
}
