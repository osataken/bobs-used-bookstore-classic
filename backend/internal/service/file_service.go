package service

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"bobs-used-bookstore-api/internal/config"

	"github.com/google/uuid"
)

type FileService interface {
	Save(reader io.Reader, filename string) (string, error)
	Delete(url string) error
}

type LocalFileService struct {
	basePath string
}

func NewLocalFileService(cfg *config.Config) *LocalFileService {
	os.MkdirAll(cfg.LocalImagePath, 0755)
	return &LocalFileService{basePath: cfg.LocalImagePath}
}

func (s *LocalFileService) Save(reader io.Reader, filename string) (string, error) {
	if reader == nil {
		return "", nil
	}

	ext := filepath.Ext(filename)
	newFilename := uuid.New().String() + ext
	fullPath := filepath.Join(s.basePath, newFilename)

	file, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return "", err
	}

	return fmt.Sprintf("/uploads/coverimages/%s", newFilename), nil
}

func (s *LocalFileService) Delete(url string) error {
	if url == "" {
		return nil
	}
	filename := filepath.Base(url)
	fullPath := filepath.Join(s.basePath, filename)
	return os.Remove(fullPath)
}

type S3FileService struct {
	bucket          string
	cloudFrontDomain string
}

func NewS3FileService(cfg *config.Config) *S3FileService {
	return &S3FileService{
		bucket:          cfg.S3Bucket,
		cloudFrontDomain: cfg.CloudFrontDomain,
	}
}

func (s *S3FileService) Save(reader io.Reader, filename string) (string, error) {
	if reader == nil {
		return "", nil
	}
	// In AWS mode, upload to S3
	ext := filepath.Ext(filename)
	newFilename := uuid.New().String() + ext
	// Placeholder: actual S3 upload would go here
	return fmt.Sprintf("https://%s/%s", s.cloudFrontDomain, newFilename), nil
}

func (s *S3FileService) Delete(url string) error {
	// Placeholder: actual S3 delete would go here
	return nil
}
