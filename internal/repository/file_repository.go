package repository

import (
	"context"
	"fmt"
	"io"
	"message-server/internal/domain"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

type fileRepository struct {
	storageClient *storage.Client
	bucketName    string
}

func NewFileRepository(credentialsPath, bucketName string) domain.FileRepository {
	ctx := context.Background()

	opt := option.WithCredentialsFile(credentialsPath)

	storageClient, err := storage.NewClient(ctx, opt)
	if err != nil {
		panic(err)
	}

	return &fileRepository{
		storageClient: storageClient,
		bucketName:    bucketName,
	}
}

func (r *fileRepository) UploadFile(file io.Reader, fileName, contentType string) (*domain.FileUploadResponse, error) {
	ctx := context.Background()

	ext := filepath.Ext(fileName)
	filename := fmt.Sprintf("listings/%s%s", uuid.New().String(), ext)

	bucket := r.storageClient.Bucket(r.bucketName)

	obj := bucket.Object(filename)

	w := obj.NewWriter(ctx)
	w.ContentType = contentType
	w.CacheControl = "public, max-age=31536000"

	if _, err := io.Copy(w, file); err != nil {
		return nil, fmt.Errorf("error copying file: %v", err)
	}

	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("error closing writer: %v", err)
	}

	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return nil, fmt.Errorf("error making file public: %v", err)
	}

	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", r.bucketName, filename)

	return &domain.FileUploadResponse{
		URL: url,
	}, nil
}

func (r *fileRepository) DeleteFile(url string) error {
	ctx := context.Background()

	bucket := r.storageClient.Bucket(r.bucketName)

	objectName := strings.TrimPrefix(url, fmt.Sprintf("https://storage.googleapis.com/%s/", r.bucketName))

	obj := bucket.Object(objectName)
	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("error deleting file: %v", err)
	}

	return nil
}
