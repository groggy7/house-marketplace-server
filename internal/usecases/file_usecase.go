package usecases

import (
	"message-server/internal/domain"

	"github.com/google/uuid"
)

type FileUseCase struct {
	fileRepo domain.FileRepository
}

func NewFileUseCase(fileRepo domain.FileRepository) *FileUseCase {
	return &FileUseCase{fileRepo: fileRepo}
}

func (s *FileUseCase) GenerateListingUploadURL(req *domain.GenerateListingUploadURLRequest) (*domain.URLResponse, error) {
	key := uuid.NewString()
	URL, err := s.fileRepo.GenerateListingUploadURL(req.ListingID, key, req.ContentType)
	if err != nil {
		return nil, err
	}

	return &domain.URLResponse{URL: URL}, nil
}

func (s *FileUseCase) GenerateAvatarUploadURL(req *domain.GenerateAvatarUploadURLRequest) (*domain.URLResponse, error) {
	URL, err := s.fileRepo.GenerateAvatarUploadURL(req.ID, req.ContentType)
	if err != nil {
		return nil, err
	}

	return &domain.URLResponse{URL: URL}, nil
}

func (s *FileUseCase) GenerateDownloadURL(key string) (*domain.URLResponse, error) {
	URL, err := s.fileRepo.GenerateDownloadURL(key)
	if err != nil {
		return nil, err
	}

	return &domain.URLResponse{URL: URL}, nil
}

func (s *FileUseCase) DeleteFile(key string) error {
	return s.fileRepo.DeleteFile(key)
}
