package funcs

import (
	"errors"
	"log"
	"reflect"
	"runtime"
	"sync"
)

var (
	ErrorFuncDirectorRegisteTypeNotFunc error = errors.New("functions director: registe item must be a function")
	ErrorFuncDirectorAlreadyRegisted    error = errors.New("functions director: function already registed")
	ErrorFuncDirectorNotRegisted        error = errors.New("functions director: function not registed")
	ErrorFuncDirectorParamsNotMatch     error = errors.New("functions director: function call params not match")
)

type FuncsDirector interface {
	// Register registe function into direcotr
	Registe(interface{}) error
	// Unregister unregiste function from director
	Unregiste(interface{}) error
	// Call call the function by name and params
	Call(funcName string, params ...interface{}) (result []reflect.Value, err error)
}

type funcsDirector struct {
	mu sync.Mutex

	funcs map[string]interface{}
}

func NewFuncsDirector() FuncsDirector {
	return &funcsDirector{
		funcs: make(map[string]interface{}),
	}
}

func (f *funcsDirector) Registe(fun interface{}) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if !isFunction(fun) {
		return ErrorFuncDirectorRegisteTypeNotFunc
	}
	name := getFuncNameByPointer(fun)
	if _, ok := f.funcs[name]; ok {
		return ErrorFuncDirectorAlreadyRegisted
	}
	f.funcs[name] = fun
	return nil
}

func (f *funcsDirector) Unregiste(fun interface{}) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if !isFunction(fun) {
		return ErrorFuncDirectorRegisteTypeNotFunc
	}
	name := getFuncNameByPointer(fun)
	if _, ok := f.funcs[name]; !ok {
		return ErrorFuncDirectorNotRegisted
	}
	delete(f.funcs, name)
	return nil
}

func (f *funcsDirector) Call(funname string, params ...interface{}) ([]reflect.Value, error) {
	f.mu.Lock()
	var fun interface{}
	var ok bool
	if fun, ok = f.funcs[funname]; !ok {
		log.Printf("function director: function %s not registed. functions: %#v", funname, f.funcs)
		return []reflect.Value{}, ErrorFuncDirectorNotRegisted
	}
	f.mu.Unlock()

	fu := reflect.ValueOf(fun)
	if len(params) != fu.Type().NumIn() {
		return []reflect.Value{}, ErrorFuncDirectorParamsNotMatch
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result := fu.Call(in)
	return result, nil
}

// utils function
func getFuncNameByPointer(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

func isFunction(fn interface{}) bool {
	typ := reflect.TypeOf(fn)
	if typ.Kind() != reflect.Func {
		return false
	}
	return true
}
