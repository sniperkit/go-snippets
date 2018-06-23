package event

import "encoding/json"

var Events chan Event

func init() {
	Events = make(chan Event)
}

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func (self Event) ToJSON() string {
	b, _ := json.Marshal(self)
	return string(b)
}
