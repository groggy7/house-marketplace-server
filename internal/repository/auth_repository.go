package repository

import (
	"context"
	"message-server/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type authRepository struct {
	pool *pgxpool.Pool
}

func NewAuthRepository(pool *pgxpool.Pool) domain.AuthRepository {
	return &authRepository{pool: pool}
}

func (r *authRepository) CreateUser(user *domain.User) error {
	query := `
		INSERT INTO users (full_name, username, email, password)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.pool.Exec(context.Background(), query, user.FullName, user.Username, user.Email, user.Password)
	return err
}

func (r *authRepository) GetUserByUsername(username string) (*domain.User, error) {
	query := `SELECT id, full_name, username, email, password FROM users WHERE username = $1`
	var user domain.User
	err := r.pool.QueryRow(context.Background(), query, username).Scan(&user.ID, &user.FullName, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) GetUserByEmail(email string) (*domain.User, error) {
	query := `SELECT id, full_name, username, email, password FROM users WHERE email = $1`
	var user domain.User
	err := r.pool.QueryRow(context.Background(), query, email).Scan(&user.ID, &user.FullName, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) CheckUserExists(userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	var exists bool
	err := r.pool.QueryRow(context.Background(), query, userID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
