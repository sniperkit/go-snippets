package mapview

import (
	"github.com/gopherjs/gopherjs/js"
)

type LatLngBounds struct {
	js.Object
}

func NewLatLngBounds(southWest, northEast *LatLng) *LatLngBounds {
	return &LatLngBounds{
		Object: L.Call("latLngBounds", southWest, northEast),
	}
}

func (llb *LatLngBounds) Extend(point *LatLng) {
	llb.Call("extend", point)
}

func (llb *LatLngBounds) Pad(percentage float64) {
	llb.Object = llb.Call("pad", percentage)
}
