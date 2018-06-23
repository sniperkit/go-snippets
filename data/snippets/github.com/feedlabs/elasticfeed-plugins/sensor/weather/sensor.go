package weather

import (
	"github.com/feedlabs/elasticfeed/common"
	"github.com/feedlabs/elasticfeed/workflow"
//	"github.com/feedlabs/elasticfeed/plugin/model"
)

type config struct {
	common.ElasticfeedConfig `mapstructure:",squash"`

	tpl *workflow.ConfigTemplate
}

type Sensor struct {
	config config
}

func (p *Sensor) Prepare(raws ...interface{}) ([]string, error) {
	return nil, nil
}

func (p *Sensor) Run(data interface{}) (interface{}, error) {
	return data, nil
}

func (p *Sensor) Cancel() {
}


