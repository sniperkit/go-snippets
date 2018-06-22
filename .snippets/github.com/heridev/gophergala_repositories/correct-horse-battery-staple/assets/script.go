// +build js

package main

import (
	"log"
	"time"

	"github.com/gophergala/correct-horse-battery-staple/common"
	"github.com/gophergala/correct-horse-battery-staple/js/encoding/json"
	"github.com/gophergala/correct-horse-battery-staple/mapview"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket"
	"honnef.co/go/js/dom"
)

var document = dom.GetWindow().Document().(dom.HTMLDocument)

var markers map[int64]*mapview.Marker = make(map[int64]*mapview.Marker, 10)
var mapView *mapview.MapView

func setup() error {
	mapView = mapview.New("map")

	return nil
}

func run() error {
	webSocketClosed := false
	ws, err := websocket.Dial(js.Global.Get("WebSocketAddress").String())
	if err != nil {
		return err
	}
	defer ws.Close()
	ws.AddEventListener("close", false, func(_ js.Object) {
		webSocketClosed = true
	})
	enc := json.NewEncoder(ws)
	dec := json.NewDecoder(ws)

	// Clear all markers (possibly left over from previous connection).
	for _, marker := range markers {
		marker.RemoveFromMap(mapView)
	}

	var bounds *mapview.LatLngBounds
	var lat, lng float64
	var accuracy float64
	foundLocation := false

	shareIcon := document.GetElementByID("share-icon").(*dom.HTMLImageElement)
	shareIcon.AddEventListener("click", false, func(event dom.Event) {
		event.PreventDefault()
		shareBox := document.GetElementByID("share-box").(*dom.HTMLInputElement)
		shareBox.Style().SetProperty("display", "initial", "")
		shareBox.Focus()
		shareBox.Select()
		shareBox.Value = dom.GetWindow().Location().Href

		shareIcon.Style().SetProperty("display", "none", "")

		shareBox.AddEventListener("blur", false, func(event dom.Event) {
			event.PreventDefault()
			shareBox.Style().SetProperty("display", "none", "")
			shareIcon.Style().SetProperty("display", "initial", "")
		})
	})

	go func() {
		for {
			if webSocketClosed {
				break
			}

			time.Sleep(time.Second)

			if !foundLocation {
				continue
			}

			clientState := common.ClientState{
				Name:     document.GetElementByID("message").(*dom.HTMLInputElement).Value,
				Lat:      lat,
				Lng:      lng,
				Accuracy: accuracy,
			}

			err = enc.Encode(clientState)
			if err != nil {
				log.Println("enc.Encode:", err)
				break
			}
		}
	}()

	mapView.OnLocFound(func(loc js.Object) {
		document.GetElementByID("spinner-container").(dom.HTMLElement).Style().SetProperty("display", "none", "")

		foundLocation = true
		latlng := loc.Get("latlng")
		lat = latlng.Get("lat").Float()
		lng = latlng.Get("lng").Float()
		accuracy = loc.Get("accuracy").Float()
	})

	mapView.StartLocate()

	for {
		var msg common.ServerUpdate
		originalIds := make(map[int64]struct{})

		for k := range markers {
			originalIds[k] = struct{}{}
		}

		err = dec.Decode(&msg)
		if err != nil || webSocketClosed {
			log.Println("dec.Decode:", err)
			break
		}

		for i, clientState := range msg.Clients {
			if markers[clientState.Id] == nil {
				markers[clientState.Id] = mapView.AddMarkerWithMessage(clientState.Lat, clientState.Lng, clientState.Accuracy, clientState.Name)
			} else {
				markers[clientState.Id].SetLatLng(clientState.Lat, clientState.Lng)
				markers[clientState.Id].SetMessage(clientState.Name)
				markers[clientState.Id].SetAccuracy(clientState.Accuracy)
				delete(originalIds, clientState.Id)
			}

			if i == 0 {
				bounds = mapview.NewLatLngBounds(
					mapview.NewLatLng(clientState.Lat, clientState.Lng),
					mapview.NewLatLng(clientState.Lat, clientState.Lng))
			} else {
				bounds.Extend(mapview.NewLatLng(clientState.Lat, clientState.Lng))
			}
		}

		for key := range originalIds {
			markers[key].RemoveFromMap(mapView)
		}

		if bounds != nil {
			bounds.Pad(0.10)
			mapView.FitBounds(bounds)
		}

		log.Printf("%#v\n", msg)
	}

	mapView.StopLocate()

	return nil
}

func main() {
	document.AddEventListener("DOMContentLoaded", false, func(_ dom.Event) {
		go func() {
			err := setup()
			if err != nil {
				log.Println(err)
			}

			for {
				err := run()
				if err != nil {
					log.Println(err)
				}

				time.Sleep(5 * time.Second)
			}
		}()
	})
}
