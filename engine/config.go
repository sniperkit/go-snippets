package engine

import (
	"errors"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

// ErrNoConfigDecoder
var ErrNoConfigDecoder = errors.New("config entry has no config decoder")

// ErrConfigDecoderFailed
var ErrConfigDecoderFailed = errors.New("config decoder failed")

// ErrNotAConfigStructure
var ErrNotAConfigStructure = errors.New("not a propper structure")

// parse config files
func parse(path string) (map[string]interface{}, error) {
	// TODO: sanitize path
	var data map[string]interface{}
	tomldata, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := toml.Unmarshal(tomldata, &data); err != nil {
		return nil, err
	}
	return data, nil
}
