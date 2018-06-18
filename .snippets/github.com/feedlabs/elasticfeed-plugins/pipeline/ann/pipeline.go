package ann

import (
	"time"
	"math/rand"

	"github.com/feedlabs/elasticfeed/common"
	"github.com/feedlabs/elasticfeed/workflow"
	//	"github.com/feedlabs/elasticfeed/plugin/model"
)

type DataA struct {
	Data string
}

type config struct {
	common.ElasticfeedConfig `mapstructure:",squash"`

	tpl *workflow.ConfigTemplate
}

type Pipeline struct {
	config config
}

func (p *Pipeline) Prepare(raws ...interface{}) ([]string, error) {
	return nil, nil
}

func (p *Pipeline) Run(data interface{}) (interface{}, error) {
	//	return RandomAnimator(data), nil
	return data, nil
}

func (p *Pipeline) Cancel() {
}

func RandomAnimator(data interface{}) interface{} {

	_data := make(map[int64]int64, len(data.(map[interface{}]interface{})))

	for i, v := range data.(map[interface{}]interface{}) {
		_data[i.(int64)] = v.(int64)*100
	}

	// PIPE DELAY SIMULATION

	amt := time.Duration(rand.Intn(1))
	time.Sleep(amt * time.Millisecond)

	return _data
}
