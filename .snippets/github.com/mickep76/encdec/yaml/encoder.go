package yaml

import (
	"gopkg.in/yaml.v2"
)

type encoder struct {
	encoder *yaml.Encoder
}

func (e *encoder) SetIndent(indent string) {}

func (e *encoder) Encode(value interface{}) error {
	return e.encoder.Encode(value)
}
