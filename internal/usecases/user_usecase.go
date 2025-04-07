package usecases

import (
	"message-server/internal/domain"
	"mime/multipart"
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

func (s *UserUseCase) UpdateUserAvatar(req *domain.UpdateUserRequest, avatarFile *multipart.FileHeader) error {
	oldUser, err := s.authRepo.GetUserByEmail(req.Email)
	if err != nil {
		return err
	}

	if oldUser.AvatarURL != "" {
		if err := s.fileRepo.DeleteFile(oldUser.AvatarURL); err != nil {
			return err
		}
	}

	src, err := avatarFile.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	uploadResponse, err := s.fileRepo.UploadProfilePicture(src, avatarFile.Filename, avatarFile.Header.Get("Content-Type"))
	if err != nil {
		return err
	}

	return s.userRepo.UpdateUser(oldUser.FullName, uploadResponse.URL, req.UserID)
}
