package client

import (
	"time"
)

type Event struct {
	ID         string      `json:"id,omitempty"`
	UserID     string      `json:"user_id,omitempty"`
	ClassID    string      `json:"class_id,omitempty"`
	InstanceID string      `json:"instance_id,omitempty"`
	SessionID  string      `json:"session_id"`
	Data       interface{} `json:"user_data,omitempty"`
	LocalTime  int         `json:"local_time"`
	LatLng     [2]float64  `json:"latlng,omitempty"`
	CreatedAt  time.Time   `json:"created_at,omitempty"`
}

type eventEnvelope struct {
	Event Event `json:"deviceEvent"`
}

type eventsEnvelope struct {
	Events []Event `json:"deviceEvents"`
}
