package json

import (
	"encoding/json"
)

type decoder struct {
	decoder *json.Decoder
}

func (d *decoder) SetMapString() {}

func (d *decoder) Decode(value interface{}) error {
	return d.decoder.Decode(value)
}
