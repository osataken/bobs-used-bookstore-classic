package service

import "io"

// ImageValidationService defines the interface for image safety validation
type ImageValidationService interface {
	IsSafe(image io.Reader) (bool, error)
}

// LocalImageValidationService always returns safe (local mode)
type LocalImageValidationService struct{}

func NewLocalImageValidationService() *LocalImageValidationService {
	return &LocalImageValidationService{}
}

func (s *LocalImageValidationService) IsSafe(image io.Reader) (bool, error) {
	return true, nil
}

// RekognitionImageValidationService uses AWS Rekognition for image moderation
type RekognitionImageValidationService struct{}

func NewRekognitionImageValidationService() *RekognitionImageValidationService {
	return &RekognitionImageValidationService{}
}

func (s *RekognitionImageValidationService) IsSafe(image io.Reader) (bool, error) {
	// TODO: Implement Rekognition moderation when AWS mode is needed
	return true, nil
}
