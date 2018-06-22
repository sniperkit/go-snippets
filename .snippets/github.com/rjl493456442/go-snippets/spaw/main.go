// Copyright 2016-2018 Hyperchain Corp.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"math/big"

	"github.com/davecgh/go-spew/spew"
	"fmt"
)

type LocalTyp struct {
	Field string
}

type ComplexTyp struct {
	Field1 string
	Field2 []byte
	Field3 float32
	Field4 *big.Int
	Local  LocalTyp
}

type StringableTyp struct {
	Field  string
	Field2 int
}

func (typ StringableTyp) String() string {
	return fmt.Sprintf("field1: %s, field2: %d", typ.Field, typ.Field2)
}

func main() {
	obj := &ComplexTyp{
		Field1: "field1",
		Field2: []byte("field2"),
		Field3: 1.46,
		Field4: big.NewInt(100),
		Local: LocalTyp{
			Field: "new field",
		},
	}
	spew.Println(obj)
	spew.Dump(obj)

	obj2 := StringableTyp{
		Field:  "hello",
		Field2: 100,
	}
	config := spew.NewDefaultConfig()
	// DisableMethods specifies whether or not error and Stringer interfaces are
	// invoked for types that implement them.
	config.DisableMethods = true
	config.Println(obj2)

	config.Dump(obj2)

	config2 := spew.NewDefaultConfig()
	config2.Println(obj2)

}
