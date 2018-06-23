package maputil

import (
	_ "encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeepSetNothing(t *testing.T) {
	assert := require.New(t)

	output := make(map[string]interface{})
	output = DeepSet(output, []string{}, "yay").(map[string]interface{})

	assert.Empty(output)
}

func TestDeepSetString(t *testing.T) {
	assert := require.New(t)

	output := make(map[string]interface{})
	testValue := "test-string"

	output = DeepSet(output, []string{"str"}, testValue).(map[string]interface{})

	value, ok := output["str"]
	assert.True(ok)
	assert.Equal(testValue, value)
}

func TestDeepSetBool(t *testing.T) {
	assert := require.New(t)

	output := make(map[string]interface{})
	testValue := true

	output = DeepSet(output, []string{"bool"}, testValue).(map[string]interface{})

	value, ok := output["bool"]
	assert.True(ok)
	assert.Equal(testValue, value)
}

func TestDeepSetArray(t *testing.T) {
	assert := require.New(t)

	output := make(map[string]interface{})
	testValues := []string{"first", "second"}

	for i, tv := range testValues {
		output = DeepSet(output, []string{"top-array", fmt.Sprint(i)}, tv).(map[string]interface{})
	}

	// output = DeepSet(output, []string{"top-array"}, 3.4).(map[string]interface{})

	topArray, ok := output["top-array"]
	assert.True(ok)

	switch topArray.(type) {
	case []interface{}:
		for i, val := range topArray.([]interface{}) {
			assert.Equal(testValues[i], val)
		}
	default:
		t.Errorf("want topArray to be []string, got: %T", topArray)
	}
}

func TestDeepSetNestedMapCreation(t *testing.T) {
	assert := require.New(t)

	output := make(map[string]interface{})
	output = DeepSet(output, []string{"deeply", "nested", "map"}, true).(map[string]interface{})
	output = DeepSet(output, []string{"deeply", "nested", "count"}, 2).(map[string]interface{})

	deeply, ok := output["deeply"]
	assert.True(ok)

	deeplyMap := deeply.(map[string]interface{})

	nested, ok := deeplyMap["nested"]
	assert.True(ok)

	nestedMap := nested.(map[string]interface{})

	_, ok = nestedMap["map"]
	assert.True(ok)

	_, ok = nestedMap["count"]
	assert.True(ok)
}

func TestDiffuseMap(t *testing.T) {
	assert := require.New(t)

	output := make(map[string]interface{})

	output["name"] = "test.thing.name"
	output["enabled"] = true
	output["cool.beans"] = "yep"
	output["tags.0"] = "base"
	output["tags.1"] = "other"
	output["devices.0.name"] = "lo"
	output["devices.1.name"] = "eth0"
	output["devices.1.peers.0"] = "0.0.0.0"
	output["devices.1.peers.1"] = "1.1.1.1"
	output["devices.1.peers.2"] = "2.2.2.2"
	output["devices.1.switch.0.name"] = "aa:bb:cc:dd:ee:ff"
	output["devices.1.switch.0.ip"] = "111.222.0.1"
	output["devices.1.switch.1.name"] = "cc:dd:ee:ff:bb:dd"
	output["devices.1.switch.1.ip"] = "111.222.0.2"

	output, err := DiffuseMap(output, ".")
	assert.NoError(err)

	//  name
	v, _ := output["name"]
	assert.Equal("test.thing.name", v)

	//  enabled
	v, _ = output["enabled"]
	assert.Equal(true, v)

	//  tags[]
	v, ok := output["tags"]
	assert.True(ok)

	assert.Len(v, 2)

	vArray := v.([]interface{})

	assert.Equal("base", vArray[0])
	assert.Equal("other", vArray[1])
}
