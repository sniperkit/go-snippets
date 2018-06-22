package gopherpaint

import "testing"
import "image/color"

func TestDistance(t *testing.T) {
	d := distance([]int{5, 5}, 5, 5)
	if 0 != d {
		t.Errorf("Expected 0, given %v", d)
	}

	d = distance([]int{2, 0}, 0, 0)
	if 2 != d {
		t.Errorf("Expected 2, given %v", d)
	}
}

func TestColorMean(t *testing.T) {
	colores := []color.Color{
		color.Black,
		color.Black,
		color.Black,
		color.Black,
	}
	media := colorMean(colores)
	r, g, b, a := media.RGBA()
	if r != 0 || g != 0 || b != 0 || a != 255 {
		t.Errorf("Expected 0 0 0 255, given %v (%v %v %v %v)", media, r, g, b, a)
	}
}
