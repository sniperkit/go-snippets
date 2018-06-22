package pipelines

import (
	"github.com/infinitbyte/framework/core/pipeline"
)

type RowStructureJoint struct {
}

func (joint RowStructureJoint) Name() string {
	return "row_structure"
}

func (joint RowStructureJoint) Process(c *pipeline.Context) error {

	return nil
}
