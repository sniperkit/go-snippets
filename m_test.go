package toolkit

import (
	"encoding/json"
	"fmt"
	"testing"
)

func ExampleM_PathGet() {
	input := []byte(`{
		"student": {
			"name": "John Doe",
			"class": {
				"name": "A5",
				"year": 2006,
				"major": {
					"name": "Biomedic"
				}
			}
		}
	}`)

	var obj M
	err := json.Unmarshal(input, &obj)
	if err != nil {
		fmt.Println("error: ", err)
	}

	majorName, err := obj.PathGet("student.class.major.name")
	if err != nil {
		fmt.Println("error: ", err)
	}

	fmt.Printf("%v", majorName)
	// Output:
	// Biomedic
}

func TestMPathGetInvalidPath(t *testing.T) {
	obj := M{
		"student": M{
			"name": "Charlie",
		},
	}

	var result interface{}
	var err error
	result, err = obj.PathGet("student.class.major.name")
	if err == nil {
		t.Fail()
	}
	if result != nil {
		t.Fail()
	}

	result, err = obj.PathGet("student.name.fullname")
	if err == nil {
		t.Fail()
	}
	if result != nil {
		t.Fail()
	}
}

func TestGetRef(t *testing.T) {
	type model struct {
		Key   string
		Value float32
	}
	var d3 model

	d1 := model{"M1", 20.53}
	m := M{}.Set("D1", &d1)

	m.GetRef("D1", new(model), &d3)
	if d3.Value != d1.Value {
		t.Errorf("Expecting %v, got %v", d1.Value, d3.Value)
	}
}

func TestGetRefNil(t *testing.T) {
	type model struct {
		Key   string
		Value float32
	}
	var d3 model

	d1 := model{"M1", 20.53}
	m := M{}.Set("D1", nil)

	m.GetRef("D1", &d1, &d3)
	if d3.Value != d1.Value {
		t.Errorf("Expecting %v, got %v", d1.Value, d3.Value)
	}
}

func TestGetRefNonPointer(t *testing.T) {
	type model struct {
		Key   string
		Value float32
	}
	var d3 model

	d1 := model{"M1", 20.53}
	m := M{}.Set("D1", nil)

	m.GetRef("D1", &d1, &d3)
	if d3.Value != d1.Value {
		t.Errorf("Expecting %v, got %v", d1.Value, d3.Value)
	}
}

type model struct {
	Key   string
	Value float32
}

var result model

func BenchmarkGetRef(b *testing.B) {
	var d3 model

	d1 := model{"M1", 20.53}
	m := M{}.Set("D1", d1)

	for n := 0; n < b.N; n++ {
		m.GetRef("D1", &d1, &d3)
	}
}

func BenchmarkGet(b *testing.B) {
	d1 := model{"M1", 20.53}
	m := M{}.Set("D1", d1)

	for n := 0; n < b.N; n++ {
		_ = m.Get("D1", d1)
	}
}

func BenchmarkGetCasted(b *testing.B) {

	d1 := model{"M1", 20.53}
	m := M{}.Set("D1", d1)

	for n := 0; n < b.N; n++ {
		result = m.Get("D1", d1).(model)
	}
}
