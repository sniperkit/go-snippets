package json

import (
	"bytes"
	"encoding/json"
)

type Coder struct{}

func (c *Coder) Marshal(v interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(v)
	return buf.Bytes(), err
}

func (c *Coder) Unmarshal(data []byte, v interface{}) error {
	buf := bytes.NewBuffer(data)
	return json.NewDecoder(buf).Decode(v)
}
