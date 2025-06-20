package domain

import "errors"

type User struct {
	ID        string `json:"id"`
	FullName  string `json:"full_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	AvatarKey string `json:"avatar_key"`
}

type RegisterRequest struct {
	FullName  string `json:"full_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	AvatarKey string `json:"avatar_key"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthRepository interface {
	CreateUser(name, username, email, password string) error
	GetUserByUsername(username string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id string) (*User, error)
	UpdateUser(name, avatarURL string, userID string) error
	CheckUserExists(userID string) (bool, error)
	CheckUserCredentialsExist(username, email string) error
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidRequest     = errors.New("invalid request")
	ErrDuplicateUsername  = errors.New("username already exists")
	ErrDuplicateEmail     = errors.New("email already exists")
	ErrDatabaseError      = errors.New("database error")
)
