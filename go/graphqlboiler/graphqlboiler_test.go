package graphqlboiler

import (
	"fmt"
	"testing"
)

func TestReflectFields(t *testing.T) {
	fields := reflectFields(SampleSchemaPerson{})
	fmt.Println(fields)
}
