package maputil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type mapTestSubstruct struct {
	Name  string
	Value interface{}
}

type mapTestStruct struct {
	Name    string
	NestedP *mapTestSubstruct
	Nested  mapTestSubstruct
}

func TestGetNil(t *testing.T) {
	assert := require.New(t)

	input := make(map[string]interface{})
	level1 := make(map[string]interface{})

	level1["nilvalue"] = nil

	input["test"] = level1

	assert.Nil(DeepGet(input, []string{"test", "nilvalue"}, "nope"))
	assert.Nil(DeepGet(input, []string{"test", "nilvalue"}, nil))
}

func TestDeepGetScalar(t *testing.T) {
	assert := require.New(t)

	input := make(map[string]interface{})

	input = DeepSet(input, []string{"deeply", "nested", "value"}, 1.4).(map[string]interface{})

	assert.NotNil(DeepGet(input, []string{"deeply", "nested", "value"}, nil))
	assert.Equal(true, DeepGet(input, []string{"deeply", "nested", "value2"}, true))
	assert.Equal(`fallback`, DeepGet(input, []string{"deeply", "nested", "value2"}, "fallback"))
}

func TestDeepGetArrayElement(t *testing.T) {
	input := make(map[string]interface{})

	input = DeepSet(input, []string{"tags", "0"}, "base").(map[string]interface{})
	input = DeepSet(input, []string{"tags", "1"}, "other").(map[string]interface{})

	if v := DeepGet(input, []string{"tags", "0"}, nil); v != "base" {
		t.Errorf("%s\n", v)
	}

	if v := DeepGet(input, []string{"tags", "1"}, nil); v != "other" {
		t.Errorf("%s\n", v)
	}
}

func TestDeepGetMapKeyInArray(t *testing.T) {
	input := make(map[string]interface{})

	input = DeepSet(input, []string{"devices", "0", "name"}, "lo").(map[string]interface{})
	input = DeepSet(input, []string{"devices", "1", "name"}, "eth0").(map[string]interface{})

	if v := DeepGet(input, []string{"devices", "0", "name"}, nil); v != "lo" {
		t.Errorf("%s\n", v)
	}

	if v := DeepGet(input, []string{"devices", "1", "name"}, nil); v != "eth0" {
		t.Errorf("%s\n", v)
	}
}

func TestDeepGetMapKeyInDeepArray(t *testing.T) {
	input := make(map[string]interface{})

	input = DeepSet(input, []string{"devices", "0", "switch", "peers", "0"}, "0.0.0.0").(map[string]interface{})
	input = DeepSet(input, []string{"devices", "0", "switch", "peers", "1"}, "0.0.1.1").(map[string]interface{})
	input = DeepSet(input, []string{"devices", "1", "switch", "peers", "0"}, "1.1.0.0").(map[string]interface{})
	input = DeepSet(input, []string{"devices", "1", "switch", "peers", "1"}, "1.1.1.1").(map[string]interface{})

	if v := DeepGet(input, []string{"devices", "0", "switch", "peers", "0"}, nil); v != "0.0.0.0" {
		t.Errorf("%s\n", v)
	}

	if v := DeepGet(input, []string{"devices", "0", "switch", "peers", "1"}, nil); v != "0.0.1.1" {
		t.Errorf("%s\n", v)
	}

	if v := DeepGet(input, []string{"devices", "1", "switch", "peers", "0"}, nil); v != "1.1.0.0" {
		t.Errorf("%s\n", v)
	}

	if v := DeepGet(input, []string{"devices", "1", "switch", "peers", "1"}, nil); v != "1.1.1.1" {
		t.Errorf("%s\n", v)
	}
}

func TestDeepGetBool(t *testing.T) {
	assert := require.New(t)
	var input interface{}

	input = make(map[string]interface{})

	input = DeepSet(input, []string{"deeply", "nested", "value"}, true)
	input = DeepSet(input, []string{"deeply", "nested", "thing"}, "nope")

	assert.True(DeepGetBool(input, []string{"deeply", "nested", "value"}))
	assert.False(DeepGetBool(input, []string{"deeply", "nested", "other"}))
	assert.False(DeepGetBool(input, []string{"deeply", "nested", "nope"}))
}

func TestDeepGetMapInMap(t *testing.T) {
	assert := require.New(t)

	in := map[string]interface{}{
		`ok`: true,
		`always`: map[string]interface{}{
			`finishing`: map[string]interface{}{
				`each_others`: `sentences`,
			},
		},
	}

	assert.Equal(`sentences`, DeepGet(in, []string{`always`, `finishing`, `each_others`}))
	assert.Nil(DeepGet(in, []string{`always`, `finishing`, `each_others`, `sandwiches`}))
}

func TestDeepStructs(t *testing.T) {
	assert := require.New(t)

	in := &mapTestStruct{
		Name: `toplevel`,
		NestedP: &mapTestSubstruct{
			Name:  `one-ptr`,
			Value: true,
		},
		Nested: mapTestSubstruct{
			Name:  `one-value`,
			Value: 3.14,
		},
	}

	assert.Equal(`toplevel`, DeepGet(in, []string{`Name`}))
	assert.Equal(`one-ptr`, DeepGet(in, []string{`NestedP`, `Name`}))
	assert.Equal(true, DeepGet(in, []string{`NestedP`, `Value`}))
	assert.Equal(`one-value`, DeepGet(in, []string{`Nested`, `Name`}))
	assert.Equal(float64(3.14), DeepGet(in, []string{`Nested`, `Value`}))
}
