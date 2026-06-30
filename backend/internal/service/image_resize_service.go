package service

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"

	// Register image decoders
	_ "image/gif"
	_ "image/png"

	"golang.org/x/image/draw"
)

const (
	targetWidth  = 400
	targetHeight = 600
)

type ImageResizeService struct{}

func NewImageResizeService() *ImageResizeService {
	return &ImageResizeService{}
}

// ResizeImage resizes the image to 400x600 preserving aspect ratio
func (s *ImageResizeService) ResizeImage(reader io.Reader) (io.Reader, error) {
	src, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	// Calculate dimensions preserving aspect ratio
	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	newW, newH := targetWidth, targetHeight
	srcRatio := float64(srcW) / float64(srcH)
	targetRatio := float64(targetWidth) / float64(targetHeight)

	if srcRatio > targetRatio {
		newH = int(float64(targetWidth) / srcRatio)
	} else {
		newW = int(float64(targetHeight) * srcRatio)
	}

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.BiLinear.Scale(dst, dst.Bounds(), src, srcBounds, draw.Over, nil)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 85}); err != nil {
		return nil, err
	}

	return &buf, nil
}
