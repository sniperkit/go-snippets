package libkv

import (
	"encoding"
	"encoding/base64"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/docker/libkv/store"
	"github.com/go-mixins/loader"
	"github.com/mitchellh/mapstructure"
)

// Load loads the target from libkv source
func (l *Loader) Load(dest interface{}) error {
	cfg := &mapstructure.DecoderConfig{
		Result:           dest,
		DecodeHook:       decodeHook,
		WeaklyTypedInput: true,
	}
	decoder, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		return loader.Errors.Wrap(err, "creating map decoder")
	}
	data, err := l.getRecursive(l.prefix)
	if err != nil {
		return err
	}
	return loader.Errors.Wrap(decoder.Decode(data), "decoding values")
}

func decodeHook(fromType reflect.Type, toType reflect.Type, data interface{}) (interface{}, error) {
	// decode hook is borrowed from the excellent package
	// "github.com/containous/staert"
	// Copyright (c) 2016 Containous SAS, Emile Vauge, emile@vauge.com

	// custom unmarshaler
	textUnmarshalerType := reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
	if toType.Implements(textUnmarshalerType) {
		object := reflect.New(toType.Elem()).Interface()
		err := object.(encoding.TextUnmarshaler).UnmarshalText([]byte(data.(string)))
		if err != nil {
			return nil, loader.Errors.Wrapf(err, "unmarshaling %v: %v", data, err)
		}
		return object, nil
	}
	switch toType.Kind() {
	case reflect.Ptr:
		if fromType.Kind() == reflect.String {
			if data == "" {
				// default value Pointer
				return make(map[string]interface{}), nil
			}
		}
	case reflect.Slice:
		if fromType.Kind() == reflect.Map {
			// Type assertion
			dataMap, ok := data.(map[string]interface{})
			if !ok {
				return data, loader.Errors.Errorf("input data is not a map : %#v", data)
			}
			// Sorting map
			indexes := make([]int, len(dataMap))
			i := 0
			for k := range dataMap {
				ind, err := strconv.Atoi(k)
				if err != nil {
					return dataMap, loader.Errors.Wrap(err, "converting index")
				}
				indexes[i] = ind
				i++
			}
			sort.Ints(indexes)
			// Building slice
			dataOutput := make([]interface{}, i)
			i = 0
			for _, k := range indexes {
				dataOutput[i] = dataMap[strconv.Itoa(k)]
				i++
			}

			return dataOutput, nil
		} else if fromType.Kind() == reflect.String {
			b, err := base64.StdEncoding.DecodeString(data.(string))
			if err != nil {
				return nil, loader.Errors.Wrap(err, "decoding base64")
			}
			return b, nil
		}
	}
	return data, nil
}

func put(dest map[string]interface{}, path []string, val interface{}) {
	switch len(path) {
	case 0:
		break
	case 1:
		dest[path[0]] = val
	default:
		tmp, ok := dest[path[0]].(map[string]interface{})
		if !ok {
			tmp = make(map[string]interface{})
		}
		put(tmp, path[1:], val)
		dest[path[0]] = tmp
	}
}

func (l *Loader) getRecursive(prefix string) (interface{}, error) {
	pairs, err := l.store.List(prefix)
	if err == store.ErrKeyNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, loader.Errors.Wrapf(err, "getting KV list for %q", prefix)
	}
	if len(pairs) == 0 {
		val, err := l.store.Get(prefix)
		if err != nil {
			return nil, loader.Errors.Wrapf(err, "getting key value for %q", prefix)
		}
		if val != nil {
			return string(val.Value), nil
		}
		return nil, nil
	}
	res := make(map[string]interface{})
	for _, p := range pairs {
		val, err := l.getRecursive(p.Key)
		if err != nil {
			return res, err
		}
		put(res, strings.Split(strings.TrimPrefix(p.Key, prefix+"/"), "/"), val)
	}
	return res, nil
}
