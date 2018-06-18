package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	// "io"
)

type Record struct {
	Author Author `json:"author"`
	Title  string `json:"title"`
	URL    string `json:"url"`
}

type Author struct {
	ID    uint64 `json:"id"`
	Email string `json:"email"`
}

type author Author

func (a *Author) UnmarshalJSON(b []byte) (err error) {
	fmt.Println("unmarshalling author")
	j, s, n := author{}, "", uint64(0)
	if err = json.Unmarshal(b, &j); err == nil {
		*a = Author(j)
		return
	}
	if err = json.Unmarshal(b, &s); err == nil {
		a.Email = s
		return
	}
	if err = json.Unmarshal(b, &n); err == nil {
		a.ID = n
	}
	return
}

// func Decode(r io.Reader) (x *Record, err error) {
// 	x = new(Record)
// 	err = json.NewDecoder(r).Decode(x)
// 	return
// }

func main() {
	var r []Record

	encodedJSON :=
		[]byte(`[{
        "author": "attila@attilaolah.eu",
        "title":  "My Blog",
        "url":    "http://attilaolah.eu"
    }, {
        "author": 1234567890,
        "title":  "Westartup",
        "url":    "http://www.westartup.eu"
    }, {
        "author": {
            "id":    1234567890,
            "email": "nospam@westartup.eu"
        },
        "title":  "Westartup",
        "url":    "http://www.westartup.eu"
    }]`)

	encodedStream := bytes.NewReader(encodedJSON)
	// fmt.Println("ERROR:", json.Unmarshal(encodedJSON, &r))
	fmt.Println("ERROR:", json.NewDecoder(encodedStream).Decode(&r))
	fmt.Println("RECORDS:", r)
}
