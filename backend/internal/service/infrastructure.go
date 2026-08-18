package service

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"bobs-used-bookstore-api/internal/config"

	"github.com/google/uuid"
	"golang.org/x/image/draw"
)

var ErrImageUnsafe = errors.New("image failed safety validation")

// FileServiceInterface for saving/deleting files
type FileServiceInterface interface {
	Save(data []byte, filename string) (string, error)
	Delete(url string) error
}

// ImageResizeServiceInterface for resizing images
type ImageResizeServiceInterface interface {
	ResizeImage(data []byte, maxWidth, maxHeight int) ([]byte, error)
}

// ImageValidationServiceInterface for validating image safety
type ImageValidationServiceInterface interface {
	IsSafe(data []byte) (bool, error)
}

// LocalFileService saves files to local filesystem
type LocalFileService struct {
	storagePath string
}

func (s *LocalFileService) Save(data []byte, filename string) (string, error) {
	dir := filepath.Join(s.storagePath, "coverimages")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	ext := filepath.Ext(filename)
	newFilename := uuid.New().String() + ext
	fullPath := filepath.Join(dir, newFilename)

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return "", err
	}

	return "/images/coverimages/" + newFilename, nil
}

func (s *LocalFileService) Delete(url string) error {
	if url == "" {
		return nil
	}
	// Convert URL path to filesystem path
	relPath := strings.TrimPrefix(url, "/images/")
	fullPath := filepath.Join(s.storagePath, relPath)
	return os.Remove(fullPath)
}

// S3FileService saves files to S3 (stub for aws mode)
type S3FileService struct {
	bucket       string
	cloudFrontURL string
}

func (s *S3FileService) Save(data []byte, filename string) (string, error) {
	// In production, this would upload to S3
	ext := filepath.Ext(filename)
	key := "coverimages/" + uuid.New().String() + ext
	return fmt.Sprintf("%s/%s", s.cloudFrontURL, key), nil
}

func (s *S3FileService) Delete(url string) error {
	// In production, this would delete from S3
	return nil
}

// NewFileService creates the appropriate file service based on config
func NewFileService(cfg *config.Config) FileServiceInterface {
	if cfg.FileServiceMode == "aws" {
		return &S3FileService{
			bucket:       cfg.S3Bucket,
			cloudFrontURL: cfg.CloudFrontURL,
		}
	}
	return &LocalFileService{storagePath: cfg.LocalStoragePath}
}

// ImageResizeService resizes images preserving aspect ratio
type ImageResizeService struct{}

func NewImageResizeService() *ImageResizeService {
	return &ImageResizeService{}
}

func (s *ImageResizeService) ResizeImage(data []byte, maxWidth, maxHeight int) ([]byte, error) {
	reader := bytes.NewReader(data)
	img, format, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	origWidth := bounds.Dx()
	origHeight := bounds.Dy()

	// Calculate new dimensions preserving aspect ratio
	newWidth, newHeight := maxWidth, maxHeight
	ratio := float64(origWidth) / float64(origHeight)
	targetRatio := float64(maxWidth) / float64(maxHeight)

	if ratio > targetRatio {
		newHeight = int(float64(maxWidth) / ratio)
	} else {
		newWidth = int(float64(maxHeight) * ratio)
	}

	// Resize
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.BiLinear.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)

	var buf bytes.Buffer
	switch format {
	case "png":
		err = png.Encode(&buf, dst)
	default:
		err = jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 85})
	}
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// LocalImageValidationService always returns safe (local mode)
type LocalImageValidationService struct{}

func (s *LocalImageValidationService) IsSafe(_ []byte) (bool, error) {
	return true, nil
}

// RekognitionImageValidationService uses AWS Rekognition (stub)
type RekognitionImageValidationService struct{}

func (s *RekognitionImageValidationService) IsSafe(_ []byte) (bool, error) {
	// In production, this would call AWS Rekognition
	return true, nil
}

// NewImageValidationService creates the appropriate validation service
func NewImageValidationService(cfg *config.Config) ImageValidationServiceInterface {
	if cfg.ImageValidationMode == "aws" {
		return &RekognitionImageValidationService{}
	}
	return &LocalImageValidationService{}
}

// Register png and jpeg decoders
func init() {
	image.RegisterFormat("png", "\x89PNG", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpeg", "\xff\xd8", jpeg.Decode, jpeg.DecodeConfig)
}

// Ensure io import is used (for potential future streaming)
var _ = io.Discard
