package mapdecor

import (
	"strings"
	"testing"

	"github.com/linkosmos/mapop"
	"github.com/stretchr/testify/assert"
)

var decorateTests = []struct {
	input     map[string]interface{}
	expected  map[string]interface{}
	decorFunc Decor
}{
	{
		input: map[string]interface{}{
			"key1": nil,
			"key2": "with",
			"val1": "1",
			"val2": "2",
			"val3": "3",
			"val4": "4",
		},
		expected: map[string]interface{}{
			"key1": nil,
			"key2": "with",
			"values": map[string]interface{}{
				"val1": "1",
				"val2": "2",
				"val3": "3",
				"val4": "4",
			},
		},
		decorFunc: func(input map[string]interface{}) (output map[string]interface{}) {
			partitonFunc := func(s string, i interface{}) bool {
				return strings.Contains(s, "val")
			}

			// For first (inputPartitioned[0]) partition we get key-values of val
			// For second (inputPartitioned[1])  partition what is left
			inputPartitioned := mapop.Partition(partitonFunc, input)

			// Assigning values key to first partition
			inputPartitioned[1]["values"] = inputPartitioned[0]

			return inputPartitioned[1]
		},
	},
}

func TestDecorateFromRegistry(t *testing.T) {
	ResetRegistry()

	for _, test := range decorateTests {
		Register(test.decorFunc)
		got := DecorateFromRegistry(test.input)

		assert.Equal(t, test.expected, got)
	}
}

func TestDecorator(t *testing.T) {
	for _, test := range decorateTests {
		got := Decorate(test.input, test.decorFunc)

		assert.Equal(t, test.expected, got)
	}
}
