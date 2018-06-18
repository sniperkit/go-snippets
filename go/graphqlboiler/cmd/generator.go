package main

import (
	"os"

	"github.com/wzulfikar/lab/go/graphqlboiler"
)

// sample struct
type Person struct {
	Name    string
	Age     int
	married bool `json:"-"`
}

// Generate graphql boilerplate from given schemas
func main() {
	schemas := []interface{}{
		Person{},
	}

	for _, schema := range schemas {
		graphqlboiler.Boil(graphqlboiler.Tpl{
			RootResolver: "RootResolver",
			Schema:       schema,
		}, os.Getenv("RESOLVERS_DIR"))
	}
}
