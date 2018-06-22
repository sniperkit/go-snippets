package mapview

import (
	"github.com/gopherjs/gopherjs/js"
)

type LatLng struct {
	js.Object
}

func NewLatLng(lat, lng float64) *LatLng {
	return &LatLng{
		Object: L.Call("latLng", lat, lng),
	}
}
