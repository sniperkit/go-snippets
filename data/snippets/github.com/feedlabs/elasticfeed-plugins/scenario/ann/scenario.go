package ann

import (
	"github.com/feedlabs/elasticfeed/common"
	"github.com/feedlabs/elasticfeed/workflow"
//	"github.com/feedlabs/elasticfeed/plugin/model"
)

type config struct {
	common.ElasticfeedConfig `mapstructure:",squash"`

	tpl *workflow.ConfigTemplate
}

type Scenario struct {
	config config
}

func (p *Scenario) Prepare(raws ...interface{}) ([]string, error) {
	return nil, nil
}

func (p *Scenario) Run(data interface{}) (interface{}, error) {
	return data, nil
}

func (p *Scenario) Cancel() {
}
