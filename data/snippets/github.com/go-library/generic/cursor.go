package generic

import (
	"reflect"
)

type M map[string]interface{}
type S []interface{}

func Must(cur *Cursor, err error) *Cursor {
	if err != nil {
		panic(err)
	}

	return cur
}

type Cursor struct {
	// interface
	parent reflect.Value
	// for map or slice
	myKey reflect.Value
}

/* NEW FUNCTIONS */

func NewCursor(root *interface{}) *Cursor {
	c := new(Cursor)
	c.parent = reflect.ValueOf(root).Elem()
	return c
}

/* ACCESS INDEX ELEMENT, IF NOT EXISTED CREATE OR PANIC */

func (c *Cursor) Index(keys ...interface{}) (nextCursor *Cursor, err error) {
	nextCursor = new(Cursor)
	nextCursor.parent = c.parent
	nextCursor.myKey = c.myKey

	for i := range keys {
		k := reflect.ValueOf(keys[i])

		err = nextCursor.prepareToNext(nextCursor.value(), k, false, false)
		if err != nil {
			return nil, err
		}
		nextCursor.parent = nextCursor.value()
		nextCursor.myKey = k
	}

	return nextCursor, nil
}

func (c *Cursor) SetIndex(keys ...interface{}) (nextCursor *Cursor, err error) {
	nextCursor = new(Cursor)
	nextCursor.parent = c.parent
	nextCursor.myKey = c.myKey

	for i := range keys {
		k := reflect.ValueOf(keys[i])

		err = nextCursor.prepareToNext(nextCursor.value(), k, true, true)
		if err != nil {
			return nil, err
		}
		nextCursor.parent = nextCursor.value()
		nextCursor.myKey = k
		if !nextCursor.IsValid() {
			nextCursor.setEmpty()
		}
	}

	return nextCursor, nil
}

/* GETTERS */
func (c *Cursor) Interface() interface{} {
	if c.value().IsValid() {
		return c.value().Interface()
	} else {
		return nil
	}
}

func (c *Cursor) String() string {
	return c.value().String()
}

func (c *Cursor) getIface() reflect.Value {
	var value reflect.Value

	switch c.parent.Kind() {
	case reflect.Interface:
		value = c.parent
	case reflect.Map:
		value = c.parent.MapIndex(c.myKey)
	case reflect.Slice:
		value = c.parent.Index(int(c.myKey.Int()))
	default:
		value = c.parent
	}

	return value
}

func (c *Cursor) IsValid() bool {
	return c.getIface().IsValid()
}

func (c *Cursor) value() reflect.Value {
	var value reflect.Value
	value = c.getIface()
	if value.Kind() == reflect.Interface {
		value = value.Elem()
	}

	return value
}

/* SLICE FUNCTIONS */

func (c *Cursor) Len() int {
	if c.value().Kind() != reflect.Slice {
		panic(Errorf("this value is not slice"))
	}

	return c.value().Len()
}

func (c *Cursor) Push(values ...interface{}) {
	var vars []reflect.Value
	for i := range values {
		if values[i] == nil {
			vars = append(vars, reflect.Zero(c.value().Type().Elem()))
		} else {
			vars = append(vars, reflect.ValueOf(values[i]))
		}
	}

	c.pushValues(vars...)
}

func (c *Cursor) pushValues(values ...reflect.Value) {
	if c.value().Kind() != reflect.Slice {
		panic(Errorf("this value is not slice"))
	}

	c.setValue(reflect.Append(c.value(), values...))
}

func (c *Cursor) Slice(i, j int) (nextCursor *Cursor) {
	v := c.value().Slice(i, j)
	nc := NewCursor(new(interface{}))
	nc.setValue(v)

	return nc
}

func (c *Cursor) Slice3(i, j, k int) (nextCursor *Cursor) {
	v := c.value().Slice3(i, j, k)
	nc := NewCursor(new(interface{}))
	nc.setValue(v)
	return nc
}

/* MAP FUNCTIONS */

func (c *Cursor) Keys() (keys []string) {
	if c.value().Kind() != reflect.Map {
		panic(Errorf("this value is not map"))
	}

	keyValues := c.value().MapKeys()
	for i := range keyValues {
		keys = append(keys, keyValues[i].String())
	}

	return keys
}

/* SETTERS */

func (c *Cursor) Set(value interface{}) {
	if value == nil {
		c.setEmpty()
	} else {
		c.setValue(reflect.ValueOf(value))
	}
}

func (c *Cursor) Delete() {
	switch c.parent.Kind() {
	case reflect.Map:
		c.setValue(reflect.Value{})
	default:
		c.setEmpty()
	}
}

func (c *Cursor) setValue(value reflect.Value) {
	switch c.parent.Kind() {
	case reflect.Interface:
		c.parent.Set(value)
	case reflect.Map:
		c.parent.SetMapIndex(c.myKey, value)
	case reflect.Slice:
		c.parent.Index(int(c.myKey.Int())).Set(value)
	default:
		panic(Errorf("unsupported value kind: %s", c.parent.Kind()))
	}
}

func (c *Cursor) setMap() {
	c.setValue(makeMap())
}

func (c *Cursor) setSlice() {
	c.setValue(makeSlice(0, 0))
}

func (c *Cursor) setEmpty() {
	if c.parent.IsValid() {
		if c.parent.Kind() == reflect.Interface {
			c.setValue(reflect.Zero(c.parent.Type()))
		} else {
			c.setValue(reflect.Zero(c.parent.Type().Elem()))
		}
	} else {
		panic(Errorf("parent value is invalid"))
	}
}

func (c *Cursor) prepareToNext(value, key reflect.Value, permitCreate bool, permitIncrease bool) (err error) {
	var (
		nextValue reflect.Value
		isCreated = false
		vk        = value.Kind()
		kk        = key.Kind()
	)

	//check permission or type

	switch kk {
	case reflect.String, reflect.Int:
	default:
		return Errorf("key should be string or integer")
	}

	switch vk {
	case reflect.Map, reflect.Slice, reflect.Invalid:
	default:
		return Errorf("value should be map, slice or invalid")
	}

	switch {
	case vk == reflect.Map && kk == reflect.String:
	case vk == reflect.Slice && kk == reflect.Int:
	case vk == reflect.Invalid:
		// check create permission
		if !permitCreate {
			return Errorf("implicated creation failure")
		}
	default:
		// check override permission
		return Errorf("value is not compatible with key: value:%v - key:%v", vk, kk)
	}

	switch {
	case kk == reflect.String && vk == reflect.Map:
		// nothing to do
		nextValue = value
	case kk == reflect.String && vk == reflect.Invalid:
		// make new map
		nextValue = makeMap()
		isCreated = true
	case kk == reflect.Int && vk == reflect.Slice:
		idx := int(key.Int())
		if idx >= value.Cap() {
			// check increase permission
			if !permitIncrease {
				return Errorf("out of slice capacity")
			}
			nextValue = makeSlice(idx+1, idx+1)
			reflect.Copy(nextValue, value)
			isCreated = true
		} else if idx < value.Cap() && idx >= value.Len() {
			// check increase permission
			if !permitIncrease {
				return Errorf("out of slice length")
			}
			nextValue = value.Slice(0, idx+1)
			isCreated = true
		}
	case kk == reflect.Int && vk == reflect.Invalid:
		// make new slice
		idx := int(key.Int())
		nextValue = makeSlice(idx+1, idx+1)
		isCreated = true
	}

	if isCreated {
		c.setValue(nextValue)
	}

	return nil
}

func makeMap() reflect.Value {
	var objectInterface map[string]interface{}
	return reflect.MakeMap(reflect.TypeOf(objectInterface))
}

func makeSlice(len int, cap int) reflect.Value {
	var (
		n uint
		i uint
	)

	for i = 1; ; i++ {
		n = 1 << i

		if int(n) >= cap {
			cap = int(n)
			break
		}
	}

	var arrayInterface []interface{}
	return reflect.MakeSlice(reflect.TypeOf(arrayInterface), len, cap)
}
