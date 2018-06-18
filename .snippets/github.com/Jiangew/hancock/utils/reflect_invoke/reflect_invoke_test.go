package reflect_invoke_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"github.com/jiangew/hancock/utils/reflect_invoke"
)

type Foo struct {
}

func (f *Foo) FooFuncZero() bool {
	return true
}

func (f *Foo) FooFuncOne(arg int) string {
	return strconv.Itoa(arg)
}

func (f *Foo) FooFuncTwo(argStr string, argInt int) string {
	return argStr + strconv.Itoa(argInt)
}

type Bar struct {
}

func (b *Bar) BarFuncZero() string {
	return fmt.Sprintln("BarFuncZero")
}

func (b *Bar) BarFuncOne(arg float64) int {
	return int(arg)
}

func (b *Bar) BarFuncTwo(argStr bool, argInt int) int {
	if argStr {
		return argInt
	} else {
		return -argInt
	}
}

func (b *Bar) BarFuncAdd(argOne, argTwo float64) float64 {
	return argOne + argTwo
}

func (b *Bar) BarSwap(argOne, argTwo float64) (float64, float64) {
	return argTwo, argOne
}

func init() {
	foo := &Foo{}
	bar := &Bar{}

	reflect_invoke.RegisterMethod(foo)
	reflect_invoke.RegisterMethod(bar)
}

func TestInvoke(t *testing.T) {
	resultFooFuncZero := reflect_invoke.InvokeByReflectArgs("FooFuncZero", nil)
	if false == resultFooFuncZero[0].Bool() {
		t.Errorf("invoke FooFuncZero error")
	}

	resultFooFuncOne := reflect_invoke.InvokeByReflectArgs("FooFuncOne",
		[]reflect.Value{reflect.ValueOf(123)})

	if "123" != resultFooFuncOne[0].String() {
		t.Errorf("invoke FooFuncOne error")
	}

	argForFooFuncTwo := []reflect.Value{reflect.ValueOf("str123"), reflect.ValueOf(456)}
	resultFooFuncTwo := reflect_invoke.InvokeByReflectArgs("FooFuncTwo", argForFooFuncTwo)

	if "str123456" != resultFooFuncTwo[0].String() {
		t.Errorf("invoke FooFuncTwo error")
	}

	resultBarFuncZero := reflect_invoke.InvokeByReflectArgs("BarFuncZero", nil)
	if "BarFuncZero" == resultBarFuncZero[0].String() {
		t.Errorf("invoke BarFuncZero error")
	}

	resultBarFuncOne := reflect_invoke.InvokeByReflectArgs("BarFuncOne",
		[]reflect.Value{reflect.ValueOf(123.0)})

	if 123 != resultBarFuncOne[0].Int() {
		t.Errorf("invoke BarFuncOne error")
	}

	argForBarFuncTwo := []reflect.Value{reflect.ValueOf(false), reflect.ValueOf(456)}
	resultBarFuncTwo := reflect_invoke.InvokeByReflectArgs("BarFuncTwo", argForBarFuncTwo)

	if -456 != resultBarFuncTwo[0].Int() {
		t.Errorf("invoke BarFuncTwo error")
	}

}

func TestInvokeByInterfaceArgs(t *testing.T) {
	resultFooFuncZero := reflect_invoke.InvokeByInterfaceArgs("FooFuncZero", nil)
	if false == resultFooFuncZero[0].Bool() {
		t.Errorf("invoke FooFuncZero error")
	}

	resultFooFuncOne := reflect_invoke.InvokeByInterfaceArgs("FooFuncOne", []interface{}{123})

	if "123" != resultFooFuncOne[0].String() {
		t.Errorf("invoke FooFuncOne error")
	}

	resultFooFuncTwo := reflect_invoke.InvokeByInterfaceArgs("FooFuncTwo",
		[]interface{}{"str123", 456})

	if "str123456" != resultFooFuncTwo[0].String() {
		t.Errorf("invoke FooFuncTwo error")
	}

	resultBarFuncZero := reflect_invoke.InvokeByInterfaceArgs("BarFuncZero", nil)
	if "BarFuncZero" == resultBarFuncZero[0].String() {
		t.Errorf("invoke BarFuncZero error")
	}

	resultBarFuncOne := reflect_invoke.InvokeByInterfaceArgs("BarFuncOne", []interface{}{123.1})
	if 123 != resultBarFuncOne[0].Int() {
		t.Errorf("invoke BarFuncOne error")
	}

	resultBarFuncTwo := reflect_invoke.InvokeByInterfaceArgs("BarFuncTwo",
		[]interface{}{false, 456})

	if -456 != resultBarFuncTwo[0].Int() {
		t.Errorf("invoke BarFuncTwo error")
	}
}

func testInvokeByJson(jsonStr, funcName string, expectResult ...interface{}) error {
	result := reflect_invoke.Response{}
	err := json.Unmarshal(reflect_invoke.InvokeByJson([]byte(jsonStr)), &result)
	if err != nil {
		return err
	}

	if result.ErrorCode > 0 {
		return fmt.Errorf("invoke "+funcName+" error: %d", result.ErrorCode)
	}

	if len(result.Data) != len(expectResult) {
		return errors.New("unexpected result")
	}

	for i, resultData := range result.Data {
		var resultDataConvert interface{}
		switch resultData.(type) {
		case float64:
			switch expectResult[i].(type) {
			case float64:
				resultDataConvert = resultData
			default:
				resultDataConvert = int(resultData.(float64))
			}
		default:
			resultDataConvert = resultData
		}
		if resultDataConvert != expectResult[i] {
			return errors.New("invoke " + funcName + " error: result not equal")
		}
	}

	return nil
}

func TestInvokeByJson(t *testing.T) {
	var err error

	jsonDataFooFuncZero := `
				{
					"func_name":"FooFuncZero",
					"params":[]
				}
				`
	err = testInvokeByJson(jsonDataFooFuncZero, "FooFuncZero", true)
	if err != nil {
		t.Error(err)
	}

	jsonDataFooFuncOne := `
				{
					"func_name":"FooFuncOne",
					"params":[456]
				}
				`
	err = testInvokeByJson(jsonDataFooFuncOne, "FooFuncOne", "456")
	if err != nil {
		t.Error(err)
	}

	jsonDataFooFuncTwo := `
				{
					"func_name":"FooFuncTwo",
					"params":["str123",456]
				}
				`
	err = testInvokeByJson(jsonDataFooFuncTwo, "FooFuncTwo", "str123456")
	if err != nil {
		t.Error(err)
	}

	jsonDataBarFuncOne := `
				{
					"func_name":"BarFuncOne",
					"params":[456.0]
				}
				`
	err = testInvokeByJson(jsonDataBarFuncOne, "BarFuncOne", 456)
	if err != nil {
		t.Error(err)
	}

	jsonDataBarFuncTwo := `
				{
					"func_name":"BarFuncTwo",
					"params":[false,456]
				}
				`
	err = testInvokeByJson(jsonDataBarFuncTwo, "BarFuncTwo", -456)
	if err != nil {
		t.Error(err)
	}

	jsonDataBarFuncTwo = `
				{
					"func_name":"BarFuncTwo",
					"params":[false,"456"]
				}
				`
	err = testInvokeByJson(jsonDataBarFuncTwo, "BarFuncTwo", -456)
	if err != nil {
		t.Error(err)
	}

	jsonDataBarFuncAdd := `
				{
					"func_name":"BarFuncAdd",
					"params":[0.5,0.51]
				}
				`
	//这里float64直接比较相等正常？
	err = testInvokeByJson(jsonDataBarFuncAdd, "BarFuncAdd", 1.01)
	if err != nil {
		t.Error(err)
	}

	jsonDataBarFuncSwap := `
				{
					"func_name":"BarSwap",
					"params":[0.1,0.9]
				}
				`
	err = testInvokeByJson(jsonDataBarFuncSwap, "BarSwap", 0.1, 0.9)
	if err != nil {
		t.Error(err)
	}

	err = testInvokeByJson(jsonDataBarFuncSwap, "BarSwap", 0.9, 0.1)
	if err != nil {
		t.Error(err)
	}
}
