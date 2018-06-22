package mapview

import (
	"github.com/gopherjs/gopherjs/js"
)

const (
	tilesUrl = "http://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
)

var L = js.Global.Get("L")

type MapView struct {
	js.Object
}

func New(id string) *MapView {
	options := make(map[string]interface{})
	options["closePopupOnClick"] = false
	mapView := L.Call("map", id, options)
	L.Call("tileLayer", tilesUrl).Call("addTo", mapView)

	return &MapView{
		Object: mapView,
	}
}

func (mv *MapView) SetView(lat, lng float64, zoom int) {
	mv.Call("setView", NewLatLng(lat, lng), zoom)
}

func (mv *MapView) AddMarker(lat, lng float64) *Marker {
	marker := NewMarker(lat, lng)
	marker.AddToMap(mv)
	return marker
}

func (mv *MapView) RemoveMarker(marker *Marker) {
	mv.Call("removeLayer", marker)
}

func (mv *MapView) AddMarkerWithMessage(lat, lng, accuracy float64, msg string) *Marker {
	marker := NewMarker(lat, lng)
	marker.SetMessage(msg)
	marker.SetAccuracy(accuracy)
	marker.AddToMap(mv)
	return marker
}

func (mv *MapView) StartLocate() {
	mv.Call("locate", js.M{
		"watch":              true,
		"enableHighAccuracy": true,
	})
}

func (mv *MapView) StopLocate() {
	mv.Call("stopLocate")
}

func (mv *MapView) OnLocFound(cb func(js.Object)) {
	mv.Call("on", "locationfound", cb)
}

func (mv *MapView) FitBounds(bounds *LatLngBounds) {
	mv.Call("fitBounds", bounds)
}
