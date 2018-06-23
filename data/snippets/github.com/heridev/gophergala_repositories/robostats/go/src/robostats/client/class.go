package client

import (
	"time"
)

type Class struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	APIKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
}

type classEnvelope struct {
	Class Class `json:"deviceClass"`
}

type classesEnvelope struct {
	Classes []Class `json:"deviceClasses"`
}
