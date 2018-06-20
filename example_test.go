package generic_test

import (
	"encoding/json"
	"fmt"
	"github.com/go-library/generic"
	"log"
)

func ExampleCursor_Index() {
	var (
		err error
		v   interface{}
		c   *generic.Cursor
	)

	stream := `{"results": [{"name": "foo" }]}`

	err = json.Unmarshal([]byte(stream), &v)
	if err != nil {
		log.Fatal(err)
	}

	c = generic.NewCursor(&v)
	generic.Must(c.Index("results", 0, "name")).Set("bar")

	fmt.Println(generic.Must(c.Index("results", 0, "name")))
	fmt.Printf("%+v\n", v)

	// Output:
	// bar
	// map[results:[map[name:bar]]]
}
