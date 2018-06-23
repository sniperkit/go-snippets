package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fatih/color"
	prettyjson "github.com/hokaccha/go-prettyjson"
)

// struct UnixTime will be used
// as custom field in Person struct
type UnixTime struct {
	time.Time
}

// Custom marshal function for UnixTime.
// It converts Go time to JSON number.
func (t UnixTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time.Unix())
}

// Custom unmarshal function for UnixTime.
// It converts JSON number to Go time.
func (t *UnixTime) UnmarshalJSON(data []byte) error {
	var i int64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	t.Time = time.Unix(i, 0)
	return nil
}

type Person struct {
	Name   string `json:"name"`
	Age    int32  `json:"age"`
	Gender string `json:"gender"`

	// Birthday demonstrates the use of custom
	// type `UnixTime` for specific field (Birthday).
	//
	// Birthday field will be converted using custom
	// UnmarshalJSON from JSON number to Go's time
	Birthday UnixTime `json:"birthday"`

	// EmptyField demonstrates the behavior of a field
	// that's defined in struct but is not in the JSON
	EmptyField string `json:"empty"`

	// A field which value is nil will not be *omitted*
	// from the encoding if it has has omitempty option.
	OmittedField string `json:"omitted,omitempty"`
}

func main() {
	// create function to print colored title
	title := color.New(color.FgGreen).Add(color.Underline).PrintlnFunc()
	note := color.New(color.FgBlue).PrintlnFunc()

	// create byte representation of JSON string
	p := []byte(`{
		"name":"John Doe",
		"age":23,
		"gender":"Male",
		"birthday":761849888,
		"omitted":""
	}`)

	title("\n→ JSON byte to struct (unmarshal):")
	var person Person
	err := json.Unmarshal(p, &person)
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(person)

	title("\n→ Convert Person struct to JSON (marshal):")
	var jsonPerson []byte
	jsonPerson, err = json.Marshal(person)
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(jsonPerson)

	title("\n→ Convert Person struct to JSON (marshal) and cast to string:")
	spew.Dump(string(jsonPerson))

	title("\n→ Pretty-print JSON struct (using hokaccha/go-prettyjson):")
	s, _ := prettyjson.Marshal(person)
	fmt.Println(string(s))
	note("👆  Note that Person field 'OmittedField' is not included in JSON (because omitempty option is set)")
}
