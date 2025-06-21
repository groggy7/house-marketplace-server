package usecases

import (
	"message-server/internal/domain"
)

type UserUseCase struct {
	userRepo domain.UserRepository
	authRepo domain.AuthRepository
}

func NewUserUseCase(userRepo domain.UserRepository, authRepo domain.AuthRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
		authRepo: authRepo,
	}
}

func (s *UserUseCase) UpdateUserInfo(req *domain.UpdateUserRequest) error {
	return s.userRepo.UpdateUser(req.FullName, req.AvatarKey, req.UserID)
}
