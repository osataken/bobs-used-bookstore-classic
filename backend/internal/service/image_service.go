package service

import (
	"bytes"
	"io"

	"github.com/disintegration/imaging"
)

const (
	TargetWidth  = 400
	TargetHeight = 600
)

type ImageValidationService interface {
	IsSafe(reader io.Reader) (bool, error)
}

type LocalImageValidationService struct{}

func (s *LocalImageValidationService) IsSafe(reader io.Reader) (bool, error) {
	return true, nil
}

type RekognitionImageValidationService struct{}

func (s *RekognitionImageValidationService) IsSafe(reader io.Reader) (bool, error) {
	// In AWS mode, call Rekognition DetectModerationLabels
	// For now, return safe
	return true, nil
}

func ResizeImage(reader io.Reader) (io.Reader, error) {
	if reader == nil {
		return nil, nil
	}

	img, err := imaging.Decode(reader)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	if bounds.Dx() == TargetWidth && bounds.Dy() == TargetHeight {
		var buf bytes.Buffer
		if err := imaging.Encode(&buf, img, imaging.JPEG); err != nil {
			return nil, err
		}
		return &buf, nil
	}

	resized := imaging.Fit(img, TargetWidth, TargetHeight, imaging.Lanczos)

	var buf bytes.Buffer
	if err := imaging.Encode(&buf, resized, imaging.JPEG); err != nil {
		return nil, err
	}

	return &buf, nil
}
