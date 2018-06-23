package mapview

import (
	"github.com/gopherjs/gopherjs/js"
)

func accuracyColor(a float64) string {
	switch true {
	case a < 10:
		return "#0055ff"
	case a < 50:
		return "#00ff55"
	case a < 100:
		return "#ff5500"
	}
	return "#ff0000"
}

type Marker struct {
	js.Object
	Message  *Popup
	Lat      float64
	Lng      float64
	Accuracy *Circle
	mapView  *MapView
}

func NewMarker(lat, lng float64) *Marker {
	return &Marker{
		Object:  L.Call("marker", NewLatLng(lat, lng)),
		Message: nil,
		Lat:     lat,
		Lng:     lng,
		mapView: nil,
	}
}

func (m *Marker) SetLatLng(lat, lng float64) {
	m.Lat = lat
	m.Lng = lng
	m.Call("setLatLng", NewLatLng(lat, lng))
	if m.Message != nil {
		m.Message.Call("setLatLng", NewLatLng(lat, lng))
	}
	if m.Accuracy != nil {
		m.Accuracy.SetLatLng(lat, lng)
	}
}

func (m *Marker) AddToMap(mapView *MapView) {
	m.mapView = mapView
	m.Call("addTo", mapView)
	if m.Message != nil {
		mapView.Call("addLayer", m.Message)
	}
	if m.Accuracy != nil {
		mapView.Call("addLayer", m.Accuracy)
	}
}

func (m *Marker) SetMessage(message string) {
	if message != "" {
		if m.Message == nil {
			m.Message = NewPopup(m.Lat, m.Lng)
			if m.mapView != nil {
				m.mapView.Call("addLayer", m.Message)
			}
		}
		m.Message.SetContent("<span class=\"popup\">" + message + "</span>")
	} else {
		if m.Message != nil && m.mapView != nil {
			m.mapView.Call("removeLayer", m.Message)
			m.Message = nil
		}
	}
}

func (m *Marker) SetAccuracy(a float64) {
	if a != 0 {
		if m.Accuracy == nil {
			m.Accuracy = NewCircle(m.Lat, m.Lng, a)
			if m.mapView != nil {
				m.mapView.Call("addLayer", m.Accuracy)
			}
		}
		m.Accuracy.SetRadius(a)
		color := accuracyColor(a)
		m.Accuracy.SetStyle(js.M{
			"color":       color,
			"opacity":     0.5,
			"fillColor":   color,
			"fillOpacity": 0.25,
		})
	} else {
		if m.Accuracy != nil && m.mapView != nil {
			m.mapView.Call("removeLayer", m.Accuracy)
			m.Accuracy = nil
		}
	}
}

func (m *Marker) RemoveFromMap(mapView *MapView) {
	if m.Message != nil {
		mapView.Call("removeLayer", m.Message)
	}

	if m.Accuracy != nil {
		mapView.Call("removeLayer", m.Accuracy)
	}

	mapView.Call("removeLayer", m)
}
