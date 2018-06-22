package plugin

import (
	. "github.com/infinitbyte/framework/core/config"
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/medcl/csv2sql/api"
	"github.com/medcl/csv2sql/config"
	"github.com/medcl/csv2sql/pipelines"
	"github.com/medcl/csv2sql/ui"
)

type CSV2SQLPlugin struct {
}

func (this CSV2SQLPlugin) Name() string {
	return "csv2sql"
}

var (
	appConfig = config.AppConfig{}
)

func (module CSV2SQLPlugin) Start(cfg *Config) {

	cfg.Unpack(&appConfig)

	config.SetAppConfig(appConfig)

	//register UI
	if appConfig.UIEnabled {
		ui.InitUI()
	}

	api.InitAPI()

	//register pipeline joints
	pipeline.RegisterPipeJoint(pipelines.LoggingJoint{})
	pipeline.RegisterPipeJoint(pipelines.ReadCsvJoint{})
	pipeline.RegisterPipeJoint(pipelines.ConvertSQLJoint{})
	pipeline.RegisterPipeJoint(pipelines.ImportSQLJoint{})
}

func (module CSV2SQLPlugin) Stop() error {
	return nil
}
