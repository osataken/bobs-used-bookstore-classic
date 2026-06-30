package service

import (
	"bobs-used-bookstore-api/internal/config"
	"io"
	"os"
	"path/filepath"
)

type FileService struct {
	cfg *config.Config
}

func NewFileService(cfg *config.Config) *FileService {
	return &FileService{cfg: cfg}
}

func (s *FileService) SaveFile(filename string, reader io.Reader) (string, error) {
	if s.cfg.IsLocalFileService() {
		return s.saveLocal(filename, reader)
	}
	// AWS S3 upload would go here
	return s.saveLocal(filename, reader)
}

func (s *FileService) DeleteFile(filename string) error {
	if s.cfg.IsLocalFileService() {
		path := filepath.Join(s.cfg.LocalImagePath, filename)
		return os.Remove(path)
	}
	// AWS S3 delete would go here
	return nil
}

func (s *FileService) saveLocal(filename string, reader io.Reader) (string, error) {
	dir := s.cfg.LocalImagePath
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	path := filepath.Join(dir, filename)
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return "", err
	}

	return "/images/" + filename, nil
}
