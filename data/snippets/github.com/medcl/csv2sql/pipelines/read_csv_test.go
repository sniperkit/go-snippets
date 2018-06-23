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
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/infinitbyte/framework/core/env"
	"github.com/infinitbyte/framework/core/global"
	"github.com/infinitbyte/framework/core/pipeline"
	"testing"
)

func TestExcelize(t *testing.T) {
	xlsx := excelize.NewFile()
	// Create a new sheet.
	index := xlsx.NewSheet("Sheet2")
	// Set value of a cell.
	xlsx.SetCellValue("Sheet2", "A2", "Hello world.")
	xlsx.SetCellValue("Sheet1", "B2", 100)
	// Set active sheet of the workbook.
	xlsx.SetActiveSheet(index)
	// Save xlsx file by the given path.
	err := xlsx.SaveAs("/tmp/Book1.xlsx")
	if err != nil {
		fmt.Println(err)
	}

	xlsx, err = excelize.OpenFile("/tmp/Book1.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	// Get value from cell by given worksheet name and axis.
	cell := xlsx.GetCellValue("Sheet2", "A2")
	fmt.Println(cell)
	// Get all the rows in the Sheet1.
	rows := xlsx.GetRows("Sheet1")
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
}
func TestSeek(t *testing.T) {
	e := env.EmptyEnv()
	global.RegisterEnv(e)
	e.IsDebug = true
	joint := ReadCsvJoint{}
	c := pipeline.Context{}
	c.Set(rowFormatKey, []string{"insert into mytable(a,b,c) values(<{colA: }>,<{colB: }>,<{colC: }>);"})
	c.Set(sheetNameKey, "Sheet1")
	c.Set(dataFromIndexKey, 2)
	c.Set(columnNameKey, []string{"colA", "colB", "colC", "colD", "colE"})
	joint.Process(&c)
}
