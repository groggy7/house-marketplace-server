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

func (r *authRepository) CreateUser(name, username, email, password string) error {
	query := `
		INSERT INTO users (full_name, username, email, password)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.pool.Exec(context.Background(), query, name, username, email, password)
	if err != nil {
		return domain.ErrDatabaseError
	}
	return nil
}

func (r *authRepository) GetUserByUsername(username string) (*domain.User, error) {
	query := `SELECT id, full_name, username, email, password, avatar_key FROM users WHERE username = $1`
	var user domain.User
	err := r.pool.QueryRow(context.Background(), query, username).Scan(&user.ID,
		&user.FullName, &user.Username, &user.Email, &user.Password, &user.AvatarKey)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) GetUserByEmail(email string) (*domain.User, error) {
	query := `SELECT id, full_name, username, email, password, avatar_key FROM users WHERE email = $1`
	var user domain.User
	err := r.pool.QueryRow(context.Background(), query, email).Scan(&user.ID,
		&user.FullName, &user.Username, &user.Email, &user.Password, &user.AvatarKey)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) GetUserByID(id string) (*domain.User, error) {
	query := `SELECT id, full_name, username, email, password, avatar_key FROM users WHERE id = $1`
	var user domain.User
	err := r.pool.QueryRow(context.Background(), query, id).Scan(&user.ID,
		&user.FullName, &user.Username, &user.Email, &user.Password, &user.AvatarKey)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) UpdateUser(name, avatarKey string, userID string) error {
	if avatarKey == "" {
		query := `UPDATE users SET full_name = $1 WHERE id = $2`
		_, err := r.pool.Exec(context.Background(), query, name, userID)
		if err != nil {
			return domain.ErrDatabaseError
		}
	} else {
		query := `UPDATE users SET full_name = $1, avatar_key = $2 WHERE id = $3`
		_, err := r.pool.Exec(context.Background(), query, name, avatarKey, userID)
		if err != nil {
			return domain.ErrDatabaseError
		}
	}
	return nil
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

func (r *authRepository) CheckUserCredentialsExist(username, email string) error {
	query := `
		SELECT 
			(SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)) as username_exists,
			(SELECT EXISTS(SELECT 1 FROM users WHERE email = $2)) as email_exists
	`
	var usernameExists, emailExists bool
	err := r.pool.QueryRow(context.Background(), query, username, email).Scan(&usernameExists, &emailExists)
	if err != nil {
		return domain.ErrDatabaseError
	}

	if usernameExists {
		return domain.ErrDuplicateUsername
	}

	if emailExists {
		return domain.ErrDuplicateEmail
	}

	return nil
}
