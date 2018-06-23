package typeutil

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/ghetzel/go-stockutil/utils"
)

// Represents an interface type with helper functions for making it easy to do
// type conversions.
type Variant struct {
	Value interface{}
}

// Shortcut for creating a Variant.
func V(value interface{}) Variant {
	return Variant{
		Value: value,
	}
}

// Returns whether the underlying value is a zero value.
func (self Variant) IsZero() bool {
	return IsZero(self.Value)
}

// Return the value as a string, or an empty string if the value could not be converted.
func (self Variant) String() string {
	if v, err := utils.ConvertToString(self.Value); err == nil {
		return v
	} else {
		return ``
	}
}

// Return true if the value can be interpreted as a boolean true value, or false otherwise.
func (self Variant) Bool() bool {
	if v, err := utils.ConvertToBool(self.Value); err == nil {
		return v
	} else {
		// use a more relaxed set of values for determining "true" because
		// the user has very explicitly asked us to try
		switch strings.ToLower(fmt.Sprintf("%v", self.Value)) {
		case `on`, `1`, `yes`, `active`, `online`:
			return true
		}

		return false
	}
}

// Return the value as a float if it can be interpreted as such, or 0 otherwise.
func (self Variant) Float() float64 {
	if v, err := utils.ConvertToFloat(self.Value); err == nil {
		return v
	} else {
		return 0
	}
}

// Return the value as an integer if it can be interpreted as such, or 0 otherwise. Float values
// will be truncated to integers.
func (self Variant) Int() int64 {
	if v, err := utils.ConvertToFloat(self.Value); err == nil {
		return int64(v)
	} else {
		return 0
	}
}

// Return the value as a slice of Variants. Scalar types will return a slice containing
// a single Variant element representing the value.
func (self Variant) Slice() []Variant {
	values := make([]Variant, 0)

	utils.SliceEach(self.Value, func(_ int, v interface{}) error {
		values = append(values, Variant{
			Value: v,
		})
		return nil
	})

	return values
}

// Return the value automaticall converted to the appropriate type.
func (self Variant) Auto() interface{} {
	return utils.Autotype(self.Value)
}

// Return the value as a time.Time if it can be interpreted as such, or zero time otherwise.
func (self Variant) Time() time.Time {
	if v, err := utils.ConvertToTime(self.Value); err == nil {
		return v
	} else {
		return time.Time{}
	}
}

// Return the value at key as a byte slice.
func (self Variant) Bytes() []byte {
	if v, err := utils.ConvertToBytes(self.Value); err == nil {
		return v
	} else {
		return []byte{}
	}
}

// Return the value as a map[Variant]Variant if it can be interpreted as such, or nil otherwise.
func (self Variant) Map() map[Variant]Variant {
	output := make(map[Variant]Variant)

	if IsMap(self.Value) {
		mapV := reflect.ValueOf(self.Value)

		for _, key := range mapV.MapKeys() {
			if key.CanInterface() {
				if value := mapV.MapIndex(key); value.CanInterface() {
					output[V(key.Interface())] = V(value.Interface())
				}
			}
		}
	}

	return output
}

// Satisfy the json.Marshaler interface
func (self Variant) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Auto())
}
