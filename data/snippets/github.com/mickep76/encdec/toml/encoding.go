package toml

import (
	"io"

	"github.com/pelletier/go-toml"

	"github.com/mickep76/encdec"
)

type encoding struct{}

func (c *encoding) NewEncoder(writer io.Writer) encdec.Encoder {
	return &encoder{encoder: toml.NewEncoder(writer)}
}

func (c *encoding) NewDecoder(reader io.Reader) encdec.Decoder {
	return &decoder{decoder: toml.NewDecoder(reader)}
}

func init() {
	encdec.Register("toml", &encoding{})
}
