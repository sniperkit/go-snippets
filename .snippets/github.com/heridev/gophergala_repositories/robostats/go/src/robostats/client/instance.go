package client

import (
	"time"
)

type Instance struct {
	ID        string      `json:"id,omitempty"`
	UserID    string      `json:"user_id,omitempty"`
	ClassID   string      `json:"class_id"`
	Data      interface{} `json:"user_data,omitempty"`
	CreatedAt time.Time   `json:"created_at,omitempty"`
}

type instanceEnvelope struct {
	Instance Instance `json:"deviceInstance"`
}

type instancesEnvelope struct {
	Instances []Instance `json:"deviceInstances"`
}
