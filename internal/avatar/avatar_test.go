package avatar

import (
	"bytes"
	"image"
	"image/png"
	"testing"
)

func TestToJPEG_minPNG(t *testing.T) {
	t.Parallel()
	img := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	jpeg, err := ToJPEG(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if len(jpeg) < 100 {
		t.Fatalf("jpeg output tiny: %d bytes", len(jpeg))
	}
}
