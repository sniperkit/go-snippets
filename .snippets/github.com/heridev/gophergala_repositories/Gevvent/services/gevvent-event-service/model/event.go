package model

import (
	rtypes "github.com/dancannon/gorethink/types"
	"time"
)

type Event struct {
	ID          string    `gorethink:"id,omitempty", json:"id"`
	UserID      string    `gorethink:"user_id", json:"user_id"`
	Name        string    `gorethink:"name", json:"name"`
	Description string    `gorethink:"description", json:"description"`
	When        TimeRange `gorethink:"when", json:"when"`
	Where       Location  `gorethink:"where", json:"where"`
	Private     bool      `gorethink:"private", json:"private"`
	PublicAddr  bool      `gorethink:"public_addr", json:"public_addr"`
}

type Location struct {
	Address string       `gorethink:"address", json:"address"`
	LatLng  rtypes.Point `gorethink:"latlng", json:"latlng"`
}

type TimeRange struct {
	From time.Time `gorethink:"from", json:"from"`
	To   time.Time `gorethink:"to", json:"to"`
}
