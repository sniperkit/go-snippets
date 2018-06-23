package images_test

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"testing"

	"github.com/sahilm/go-tour/images"
)

func TestImage(t *testing.T) {
	i := images.Image{Width: 1, Height: 2}
	got := render(i)
	want := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAACCAIAAAAW4yFwAAAAFElEQVR4nGKaePMSEwMDAyAAAP//DUcCQYPeckYAAAAASUVORK5CYII="

	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func render(m image.Image) (base64Image string) {
	var buf bytes.Buffer
	err := png.Encode(&buf, m)
	if err != nil {
		panic(err)
	}
	base64Image = base64.StdEncoding.EncodeToString(buf.Bytes())
	return
}
