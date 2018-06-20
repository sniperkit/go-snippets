package json

import "encoding/json"

type encoder struct {
	encoder *json.Encoder
}

func (e *encoder) SetIndent(indent string) {
	e.encoder.SetIndent("", indent)
}

func (e *encoder) Encode(value interface{}) error {
	return e.encoder.Encode(value)
}
