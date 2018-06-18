package toolkittest

import (
	"reflect"
	"testing"

	"github.com/eaciit/toolkit"

	. "github.com/eaciit/toolkit"
)

type Data struct {
	ID   string
	Name string
	Age  int
}

func TestSerdeObj(t *testing.T) {
	var d1, d2 Data

	d1 = Data{"d01", "name", 10}
	err := Serde(d1, &d2, "")
	if err != nil {
		check(t, err, "failed to serde")
	}

	if !reflect.DeepEqual(d1, d2) {
		check(t, toolkit.Errorf("object is not same: %v, %v", d1, d2), "")
	}
}

func TestSerdePtr(t *testing.T) {
	var d1, d2 *Data

	d1 = &Data{"d01", "name", 10}
	d2 = new(Data)
	err := Serde(d1, d2, "")
	if err != nil {
		check(t, err, "failed to serde")
	}

	if !reflect.DeepEqual(d1, d2) {
		check(t, toolkit.Errorf("object is not same: %v, %v", d1, d2), "")
	}
}

func TestSerdeMtoM(t *testing.T) {
	var m1, m2 M
	m1 = toolkit.M{}.Set("id", 1000).Set("name", "name 2")
	err := Serde(m1, &m2, "")
	if err != nil {
		check(t, err, "failed to serde")
	}

	for k, v := range m1 {
		v2 := m2.Get(k)
		if v != v2 {
			check(t, toolkit.Errorf("object is not same for %s: %v and %v", k, v, v2), "")
		}
	}
}

func TestSerdeMtoMPtr(t *testing.T) {
	var m1, m2 M
	m1 = toolkit.M{}.Set("id", "1000").Set("name", "name 2")
	err := Serde(&m1, &m2, "")
	if err != nil {
		check(t, err, "failed to serde")
	}

	for k, v := range m1 {
		v2 := m2.Get(k)
		if v != v2 {
			check(t, toolkit.Errorf("object is not same for %s: %v and %v", k, v, v2), "")
		}
	}
}

func TestMsPtrToSlice(t *testing.T) {
	d1 := []toolkit.M{toolkit.M{}.Set("ID", "1000").Set("Name", "Name1"),
		toolkit.M{}.Set("ID", "1001").Set("Name", "Name2")}
	var d2 []Data

	err := Serde(&d1, &d2, "")
	if err != nil {
		check(t, err, "failed to serde")
	}

	if d1[1].Get("Name") != d2[1].Name {
		check(t, toolkit.Errorf("object is not same: %v\n%v", d1, d2), "")
	} else {
		toolkit.Printfn("d1=d2\n%v\n%v", d1, d2)
	}
}

func check(t *testing.T, err error, pre string) {
	if err != nil {
		if pre == "" {
			t.Fatalf(err.Error())
		} else {
			t.Fatalf(Sprintf("%s: %s", pre, err.Error()))
		}
	}
}
