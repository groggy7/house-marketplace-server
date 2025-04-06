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
	if req.FullName == "" || req.Username == "" || req.Email == "" || req.Password == "" {
		return domain.ErrInvalidRequest
	}

	if len(req.Password) < 8 {
		return domain.ErrInvalidRequest
	}

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

func (s *AuthUseCase) Login(req *domain.LoginRequest) (*domain.User, string, error) {
	if req.Username != "" {
		user, err := s.authRepo.GetUserByUsername(req.Username)
		if err != nil {
			return nil, "", domain.ErrUserNotFound
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			return nil, "", domain.ErrInvalidCredentials
		}

		token, err := auth.GenerateToken(user.Username, user.Email, user.ID)
		if err != nil {
			return nil, "", err
		}
		return user, token, nil
	}

	if req.Email != "" {
		user, err := s.authRepo.GetUserByEmail(req.Email)
		if err != nil {
			return nil, "", domain.ErrUserNotFound
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			return nil, "", domain.ErrInvalidCredentials
		}

		token, err := auth.GenerateToken(user.Username, user.Email, user.ID)
		if err != nil {
			return nil, "", err
		}
		return user, token, nil
	}

	return nil, "", nil
}

func (s *AuthUseCase) CheckUserExists(userID string) (bool, error) {
	return s.authRepo.CheckUserExists(userID)
}

func (s *AuthUseCase) GetUserByUsername(username string) (*domain.User, error) {
	return s.authRepo.GetUserByUsername(username)
}

func (s *AuthUseCase) GetUserByEmail(email string) (*domain.User, error) {
	return s.authRepo.GetUserByEmail(email)
}

func hashPassword(password string) (string, error) {
	const cost = 12

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
