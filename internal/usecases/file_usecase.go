package usecases

import (
	"message-server/internal/domain"
	"mime/multipart"
)

type FileUseCase struct {
	fileRepo domain.FileRepository
}

func NewFileUseCase(fileRepo domain.FileRepository) *FileUseCase {
	return &FileUseCase{fileRepo: fileRepo}
}

func (s *FileUseCase) UploadListingPicture(file *multipart.FileHeader) (*domain.FileUploadResponse, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	req := &domain.FileUploadRequest{
		File:        src,
		FileName:    file.Filename,
		ContentType: file.Header.Get("Content-Type"),
	}

	return s.fileRepo.UploadListingPicture(req.File, req.FileName, req.ContentType)
}

func (s *FileUseCase) UploadProfilePicture(file *multipart.FileHeader) (*domain.FileUploadResponse, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	req := &domain.FileUploadRequest{
		File:        src,
		FileName:    file.Filename,
		ContentType: file.Header.Get("Content-Type"),
	}

	return s.fileRepo.UploadProfilePicture(req.File, req.FileName, req.ContentType)
}

func (s *FileUseCase) DeleteFile(url string) error {
	return s.fileRepo.DeleteFile(url)
}
