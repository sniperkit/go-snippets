package bamboo

import (
	"context"
	"errors"
	"log"
	"reflect"
	"runtime"
	"sync"

	tomb "gopkg.in/tomb.v1"
)

var (
	ErrorBambooAlreadyRunning  error = errors.New("bamboo: already running")
	ErrorBambooAlreadyCanceled error = errors.New("bamboo: already cancel")
	ErrorBambooCancel          error = errors.New("bamboo: cancel")

	ErrorBambooJoinNotFunc                     error = errors.New("bamboo: not function")
	ErrorBambooJoinFuncParamsInNotMatch        error = errors.New("bamboo: params not match")
	ErrorBambooJoinFuncResultEmpty             error = errors.New("bamboo: function result is empty")
	ErrorBambooJoinFuncFirstInParamNotContext  error = errors.New("bamboo: function first param is not context.Context")
	ErrorBambooJoinFuncFirstOutParamNotContext error = errors.New("bamboo: function first result is not context.Context")
	ErrorBambooPieceInvokeFirstOutNotContext   error = errors.New("bamboo: function invoke first result is not context.Context")
	ErrorBambooPieceInvokeCallResultError      error = errors.New("bamboo: function invoke results is error")
)

// context.Context type
var contextType reflect.Type = reflect.TypeOf((*context.Context)(nil)).Elem()

// add functions and consist these into a pipe
type Bamboo interface {
	// Join add a piece of bamboo
	Join(fn interface{}, params ...interface{}) Bamboo
	// Go run these functions
	Go() (context.Context, error)
	// Cancel abort this operation
	Cancel() error
}

type bamboo struct {
	mut  sync.Mutex
	tomb *tomb.Tomb

	trunk []*piece
	err   error
	run   bool
}

func NewBamboo() *bamboo {
	return &bamboo{
		tomb: new(tomb.Tomb),
	}
}

func (b *bamboo) Join(fn interface{}, params ...interface{}) *bamboo {
	if b.err != nil || b.run {
		return b
	}
	select {
	case <-b.tomb.Dying():
		return b
	default:
	}

	fntyp := getInterfaceType(fn)
	if reflect.Func != fntyp.Kind() {
		b.errSet(ErrorBambooJoinNotFunc)
		return b
	}

	// check the first argument is context.Context
	if fntyp.NumIn() > 0 {
		if fntyp.In(0) != contextType {
			b.errSet(ErrorBambooJoinFuncFirstInParamNotContext)
			return b
		}
	} else {
		b.errSet(ErrorBambooJoinFuncParamsInNotMatch)
		return b
	}

	// check function return results
	if fntyp.NumOut() == 0 {
		b.errSet(ErrorBambooJoinFuncResultEmpty)
		return b
	}
	// check the first result is context.Context
	if fntyp.Out(0) != contextType {
		b.errSet(ErrorBambooJoinFuncFirstOutParamNotContext)
		return b
	}

	b.mut.Lock()
	defer b.mut.Unlock()
	// check the previous function return results count
	if len(b.trunk) > 0 {
		var last *piece
		last = b.trunk[len(b.trunk)-1]
		lastType := getInterfaceType(last.fn)
		if fntyp.NumIn() != (len(params) + lastType.NumOut()) {
			b.errSet(ErrorBambooJoinFuncParamsInNotMatch)
			return b
		}
	} else if (len(params) + 1) != fntyp.NumIn() {
		b.errSet(ErrorBambooJoinFuncParamsInNotMatch)
		return b
	}

	p := NewPiece(fn, params...)
	b.trunk = append(b.trunk, p)

	return b
}

func (b *bamboo) errSet(err error) {
	if b.err == nil {
		b.err = err
		return
	}
	return
}

func (b *bamboo) Cancel() error {
	b.tomb.Kill(ErrorBambooCancel)
	return nil
}

func (b *bamboo) Go() error {
	b.mut.Lock()
	err := b.err
	if err != nil {
		b.mut.Unlock()
		return err
	}
	if b.run {
		b.mut.Unlock()
		return ErrorBambooAlreadyRunning
	}
	b.run = true
	b.mut.Unlock()
	go func() {
		defer b.tomb.Done()
		b.tomb.Kill(b.runLoop())
	}()
	err = b.tomb.Wait()
	return err
}

func (b *bamboo) runLoop() error {
	b.mut.Lock()
	defer b.mut.Unlock()

	trunk := b.trunk

	ctx, cancel := context.WithCancel(context.Background())
	var resume chan struct{} = make(chan struct{}, 1)
	var errch chan error = make(chan error, 1)
	var results []reflect.Value
	var index int
	var err error

	resume <- struct{}{}

	for {
		select {
		case <-b.tomb.Dying():
			cancel()
			return nil
		case <-resume:
			if index == len(trunk) {
				log.Println("bamboo: run loop over")
				return nil
			}
			pi := trunk[index]
			go func(p *piece) {
				log.Printf("bamboo: function '%s' prepare to run", p.Name())
				ctx, results, err = p.Invoke(ctx, results)
				if err != nil {
					log.Printf("bamboo: function invoke error: %s", err)
					errch <- err
					return
				}
				resume <- struct{}{}
			}(pi)
			index++
		case erro := <-errch:
			log.Printf("bamboo: fucntion call error: %s", erro)
			return err
		}
	}
}

type piece struct {
	fn     interface{}
	params []reflect.Value
}

func NewPiece(fn interface{}, params ...interface{}) *piece {
	p := &piece{}

	params_ := make([]reflect.Value, len(params)+1)
	for index, param := range params {
		params_[index+1] = reflect.ValueOf(param)
	}
	p.params = params_
	p.fn = fn

	return p
}

func (p *piece) Name() string {
	return getFunctionName(p.fn)
}

func (p *piece) Invoke(ctx context.Context, params []reflect.Value) (context.Context, []reflect.Value, error) {
	// function params prepare
	p.params[0] = reflect.ValueOf(ctx)
	p.params = append(p.params, params...)

	// function call
	fu := getInterfaceValue(p.fn)
	results := fu.Call(p.params)

	// result
	if results[0].IsValid() {
		var ctx_ context.Context
		var ok bool
		if ctx_, ok = results[0].Interface().(context.Context); ok {
			return ctx_, results[1:], nil
		}
		return nil, results[1:], ErrorBambooPieceInvokeFirstOutNotContext
	}
	return nil, results[1:], ErrorBambooPieceInvokeCallResultError
}

func InterfaceOf(value interface{}) reflect.Type {
	typ := reflect.TypeOf(value)

	for reflect.Ptr == typ.Kind() {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Interface {
		panic("bamboo: InterfaceOf value not a interface")
	}

	return typ
}

func getFunctionName(fn interface{}) string {
	val := getInterfaceValue(fn)
	return runtime.FuncForPC(val.Pointer()).Name()
}

func getInterfaceType(value interface{}) reflect.Type {
	typ := reflect.TypeOf(value)
	for reflect.Ptr == typ.Kind() {
		typ = typ.Elem()
	}
	return typ
}

func getInterfaceValue(value interface{}) reflect.Value {
	val := reflect.ValueOf(value)
	for reflect.Ptr == val.Kind() {
		val = val.Elem()
	}
	return val
}
