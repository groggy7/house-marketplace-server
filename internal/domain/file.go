package domain

import "io"

type FileUploadRequest struct {
	File        io.Reader
	FileName    string
	ContentType string
}

type FileUploadResponse struct {
	URL string `json:"url"`
}

type DeleteFileRequest struct {
	URL string `json:"url"`
}

type FileRepository interface {
	UploadFile(file io.Reader, fileName, contentType string) (*FileUploadResponse, error)
	DeleteFile(url string) error
}
