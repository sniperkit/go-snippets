package mapview

import (
	"github.com/gopherjs/gopherjs/js"
)

type Circle struct {
	js.Object
}

func NewCircle(lat, lng, radius float64) *Circle {
	return &Circle{
		Object: L.Call("circle", NewLatLng(lat, lng), radius),
	}
}

func (c *Circle) LatLng() (float64, float64) {
	latLng := c.Call("getLatLng")
	return latLng.Get("lat").Float(), latLng.Get("lng").Float()
}

func (c *Circle) SetLatLng(lat, lng float64) {
	c.Call("setLatLng", NewLatLng(lat, lng))
}

func (c *Circle) Radius() float64 {
	return c.Call("getRadius").Float()
}

func (c *Circle) SetRadius(radius float64) {
	c.Call("setRadius", radius)
}

func (c *Circle) SetStyle(style js.M) {
	c.Call("setStyle", style)
}
