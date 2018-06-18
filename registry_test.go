package mapdecor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	ResetRegistry()
	f1 := func(input map[string]interface{}) map[string]interface{} { return input }
	f2 := func(input map[string]interface{}) map[string]interface{} { return input }

	Register(f1, f2)
	assert.Equal(t, 2, RegistryCount())
}

func TestDecorators(t *testing.T) {
	ResetRegistry()
	f1 := func(input map[string]interface{}) map[string]interface{} { return input }

	Register(f1)
	assert.Len(t, Decorators(), 1)
}
