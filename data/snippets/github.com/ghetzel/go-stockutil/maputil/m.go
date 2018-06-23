package maputil

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/ghetzel/go-stockutil/typeutil"
)

// A Map object (or "M" object) is a utility struct that makes it straightforward to
// work with interface data types that contain map-like data (has a reflect.Kind equal
// to reflect.Map).
type Map struct {
	data interface{}
}

// Create a new Variant map object from the given value (which should be a map of some kind).
func M(data interface{}) *Map {
	if dataV, ok := data.(typeutil.Variant); ok {
		data = dataV.Value
	} else if dataM, ok := data.(*Map); ok {
		return dataM
	} else if dataM, ok := data.(Map); ok {
		return &dataM
	} else if data == nil {
		data = make(map[string]interface{})
	}

	return &Map{
		data: data,
	}
}

// Return the underlying value the M-object was created with.
func (self *Map) Value() interface{} {
	return self.data
}

// Set a value in the Map at the give dot.separated key.
func (self *Map) Set(key string, value interface{}) typeutil.Variant {
	vv := typeutil.V(value)
	self.data = DeepSet(self.data, strings.Split(key, `.`), vv)

	return vv
}

// Retrieve a value from the Map by the given dot.separated key, or return a fallback
// value.  Return values are a typeutil.Variant, which can be easily coerced into
// various types.
func (self *Map) Get(key string, fallbacks ...interface{}) typeutil.Variant {
	if v := DeepGet(self.data, strings.Split(key, `.`), fallbacks...); v != nil {
		return typeutil.Variant{
			Value: v,
		}
	} else {
		return typeutil.Variant{}
	}
}

// Return the value at key as an automatically converted value.
func (self *Map) Auto(key string, fallbacks ...interface{}) interface{} {
	return self.Get(key, fallbacks...).Auto()
}

// Return the value at key as a string.
func (self *Map) String(key string, fallbacks ...interface{}) string {
	return self.Get(key, fallbacks...).String()
}

// Return the value at key interpreted as a Time.
func (self *Map) Time(key string, fallbacks ...interface{}) time.Time {
	return self.Get(key, fallbacks...).Time()
}

// Return the value at key as a bool.
func (self *Map) Bool(key string) bool {
	return self.Get(key).Bool()
}

// Return the value at key as an integer.
func (self *Map) Int(key string, fallbacks ...interface{}) int64 {
	return self.Get(key, fallbacks...).Int()
}

// Return the value at key as a float.
func (self *Map) Float(key string, fallbacks ...interface{}) float64 {
	return self.Get(key, fallbacks...).Float()
}

// Return the value at key as a byte slice.
func (self *Map) Bytes(key string) []byte {
	return self.Get(key).Bytes()
}

// Return the value at key as a slice.  Scalar values will be returned as a slice containing
// only that value.
func (self *Map) Slice(key string) []typeutil.Variant {
	return self.Get(key).Slice()
}

// Return the value at key as a Map.  If the resulting value is nil or not a
// map type, a null Map will be returned.  All values retrieved from a null
// Map will return that type's zero value.
func (self *Map) Map(key string) map[typeutil.Variant]typeutil.Variant {
	return self.Get(key).Map()
}

func (self *Map) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.data)
}
