package repository

import (
	"context"
	"fmt"
	"log"
	"message-server/internal/domain"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

type fileRepository struct {
	client *s3.Client
	bucket string
}

type FileRepositoryConfig struct {
	AccountID string
	AccessKey string
	SecretKey string
	Bucket    string
}

func NewFileRepository(c *FileRepositoryConfig) domain.FileRepository {
	log.Println(c)
	r2config, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("auto"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.AccessKey, c.SecretKey, "")),
	)
	if err != nil {
		panic(err)
	}

	client := s3.NewFromConfig(r2config, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", c.AccountID))
	})

	return &fileRepository{
		client: client,
		bucket: c.Bucket,
	}
}

func (r *fileRepository) GenerateListingUploadURL(listingID, key, contentType string) (string, error) {
	presignClient := s3.NewPresignClient(r.client)

	request, err := presignClient.PresignPutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      &r.bucket,
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
		Metadata:    map[string]string{"listing_id": listingID},
	}, s3.WithPresignExpires(3*time.Minute))

	if err != nil {
		return "", fmt.Errorf("failed to generate upload URL: %w", err)
	}

	return request.URL, nil
}

func (r *fileRepository) GenerateAvatarUploadURL(key, contentType string) (string, error) {
	presignClient := s3.NewPresignClient(r.client)

	request, err := presignClient.PresignPutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      &r.bucket,
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(3*time.Minute))

	if err != nil {
		return "", fmt.Errorf("failed to generate upload URL: %w", err)
	}

	return request.URL, nil
}

func (r *fileRepository) GenerateDownloadURL(key string) (string, error) {
	presignClient := s3.NewPresignClient(r.client)

	request, err := presignClient.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: &r.bucket,
		Key:    aws.String(key),
	}, s3.WithPresignExpires(7*24*time.Hour))

	if err != nil {
		return "", fmt.Errorf("failed to generate download URL: %w", err)
	}

	return request.URL, nil
}

func (r *fileRepository) DeleteFile(key string) error {
	_, err := r.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: &r.bucket,
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}
