package plugin

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var defCfg = PList{
	{
		Plugin: "gocolor",
		Uri:    "github.com/gizak/gopp/plugin/gocolor"},
	{
		Plugin: "parsepanic",
		Uri:    "github.com/gizak/gopp/plugin/parsepanic",
		Deps:   []string{"github.com/maruel/panicparse"}},
	{
		Plugin: "rmpkg",
		Uri:    "github.com/gizak/gopp/plugin/rm",
	}}

var ConfigFile = os.Getenv("HOME") + "/.gopp"

func cfg() (PList, error) {
	pl := defCfg
	cfg := Config{}

	data, err := ioutil.ReadFile(ConfigFile)
	if !os.IsNotExist(err) {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, err
		}
		return cfg.Plugins, nil
	}
	return pl, os.ErrNotExist
}
