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
	"container/list"
	"fmt"
	"sync"
)

const notExist int = -1

type Elem struct {
	k, v int
}

type LRUCache struct {
	order    *list.List
	data     map[int]*list.Element
	capacity int
	mu       sync.RWMutex
}

func Constructor(capacity int) LRUCache {
	return LRUCache{
		capacity: capacity,
		data:     make(map[int]*list.Element),
		order:    list.New(),
	}
}

func (this *LRUCache) Get(key int) int {
	this.mu.RLock()
	defer this.mu.RUnlock()
	if val, exist := this.data[key]; exist {
		// Update the order list
		this.order.MoveToBack(val)
		return val.Value.(Elem).v
	}
	return notExist
}

func (this *LRUCache) Put(key int, value int) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if val, exist := this.data[key]; exist {
		// Update order only.
		val.Value = Elem{key, value}
		this.order.MoveToBack(val)
		return
	}
	this.data[key] = this.order.PushBack(Elem{key, value})
	if len(this.data) > this.capacity {
		headElem := this.order.Front()
		delete(this.data, headElem.Value.(Elem).k)
		this.order.Remove(headElem)
	}
}

func main() {
	cache := Constructor(2)
	cache.Put(2, 1)
	cache.Put(1, 1)
	cache.Put(2, 3)
	cache.Put(4, 1)
	fmt.Println(cache.Get(1))
	fmt.Println(cache.Get(2))
}
