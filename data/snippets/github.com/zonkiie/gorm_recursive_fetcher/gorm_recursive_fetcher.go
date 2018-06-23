/****** MIT License **********
Copyright (c) 2017 Datzer Rainer
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
***************************/

package gorm_recursive_fetcher

import (
	"github.com/jinzhu/gorm"
	"fmt"
	"reflect"
	"regexp"
	"encoding/json"
	"encoding/xml"
	"os"
	"strings"
)

func PStdErr(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func JsonMarshal(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		return "Error"
	}
	return string(b[:])
}

func XmlMarshal(data interface{}) string {
	b, err := xml.Marshal(data)
	if err != nil {
		return "Error"
	}
	return string(b[:])
}

// This function fetches all related objects from a given object in the data parameter.
// The struct must be fully tagged, we don't recognize automatically related IDs and so on.
// The function works only with not combined keys.
// Every field which should be fetched must be tagged with:
// walkrec:"true" gorm:"ForeignKey:ID;AssociationForeignKey:ForeignKey"
// See: http://stackoverflow.com/questions/24537525/reflect-value-fieldbyname-causing-panic
// See: http://stackoverflow.com/questions/34493062/how-to-reflect-struct-recursive-in-golang
func fetchRec(db *gorm.DB, data interface{}) {
	// With data *rs: Type: *main.rs
	// With data interface{}: *main.rs
	var ref reflect.Value
	if reflect.TypeOf(data).Kind() == reflect.Struct {
		ref = reflect.ValueOf(data)
	} else if reflect.TypeOf(data).Kind() == reflect.Ptr {
		ref = reflect.Indirect(reflect.ValueOf(data))
	}
	if ref.Type().Kind() == reflect.Slice {
		for i := 0; i < ref.Len(); i++ {
			if ref.Index(i).Type().Kind() == reflect.Ptr {
				fetchRec(db, ref.Index(i).Elem().Addr().Interface())
			} else if ref.Index(i).Type().Kind() == reflect.Struct {
				// What should we do here?
			}
		}
		
	} else if ref.Type().Kind() == reflect.Struct {
		for i := 0; i < ref.NumField(); i++ {
			var IDFieldRaw string
			var IDFields []string
			var RefFieldRaw string
			var RefFields []string
			var re *regexp.Regexp
			var matches []string
			
			if ref.Field(i).CanAddr() && strings.EqualFold(ref.Type().Field(i).Tag.Get("walkrec"), "true") {
				gormflags := ref.Type().Field(i).Tag.Get("gorm")
				if gormflags == "" {
					panic("No gorm flags found!")
				} else {
					re = regexp.MustCompile(`\bForeignKey:([a-zA-Z0-9_,]+)\b`)
					matches = re.FindStringSubmatch(gormflags)
					if len(matches) == 2 {
						IDFieldRaw = matches[1]
						IDFields = strings.Split(IDFieldRaw, ",")
					}
					re = regexp.MustCompile(`\bAssociationForeignKey:([a-zA-Z0-9_,]+)\b`)
					matches = re.FindStringSubmatch(gormflags)
					if len(matches) == 2 {
						RefFieldRaw = matches[1]
						RefFields = strings.Split(RefFieldRaw, ",")
					}
				}
				if len(IDFields) == 0 { continue }
				if len(RefFields) != 0 {
					WhereMap := make(map[string]interface{})
					for fk := 0; fk < len(RefFields); fk++ {
						WhereMap[RefFields[fk]] = fmt.Sprint(ref.FieldByName(IDFields[fk]))
					}
					db.Where(WhereMap).Find(ref.Field(i).Addr().Interface())
					if ref.Field(i).Addr().Interface() != nil {
						fetchRec(db, ref.Field(i).Addr().Interface())
					}
				} else {
					panic("AssociationForeignKey empty!")
				}
			}
		}
	}
}

