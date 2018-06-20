package yaml

import (
	"io"

	"github.com/mickep76/encdec"

	"gopkg.in/yaml.v2"
)

type encoding struct{}

func (c *encoding) NewEncoder(writer io.Writer) encdec.Encoder {
	return &encoder{encoder: yaml.NewEncoder(writer)}
}

func (c *encoding) NewDecoder(reader io.Reader) encdec.Decoder {
	return &decoder{decoder: yaml.NewDecoder(reader)}
}

func init() {
	encdec.Register("yaml", &encoding{})
}
