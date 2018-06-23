package mapview

import (
	"github.com/gopherjs/gopherjs/js"
)

type Popup struct {
	js.Object
}

func NewPopup(lat, lng float64) *Popup {
	options := make(map[string]interface{})
	options["offset"] = NewPoint(0, -24)
	options["closeButton"] = false
	popup := &Popup{
		Object: L.Call("popup", options),
	}
	popup.Call("setLatLng", NewLatLng(lat, lng))
	return popup
}

func (popup *Popup) SetContent(msg string) {
	popup.Call("setContent", msg)
}
