package client

import (
	"time"
)

type Session struct {
	ID         string      `json:"id,omitempty"`
	UserID     string      `json:"user_id,omitempty"`
	ClassID    string      `json:"class_id,omitempty"`
	InstanceID string      `json:"instance_id"`
	SessionKey string      `json:"session_key,omitempty"`
	Data       interface{} `json:"user_data,omitempty"`
	StartTime  time.Time   `json:"start_time"`
	EndTime    time.Time   `json:"end_time,omitempty"`
	CreatedAt  time.Time   `json:"created_at,omitempty"`

	internalTime int
	client       *Client
}

type sessionEnvelope struct {
	Session Session `json:"deviceSession"`
}

type sessionsEnvelope struct {
	Sessions []Session `json:"deviceSessions"`
}
