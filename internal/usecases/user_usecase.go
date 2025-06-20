package usecases

import (
	"message-server/internal/domain"
)

type UserUseCase struct {
	userRepo domain.UserRepository
	fileRepo domain.FileRepository
	authRepo domain.AuthRepository
}

func NewUserUseCase(userRepo domain.UserRepository, fileRepo domain.FileRepository, authRepo domain.AuthRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
		fileRepo: fileRepo,
		authRepo: authRepo,
	}
}

func (s *UserUseCase) UpdateUserInfo(req *domain.UpdateUserRequest) error {
	if req.FullName == "" {
		return domain.ErrInvalidRequest
	}

	return s.userRepo.UpdateUser(req.FullName, "", req.UserID)
}

func (s *UserUseCase) UpdateUserAvatar(req *domain.UpdateUserRequest) error {
	oldUser, err := s.authRepo.GetUserByEmail(req.Email)
	if err != nil {
		return err
	}

	if oldUser.AvatarKey != "" {
		if err := s.fileRepo.DeleteFile(oldUser.AvatarKey); err != nil {
			return err
		}
	}

	return s.userRepo.UpdateUser(oldUser.FullName, req.AvatarKey, req.UserID)
}
