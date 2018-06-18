package bashconf

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/uoregon-libraries/gopkg/fileutil"
)

// Store populates a tagged structure with the Config data.  Tagged
// settings must be given a bash variable name, and can optionally specify the
// type of "path" or "url" for some basic sanity checking.
func (c *Config) Store(dest interface{}) error {
	var errors = c.readTaggedFields(dest)
	if len(errors) > 0 {
		return fmt.Errorf("invalid configuration: %s", strings.Join(errors, ", "))
	}

	return nil
}

// readTaggedFields iterates over the tagged fields in dest and pulls settings
// from c.  If a tagged field has a type, it's used to process/validate the
// raw string value.
func (c *Config) readTaggedFields(dest interface{}) (errors []string) {
	var rType = reflect.TypeOf(dest).Elem()
	var rVal = reflect.ValueOf(dest).Elem()

	for i := 0; i < rType.NumField(); i++ {
		var sf = rType.Field(i)

		// Ignore fields we can't set, regardless of tagging
		if !rVal.Field(i).CanSet() {
			continue
		}

		// If there's no "setting" tag, we have nothing to do here
		var sKey = sf.Tag.Get("setting")
		if sKey == "" {
			continue
		}

		var val = c.Get(sKey)
		var sType = sf.Tag.Get("type")
		switch sType {
		case "":
			rVal.Field(i).SetString(val)
		case "float":
			var num, err = strconv.ParseFloat(val, 64)
			if err != nil {
				errors = append(errors, fmt.Sprintf("%#v (%#v) is not a valid float", sKey, val))
			}
			rVal.Field(i).SetFloat(num)
		case "int":
			var num, err = strconv.ParseInt(val, 10, 64)
			if err != nil {
				errors = append(errors, fmt.Sprintf("%#v (%#v) is not a valid integer", sKey, val))
			}
			rVal.Field(i).SetInt(num)
		case "url":
			var u, err = url.Parse(val)
			if err != nil {
				errors = append(errors, fmt.Sprintf("%#v (%#v) is not a valid URL: %s", sKey, val, err))
			}
			if u == nil {
				return errors
			}
			if u.Host == "" {
				errors = append(errors, fmt.Sprintf("%#v (%#v) is not a valid URL: missing host", sKey, val))
			}
			if !strings.HasPrefix(u.Scheme, "http") {
				errors = append(errors, fmt.Sprintf("%#v (%#v) is not a valid URL: must be http(s)", sKey, val))
			}
			rVal.Field(i).SetString(val)
		case "path":
			rVal.Field(i).SetString(val)
			if !fileutil.IsDir(val) {
				errors = append(errors, fmt.Sprintf("%#v (%#v) is not a directory", sKey, val))
				continue
			}
		case "file":
			rVal.Field(i).SetString(val)
			if !fileutil.IsFile(val) {
				errors = append(errors, fmt.Sprintf("%#v (%#v) is not a file", sKey, val))
				continue
			}
		}
	}

	return errors
}
