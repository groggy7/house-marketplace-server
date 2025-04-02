package usecases

import (
	"message-server/internal/controller/auth"
	"message-server/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	authRepo domain.AuthRepository
}

func NewAuthUseCase(authRepo domain.AuthRepository) *AuthUseCase {
	return &AuthUseCase{authRepo: authRepo}
}

func (s *AuthUseCase) Register(req *domain.RegisterRequest) error {
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &domain.User{
		FullName: req.FullName,
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	return s.authRepo.CreateUser(user)
}

func (s *AuthUseCase) Login(req *domain.LoginRequest) (string, error) {
	if req.Username != "" {
		user, err := s.authRepo.GetUserByUsername(req.Username)
		if err != nil {
			return "", domain.ErrUserNotFound
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			return "", domain.ErrInvalidCredentials
		}

		return auth.GenerateToken(user.Username, user.Email)
	}

	if req.Email != "" {
		user, err := s.authRepo.GetUserByEmail(req.Email)
		if err != nil {
			return "", domain.ErrUserNotFound
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			return "", domain.ErrInvalidCredentials
		}

		return auth.GenerateToken(user.Username, user.Email)
	}

	return "", nil
}

func hashPassword(password string) (string, error) {
	const cost = 12

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
