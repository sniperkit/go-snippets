package csj

import (
	"encoding/csv"
	"encoding/json"
	"io"
)

// Writer writes CSJ
type Writer struct {
	csvWriter *csv.Writer
}

// NewWriter creates a writer
func NewWriter(writer io.Writer) *Writer {
	result := new(Writer)

	result.csvWriter = csv.NewWriter(writer)

	return result
}

// Write when passed an array of objects
// this will write a line of comma-deliminated json
func (writer *Writer) Write(objArray []interface{}) {
	outArray := make([]string, len(objArray))

	for n, obj := range objArray {
		objJsonBytes, err := json.Marshal(obj)
		if err != nil {
			panic(err)
		}
		outArray[n] = string(objJsonBytes)
	}

	writer.csvWriter.Write(outArray)
	writer.csvWriter.Flush()
}
