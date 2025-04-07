package repository

import (
	"context"
	"message-server/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) domain.UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) UpdateUser(name, avatarURL string, userID string) error {
	if avatarURL == "" {
		query := `
			UPDATE users SET full_name = $1 WHERE id = $2
		`
		_, err := r.pool.Exec(context.Background(), query, name, userID)
		if err != nil {
			return domain.ErrDatabaseError
		}
		return nil
	}
	query := `
		UPDATE users SET full_name = $1, avatar_url = $2 WHERE id = $3
	`
	_, err := r.pool.Exec(context.Background(), query, name, avatarURL, userID)
	if err != nil {
		return domain.ErrDatabaseError
	}
	return nil
}
