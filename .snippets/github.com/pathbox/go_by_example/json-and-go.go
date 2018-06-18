import (
	"encoding/json"
)

func Marshal(v interface{}) ([]byte, error)

type Message struct {
	Name string
	Body string
	Time int64
}

m := Message{"Alice", "Hello", 1294706395881547000}

b, err := json.Marshal(m)
b == []byte(`{"Name":"Alice","Body":"Hello","Time":1294706395881547000}`)

func Unmarshal(data []byte, v interface{}) error

var m Message

err := json.Unmarshal(b, &m)

m = Message{
    Name: "Alice",
    Body: "Hello",
    Time: 1294706395881547000,
}

var i interface{}

i = "a string"
i = 2026
i = 2.333

r := i.(float64)
fmt.Println("the circle's area", math.Pi*r*r)

switch v := i.(type) {
case int:
    fmt.Println("twice i is", v*2)
case float64:
    fmt.Println("the reciprocal of i is", 1/v)
case string:
    h := len(v) / 2
    fmt.Println("i swapped by halves is", v[h:]+v[:h])
default:
    // i isn't one of the types above
}

b := []byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`)  // json is a string

var f interface{}

err := json.Unmarshal(b, &f)

f = map[string]interface{}{
    "Name": "Wednesday",
    "Age":  6,
    "Parents": []interface{}{
        "Gomez",
        "Morticia",
    },
}

m := f.(map[string]interface{})

for k, v := range m {
	switch vv := v.(type) {
	case string:
		fmt.Println(k, "is string", vv)
	case int:
		fmt.Println(k, "is int", vv)
	case []interface{}:
		fmt.Println(k, "is an array:")
        for i, u := range vv {
            fmt.Println(i, u)
        }
  default:
    fmt.Println(k, "is of a type I don't know how to handle")
	}
}

type FamilyMember struct {
	Name string
	Age int
	Parents []string
}

var m FamilyMember
err := json.Unmarshal(b, &m)

type IncomingMessage struct {
	Cmd *Command
	Msg *Message
}

func NewDecoder(r io.Reader) *Decorder
func NewEncoder(r io.Writer) *Encoder

package main

import (
    "encoding/json"
    "log"
    "os"
)

func main(){
	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)
	for {
		var v map[string]interface{}
		if err := dec.Decorder(&v); err != nil {
			log.Println(err)
			return
		}
		for k := range v {
			if k!= "Name" {
				delete(v, k)
			}
		}
		if err := enc.Encode(&v); err != nil {
			log.Println(err)
		}
	}
}

// Due to the ubiquity of Readers and Writers, these Encoder and Decoder types can be used in a broad 
// range of scenarios, such as reading and writing to HTTP connections, WebSockets, or files.




































