package service

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"bobs-used-bookstore-api/internal/config"
)

// FileService defines the interface for file storage
type FileService interface {
	Save(contents io.Reader, filename string) (string, error)
	Delete(filePath string) error
}

// LocalFileService saves files to the local filesystem
type LocalFileService struct {
	basePath string
}

func NewLocalFileService(basePath string) *LocalFileService {
	// Ensure directory exists
	_ = os.MkdirAll(basePath, 0755)
	return &LocalFileService{basePath: basePath}
}

func (s *LocalFileService) Save(contents io.Reader, filename string) (string, error) {
	ext := filepath.Ext(filename)
	newFilename := uuid.New().String() + ext

	fullPath := filepath.Join(s.basePath, newFilename)
	file, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := io.Copy(file, contents); err != nil {
		return "", err
	}

	return fmt.Sprintf("/uploads/coverimages/%s", newFilename), nil
}

func (s *LocalFileService) Delete(filePath string) error {
	// Extract filename from URL path
	filename := filepath.Base(filePath)
	fullPath := filepath.Join(s.basePath, filename)
	return os.Remove(fullPath)
}

// S3FileService saves files to AWS S3
type S3FileService struct {
	bucketName     string
	cloudFrontDomain string
}

func NewS3FileService(cfg *config.Config) *S3FileService {
	return &S3FileService{
		bucketName:       cfg.S3BucketName,
		cloudFrontDomain: cfg.CloudFrontDomain,
	}
}

func (s *S3FileService) Save(contents io.Reader, filename string) (string, error) {
	// TODO: Implement S3 upload when AWS mode is needed
	ext := filepath.Ext(filename)
	key := "coverimages/" + uuid.New().String() + ext
	return fmt.Sprintf("https://%s/%s", s.cloudFrontDomain, key), nil
}

func (s *S3FileService) Delete(filePath string) error {
	// TODO: Implement S3 delete when AWS mode is needed
	return nil
}
