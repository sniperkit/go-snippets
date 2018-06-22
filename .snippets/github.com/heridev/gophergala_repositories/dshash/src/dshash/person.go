package dshash

import (
	"encoding/json"
)

type Person struct {
	Handler   string
	Locations []string
}

func (p *Person) Unmarshal(b []byte) error {
	return json.Unmarshal(b, p)
}

func (p *Person) Marshal() ([]byte, error) {
	return json.Marshal(p)
}
