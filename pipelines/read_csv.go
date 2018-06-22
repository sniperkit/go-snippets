/*
Copyright 2016 Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pipelines

import (
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/framework/core/util"
)

var fileName pipeline.ParaKey = "file_name"

type ReadCsvJoint struct {
	pipeline.Parameters
}

func (joint ReadCsvJoint) Name() string {
	return "read_csv"
}

func (joint ReadCsvJoint) Process(c *pipeline.Context) error {

	excelBytes, err := util.FileGetContent(joint.MustGetString(fileName))
	if err != nil {
		log.Error(err)
		panic(err)
		return err
	}

	c.Set(excelBytesKey, excelBytes)

	return nil
}
