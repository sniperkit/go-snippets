package otp

import (
	"errors"
	"sync"
)

// once control
var _once sync.Once
var _global globalNames

type globalNames struct {
	mutex     sync.Mutex
	registers map[Name]*gen
}

func init() {
	_once.Do(func() {
		_global = globalNames{
			mutex:     sync.Mutex{},
			registers: make(map[Name]*gen),
		}
	})
}

var (
	ErrorNameAlreadyRegister error = errors.New("register: name already register")
	ErrorNameNotFound        error = errors.New("register: name not found")
)

func registerName(name Name, gen *gen) error {
	_global.mutex.Lock()
	defer _global.mutex.Unlock()
	if _, ok := _global.registers[name]; ok {
		return ErrorNameAlreadyRegister
	}
	// register it
	_global.registers[name] = gen

	return nil
}

func unregisterName(name Name) error {
	_global.mutex.Lock()
	defer _global.mutex.Unlock()
	if _, ok := _global.registers[name]; !ok {
		return ErrorNameNotFound
	}
	// delete it
	delete(_global.registers, name)

	return nil
}

func getGenByName(name Name) (*gen, error) {
	_global.mutex.Lock()
	defer _global.mutex.Unlock()
	if gen, ok := _global.registers[name]; ok {
		return gen, nil
	}
	return nil, ErrorNameNotFound
}
