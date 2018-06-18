package csj

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
)

// Reader reads CSJ
type Reader struct {
	csvReader *csv.Reader
}

// New Reader creates a CSJ reader
func NewReader(reader io.Reader) *Reader {
	result := new(Reader)

	result.csvReader = csv.NewReader(reader)
	result.csvReader.LazyQuotes = true
	result.csvReader.FieldsPerRecord = -1

	return result
}

// MakeEmptyObject creates the object to be filled in by json unmarshal
type MakeEmptyObject func() interface{}

// Write when passed an array of MakeEmptyObject fn's
// this will return a an array of objects (as interface{})
// typeArray
func (reader *Reader) Read(typeArray []MakeEmptyObject) (objArray []interface{}, err error) {
	valueArray, err := reader.csvReader.Read()
	if err != nil {
		return
	}

	if len(typeArray) == 0 {
		log.Fatal("typeArray of MakeEmptyObject fn's must have at least one entry")
	}

	objArray = make([]interface{}, len(valueArray))

	for n, valueJson := range valueArray {
		var obj = typeArray[Min(n, len(typeArray)-1)]()

		if err = json.Unmarshal([]byte(valueJson), &obj); err != nil {
			log.Fatal(err)
			return
		}

		objArray[n] = obj
	}

	return
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
