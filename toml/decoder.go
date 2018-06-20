package toml

import (
	"github.com/pelletier/go-toml"
)

type decoder struct {
	decoder *toml.Decoder
}

func (d *decoder) SetMapString() {}

func (d *decoder) Decode(value interface{}) error {
	return d.decoder.Decode(value)
}
