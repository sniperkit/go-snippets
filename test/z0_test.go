package toolkittest

import (
	. "github.com/eaciit/toolkit"
	"math"
	"os"
	"testing"
	"time"
)

func killApp(code int) {
	defer os.Exit(code)
}

func TestMtoStruct(t *testing.T) {
	type SubStruct struct {
		Info string
	}

	type SomeStruct struct {
		Title string
		Value int
		IsOK  bool
		Sub   SubStruct
	}

	data := M{
		"Title": "Test",
		"Value": 1232,
		"IsOK":  true,
		"Sub": M{
			"Info": "hahahah",
		},
	}

	// M to Struct
	var result1 = new(SomeStruct)
	err := MtoStruct(data, result1)
	if err != nil {
		t.Fatalf(err.Error())
		return
	}
	t.Logf("M to Struct\n%#v\n", result1)

	// Struct to M
	result2 := &(M{})
	err = StructToM(result1, result2)
	if err != nil {
		t.Fatalf(err.Error())
		return
	}
	t.Logf("Struct to M\n%#v\n", result2)

	// Set property by name
	err = SetPropByName(&result1, "Title", "OK")
	if err != nil {
		t.Fatalf(err.Error())
		return
	}
	t.Logf("Set property by name\n%#v\n", result1)

	// Set property by name which is object
	sub2 := new(SubStruct)
	sub2.Info = "test"
	err = SetPropByName(&result1, "Sub", sub2)
	if err != nil {
		t.Fatalf(err.Error())
		return
	}
	t.Logf("Set property (sub struct) by name\n%#v\n", result1)
}

func TestExecFn(t *testing.T) {
	defer killApp(100)

	type Obj struct {
		Name string
		Len  int
	}

	rvobj, _ := ExecFunc(func(s string) *Obj {
		o := new(Obj)
		o.Name = s
		o.Len = len(s)
		return o
	}, "Arief Darmawan")

	obj := rvobj[0].Interface().(*Obj)
	Println("Got:", JsonString(obj))

	if obj.Name != "Arief Darmawan" {
		t.Fatalf("error")
	}
}

func TestFormulaSimple(t *testing.T) {
	//defer killApp(50)
	yr := float64(2016)
	xwant := float64((300.0 + yr) * 2.0)
	f, e := NewFormula("(300+@yr)*2")
	if e != nil {
		t.Fatalf("Error build formula. %s", e.Error())
	}
	xgot := f.Run(M{}.Set("@yr", yr))
	if xgot != xwant {
		t.Fatalf("Want %.2f got %.2f", xwant, xgot)
	} else {
		t.Logf("Want %.2f got %.2f", xwant, xgot)
	}
}

func TestFormulaComplex(t *testing.T) {
	//t.Skip()
	//defer killApp(50)
	yr := float64(2016)
	xwant := float64(-2.00 + 3*(200.0+300.00)/500.0 + yr)
	f, e := NewFormula("-2+3*(200+300.00)/500+@yr")
	if e != nil {
		t.Fatalf("Error build formula. %s", e.Error())
	}
	xgot := f.Run(M{}.Set("@yr", yr))
	if xgot != xwant {
		t.Fatalf("Want %.2f got %.2f", xwant, xgot)
	} else {
		t.Logf("Want %.2f got %.2f", xwant, xgot)
	}
}

func TestFormulaPower(t *testing.T) {
	//t.Skip()
	//defer killApp(50)
	yr := float64(5)
	xwant := float64(yr + (2.0 * math.Pow(4.0, 3.0) / 2.0))
	f, e := NewFormula("@yr+(2.0*4^3/2)")
	if e != nil {
		t.Fatalf("Error build formula. %s", e.Error())
	}
	xgot := f.Run(M{}.Set("@yr", yr))
	if xgot != xwant {
		t.Fatalf("Want %.2f got %.2f", xwant, xgot)
	} else {
		t.Logf("Want %.2f got %.2f", xwant, xgot)
	}
}

func TestID(t *testing.T) {
	type Data struct {
		ID         string
		Name       string
		Desciption string
	}

	data := &Data{"EACIIT", "EACIIT", "Enhance Enterprise Value"}
	Printf("Data before ID set:\n%s\n\n", JsonString(data))
	f, i := IdInfo(data)
	if f == "" {
		t.Errorf("Unable to find ID")
		return
	}
	Printf("Field: %s\nValue: %v\n\n", f, i)
	e := SetId(data, "EC")
	if e != nil {
		t.Errorf("Unable to set ID: " + e.Error())
		return
	}
	Printf("Data after ID set:\n%s\n\n", JsonString(data))
}

func TestCast(t *testing.T) {
	var (
		f64        float64 = 405.21
		f64other   float64
		int64      int = 201
		int64other int

		date       time.Time
		dateOther  time.Time
		dateString string
	)

	date = time.Date(1980, 4, 1, 0, 0, 0, 0, time.UTC)
	dateOther = String2Date("1-April-1980", "d-MMMM-yyyy")
	if dateOther != date {
		t.Errorf("Strign2Date fail. Expect %v got %v", date, dateOther)
	} else {
		t.Logf("String2Date: %v", dateOther)
	}

	dateString = Date2String(date, "dd MMM yy")
	if dateString != "01 Apr 80" {
		t.Errorf("Date2String fail. Expect %s got %s", "01 Apr 80", dateString)
	} else {
		t.Logf("Date2String: %s", dateString)
	}

	int64other = ToInt(f64, RoundingAuto)
	if int64other != 405 {
		t.Errorf("ToInt fail. Expect %d got %d", 405, int64other)
	}

	f64other = ToFloat64(&int64, 0, RoundingAuto) + 0.758
	if f64other != 201.758 {
		t.Errorf("ToFloat64 fail. Expect %f got %f", 201.758, f64other)
	}
}
