package json

import (
	"encoding/json"
	"io"

	"github.com/mickep76/encdec"
)

type encoding struct{}

func (c *encoding) NewEncoder(writer io.Writer) encdec.Encoder {
	return &encoder{encoder: json.NewEncoder(writer)}
}

func (c *encoding) NewDecoder(reader io.Reader) encdec.Decoder {
	return &decoder{decoder: json.NewDecoder(reader)}
}

func init() {
	encdec.Register("json", &encoding{})
}
