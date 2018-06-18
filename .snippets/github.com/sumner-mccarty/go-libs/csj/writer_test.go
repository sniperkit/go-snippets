package csj

import (
	"bytes"
	"testing"
)

func TestWrite(t *testing.T) {
	validCsjString := `"""key1""","[""val1-1"",""val1-2""]"
"""key2""","[""val2-1"",""val2-2"",""val2-3""]"
`

	var buffer bytes.Buffer
	csjWriter := NewWriter(&buffer)

	lineArray := make([]interface{}, 2)

	lineArray[0] = "key1"

	var stringArray = make([]string, 2)
	stringArray[0] = "val1-1"
	stringArray[1] = "val1-2"
	lineArray[1] = stringArray
	csjWriter.Write(lineArray)

	lineArray[0] = "key2"

	stringArray = make([]string, 3)
	stringArray[0] = "val2-1"
	stringArray[1] = "val2-2"
	stringArray[2] = "val2-3"
	lineArray[1] = stringArray
	csjWriter.Write(lineArray)

	if buffer.String() != validCsjString {
		t.Error("The csj.Writer did not write the lines correctly")
	}
}
