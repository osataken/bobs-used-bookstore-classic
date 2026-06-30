package service

import (
	"bobs-used-bookstore-api/internal/config"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ImageService struct {
	cfg *config.Config
}

func NewImageService(cfg *config.Config) *ImageService {
	return &ImageService{cfg: cfg}
}

// ValidateImage checks if the image is safe (Rekognition or always-safe in local mode)
func (s *ImageService) ValidateImage(reader io.Reader) (bool, error) {
	if s.cfg.IsLocalImageValidation() {
		return true, nil
	}
	// AWS Rekognition would go here
	return true, nil
}

// ResizeImage resizes to fit within 400x600 preserving aspect ratio
// Uses stdlib only - basic nearest-neighbor resize
func (s *ImageService) ResizeImage(inputPath, outputPath string) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	maxW, maxH := 400, 600
	ratio := float64(srcW) / float64(srcH)
	targetRatio := float64(maxW) / float64(maxH)

	var newW, newH int
	if ratio > targetRatio {
		newW = maxW
		newH = int(float64(maxW) / ratio)
	} else {
		newH = maxH
		newW = int(float64(maxH) * ratio)
	}

	// Simple nearest-neighbor resize
	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	for y := 0; y < newH; y++ {
		for x := 0; x < newW; x++ {
			srcX := x * srcW / newW
			srcY := y * srcH / newH
			dst.Set(x, y, img.At(srcX+bounds.Min.X, srcY+bounds.Min.Y))
		}
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	ext := strings.ToLower(filepath.Ext(outputPath))
	if ext == ".jpg" || ext == ".jpeg" || format == "jpeg" {
		return jpeg.Encode(out, dst, &jpeg.Options{Quality: 85})
	}
	return png.Encode(out, dst)
}
