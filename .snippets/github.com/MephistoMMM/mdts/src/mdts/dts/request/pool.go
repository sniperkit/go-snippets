package request

import (
	"errors"
	"sync"
)

var (
	// ErrPoolNotInit 代表池未初始化
	ErrPoolNotInit = errors.New("Pool Not Init")
	// ErrPoolInvalidSize 代表池初始化时参数错误
	ErrPoolInvalidSize = errors.New("Pool Invalid Size")
	// ErrPoolInvalidGenerateFunc 代表对象生成函数无效
	ErrPoolInvalidGenerateFunc = errors.New("Pool Invalid Generate Function")
)

// Object 池中的元素接口
type Object interface{}

// Pool 池接口
type Pool interface {
	// Get 得到一个池中元素
	Get() (Object, error)
	// Put 放入一个池中元素
	Put(Object) error

	// Init 初始化池
	Init(cap int) error
}

// HoldPool 保存池中元素，无超时机制，且在元素不够时会阻塞 goroutine
type HoldPool struct {
	m        sync.RWMutex
	pool     chan Object
	generate func() (Object, error)
}

// NewHoldPool 创建并初始化 HoldPool
func NewHoldPool(cap int, generate func() (Object, error)) (p Pool, err error) {
	p = &HoldPool{
		generate: generate,
	}

	err = p.Init(cap)
	return
}

// Get ...
func (hp *HoldPool) Get() (obj Object, err error) {
	hp.m.RLock()
	defer hp.m.RUnlock()

	if hp.pool == nil {
		err = ErrPoolNotInit
		return
	}

	obj = <-hp.pool
	return
}

// Put ...
func (hp *HoldPool) Put(obj Object) error {
	hp.m.RLock()
	defer hp.m.RUnlock()

	if hp.pool == nil {
		return ErrPoolNotInit

	}

	if len(hp.pool) == cap(hp.pool) {
		return nil
	}

	hp.pool <- obj
	return nil
}

// Init ...
func (hp *HoldPool) Init(cap int) error {
	hp.m.Lock()
	defer hp.m.Unlock()

	if hp.generate == nil {
		return ErrPoolInvalidGenerateFunc
	}

	if cap <= 0 {
		return ErrPoolInvalidSize
	}

	p := make(chan Object, cap)
	for i := 0; i < cap; i++ {
		obj, err := hp.generate()
		if err != nil {
			return err
		}
		p <- obj
	}

	hp.pool = p
	return nil
}
