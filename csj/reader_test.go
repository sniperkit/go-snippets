package csj

import (
	"strings"
	"testing"
)

func NewString() interface{} {
	var s string
	return &s
}

func NewStringArray() interface{} {
	var stringArray = make([]string, 0)
	return &stringArray
}

func TestRead(t *testing.T) {
	validCsjString := `"""key1""","[""val1-1"",""val1-2""]"
"""key2""","[""val2-1"",""val2-2"",""val2-3""]"
`

	rowTypesArray := []MakeEmptyObject{NewString, NewStringArray}

	strReader := strings.NewReader(validCsjString)
	reader := NewReader(strReader)

	valueArray, err := reader.Read(rowTypesArray)

	if err != nil ||
		(*valueArray[0].(*string)) != "key1" ||
		(*valueArray[1].(*[]string))[0] != "val1-1" ||
		(*valueArray[1].(*[]string))[1] != "val1-2" {
		t.Error("The csj.Reader did not read the 1st line's object array correctly")
	}

	valueArray, err = reader.Read(rowTypesArray)

	if err != nil ||
		(*valueArray[0].(*string)) != "key2" ||
		(*valueArray[1].(*[]string))[0] != "val2-1" ||
		(*valueArray[1].(*[]string))[1] != "val2-2" ||
		(*valueArray[1].(*[]string))[2] != "val2-3" {
		t.Error("The csj.Reader did not read the 2nd line's object array correctly")
	}
}
