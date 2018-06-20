package yaml

import "gopkg.in/yaml.v2"

type decoder struct {
	mapString bool
	decoder   *yaml.Decoder
}

func (d *decoder) SetMapString() {
	d.mapString = true
}

func (d *decoder) Decode(value interface{}) error {
	if d.mapString {
		var res interface{}
		if err := d.decoder.Decode(&res); err != nil {
			return err
		}

		*value.(*interface{}) = cleanupMapValue(res)
		return nil
	}

	return d.decoder.Decode(value)
}
