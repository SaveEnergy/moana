// Package avatar decodes and normalizes uploaded images into small JPEGs for user profile avatars.
package avatar

import (
	"bytes"
	"errors"
	"image"
	"image/draw" //nolint:goimports
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	xdraw "golang.org/x/image/draw"
	_ "golang.org/x/image/webp" // register webp
)

const (
	// MaxUploadBytes is the largest raw file size accepted (before square crop and resize).
	MaxUploadBytes = 1 << 20 // 1 MiB
	jpegQuality    = 86
	outSide        = 256
	maxEdge        = 8000
	maxPixels      = 20_000_000
)

var (
	// ErrEmpty is returned for a zero-length upload.
	ErrEmpty = errors.New("empty file")
	// ErrTooLarge is returned when the raw bytes exceed [MaxUploadBytes].
	ErrTooLarge = errors.New("file too large")
	// ErrInvalid is returned for corrupt or unknown image data.
	ErrInvalid = errors.New("not a displayable image")
	// ErrTooBig is returned when width×height is excessive (decode / DoS guard).
	ErrTooBig = errors.New("image dimensions too large")
)

// ToJPEG decodes a PNG/JPEG/GIF/WEBP (and similar) image, center-crops to a square, scales to
// 256px, and re-encodes as JPEG. It returns a short byte slice safe to write to disk.
func ToJPEG(src []byte) ([]byte, error) {
	if len(src) == 0 {
		return nil, ErrEmpty
	}
	if int64(len(src)) > MaxUploadBytes {
		return nil, ErrTooLarge
	}
	lim := io.LimitReader(bytes.NewReader(src), MaxUploadBytes+1)
	b, err := io.ReadAll(lim)
	if err != nil {
		return nil, err
	}
	if int64(len(b)) > MaxUploadBytes {
		return nil, ErrTooLarge
	}
	img, err := decodeAny(b)
	if err != nil {
		return nil, err
	}
	if err := checkBounds(img); err != nil {
		return nil, err
	}
	cropped := centerSquareToSide(img, outSide)
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, cropped, &jpeg.Options{Quality: jpegQuality}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decodeAny(b []byte) (image.Image, error) {
	// Order: GIF, PNG, JPEG (cheap magic checks), then generic decode (WebP, etc.).
	if isGIF(b) {
		img, err := gif.Decode(bytes.NewReader(b))
		if err != nil {
			return nil, ErrInvalid
		}
		return img, nil
	}
	if isPNG(b) {
		img, err := png.Decode(bytes.NewReader(b))
		if err != nil {
			return nil, ErrInvalid
		}
		return img, nil
	}
	if isJPEG(b) {
		img, err := jpeg.Decode(bytes.NewReader(b))
		if err != nil {
			return nil, ErrInvalid
		}
		return img, nil
	}
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, ErrInvalid
	}
	return img, nil
}

func isGIF(b []byte) bool {
	return len(b) >= 6 && (string(b[:6]) == "GIF89a" || string(b[:6]) == "GIF87a")
}

func isPNG(b []byte) bool { return len(b) >= 8 && string(b[:8]) == "\x89PNG\r\n\x1a\n" }
func isJPEG(b []byte) bool {
	return len(b) >= 3 && b[0] == 0xff && b[1] == 0xd8 && b[2] == 0xff
}

func checkBounds(img image.Image) error {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	if w < 1 || h < 1 {
		return ErrInvalid
	}
	if w > maxEdge || h > maxEdge {
		return ErrTooBig
	}
	if int64(w)*int64(h) > maxPixels {
		return ErrTooBig
	}
	return nil
}

// centerSquareToSide center-crops to a square (min width/height) and scales to side×side.
func centerSquareToSide(src image.Image, side int) *image.NRGBA {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	if w < 1 || h < 1 {
		return image.NewNRGBA(image.Rect(0, 0, 1, 1))
	}
	d := w
	if h < w {
		d = h
	}
	x0 := b.Min.X + (w-d)/2
	y0 := b.Min.Y + (h-d)/2
	srcRect := image.Rect(x0, y0, x0+d, y0+d)
	dst := image.NewNRGBA(image.Rect(0, 0, side, side))
	xdraw.CatmullRom.Scale(dst, dst.Bounds(), src, srcRect, draw.Over, nil)
	return dst
}
