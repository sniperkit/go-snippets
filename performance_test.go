package skiplist

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
	"unsafe"
)

type List struct {
	list map[string]bool
	mx   *sync.RWMutex
}

var keys [][]byte
var keysString []string
var list *Skiplist
var list1 map[string]bool
var list2 *List

const defaultKeys = 100000

func init() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	//初始化10000000个keys
	keys = make([][]byte, 0, defaultKeys)
	keysString = make([]string, 0, defaultKeys)

	for i := 0; i < defaultKeys; i++ {

		//随机生成key
		var length = r.Intn(1000)
		var buff []byte
		for j := 0; j < length; j++ {
			buff = make([]byte, 0, length)

			buff = append(buff, byte(r.Uint32()))
		}

		keys = append(keys, buff)
		keysString = append(keysString, string(buff))
	}

	list = NewSkiplist(32)
	list1 = make(map[string]bool)
	list2 = &List{
		list: make(map[string]bool),
		mx:   &sync.RWMutex{},
	}

	fmt.Println("init success")
}

//单线程
func BenchmarkSkiplistPut(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var key = keys[i%defaultKeys]
		var value = 111
		list.Put(key, unsafe.Pointer(&value))
	}
}

func BenchmarkMapPut(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var key = keysString[i%defaultKeys]
		list1[key] = true
	}
}

func BenchmarkLockMapPut(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var key = keysString[i%defaultKeys]
		list2.mx.Lock()
		list2.list[key] = true
		list2.mx.Unlock()
	}
}

var threads = 1000

//多线程
func TestSkiplistPut(t *testing.T) {
	var wg = &sync.WaitGroup{}
	wg.Add(threads * 100000)

	var start = time.Now()

	for i := 0; i < threads; i++ {
		go func() {

			for j := 0; j < 100000; j++ {
				var key = keys[j%defaultKeys]
				var value = true
				list.Put(key, unsafe.Pointer(&value))
				wg.Done()
			}

		}()
	}

	wg.Wait()

	t.Logf("time: %f", time.Since(start).Seconds())
}

func TestLockMapPut(t *testing.T) {
	var wg = &sync.WaitGroup{}
	wg.Add(threads * 100000)

	var start = time.Now()

	for i := 0; i < threads; i++ {
		go func() {
			for j := 0; j < 100000; j++ {
				var key = keysString[j%defaultKeys]
				list2.mx.Lock()
				list2.list[key] = true
				list2.mx.Unlock()
				wg.Done()
			}

		}()
	}

	wg.Wait()

	t.Logf("time: %f", time.Since(start).Seconds())
}

func TestSkiplistGet(t *testing.T) {
	var wg = &sync.WaitGroup{}
	wg.Add(threads * 100000)

	var start = time.Now()

	for i := 0; i < threads; i++ {
		go func() {

			for j := 0; j < 100000; j++ {
				var key = keys[j%defaultKeys]
				list.Get(key)
				wg.Done()
			}

		}()
	}

	wg.Wait()

	t.Logf("time: %f", time.Since(start).Seconds())
}

func TestLockMapGet(t *testing.T) {
	var wg = &sync.WaitGroup{}
	wg.Add(threads * 100000)

	var start = time.Now()

	for i := 0; i < threads; i++ {
		go func() {
			for j := 0; j < 100000; j++ {
				var key = keysString[j%defaultKeys]
				list2.mx.RLock()
				_ = list2.list[key]
				list2.mx.RUnlock()
				wg.Done()
			}

		}()
	}

	wg.Wait()

	t.Logf("time: %f", time.Since(start).Seconds())
}
