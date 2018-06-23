package pipelines

import (
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/framework/core/util"
	"time"
)

type LoggingJoint struct {
}

func (joint LoggingJoint) Name() string {
	return "logging"
}

func (joint LoggingJoint) Process(c *pipeline.Context) error {

	sqlTextMap := c.MustGetMap(sqlKey)

	util.FileAppendNewLine("log/executed_sql.log", "")
	util.FileAppendNewLine("log/executed_sql.log", time.Now().String())

	for k, sqlText := range sqlTextMap {
		util.FileAppendNewLine("log/executed_sql.log", "process data sheet: "+k)
		util.FileAppendNewLine("log/executed_sql.log", sqlText.(string))
		util.FileAppendNewLine("log/executed_sql.log", "")
	}

	return nil
}
