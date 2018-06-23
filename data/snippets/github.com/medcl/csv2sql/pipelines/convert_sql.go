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
	"bytes"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/framework/core/util"
	"strings"
)

type ConvertSQLJoint struct {
	pipeline.Parameters
}

func (joint ConvertSQLJoint) Name() string {
	return "convert_sql"
}

const excelBytesKey = "excelBytes"

var rowFormatKey pipeline.ParaKey = "row_format"
var sheetNameKey pipeline.ParaKey = "sheet_name"
var columnNameKey pipeline.ParaKey = "column_name"
var dataFromIndexKey pipeline.ParaKey = "data_start_from_index"

var sqlKey pipeline.ParaKey = "sql"

func (joint ConvertSQLJoint) Process(c *pipeline.Context) error {

	excelBytes := c.MustGetBytes(excelBytesKey)
	r := bytes.NewReader(excelBytes)

	templates := joint.MustGetStringArray(rowFormatKey)
	log.Debug("row templates: ", templates)

	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		log.Error(err)
		return err
	}
	sheetMap := xlsx.GetSheetMap()
	log.Debug("sheets: ", sheetMap)

	colNames := joint.MustGetStringArray(columnNameKey)
	dataOffset := joint.MustGetInt(dataFromIndexKey)

	sheetName := joint.MustGetString(sheetNameKey)

	rows := xlsx.GetRows(sheetName)

	var sqlBuffer bytes.Buffer
	for offset, row := range rows {
		if offset < dataOffset {
			log.Debugf("%v < data offset: %v,　ignore", offset, dataOffset)
			continue
		}

		colMap := map[string]string{}

		hit := false
		for k, colCell := range row {
			if colCell != "" {
				hit = true
			}
			colName := colNames[k]
			colMap[colName] = colCell
			log.Trace("row:", offset, ": ", colName, "-", k, "-", colCell)
		}

		//ignore empty row
		if !hit {
			continue
		}

		for _, x := range templates {
			line := x
			log.Debug("template:", line)
			for k, v := range colMap {
				log.Debug(fmt.Sprintf("<{%v: }>", k), ",", formatString(v))
				line = strings.Replace(line, fmt.Sprintf("<{%v: }>", k), formatString(v), -1)
			}
			log.Debug(line)
			sqlBuffer.WriteString(line)
		}
	}

	sqlMap, ok := c.GetMap(sqlKey)
	if !ok {
		sqlMap = map[string]interface{}{}
	}

	sqlText := sqlBuffer.String()
	sqlMap[sheetName] = sqlText
	c.Set(sqlKey, sqlMap)

	log.Trace(sqlBuffer.String())
	return nil
}

func formatString(str string) string {
	str = util.TrimSpaces(str)
	if str == "" {
		return "NULL"
	}
	str = strings.Replace(str, "\"", "", -1)
	str = strings.Replace(str, "'", "", -1)
	str = fmt.Sprintf("'%s'", str)
	return str
}
