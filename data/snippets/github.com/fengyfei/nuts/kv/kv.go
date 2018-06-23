/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd..
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2017/09/16        Yang Chenglong
 */

package kv

import (
	"sync"
)

const (
	DRBbolt DriverType = 0
)

type (
	DriverType int
	InitFun func(name string) Store

	Store interface {
		DB(name string)
		Put(bucket string, key interface{}, value interface{}) error //set session value
		Get(bucket string, key interface{}) ([]byte, error)          //get session value
		Delete(bucket string, key interface{}) error                 //delete session value
	}
)

var (
	DBStore Store
	drivers = map[string]DriverType{
		"bbolt": DRBbolt,
	}

	dbManager = sync.Map{}
)

func Register(driver string, init InitFun) {
	if init == nil {
		panic("store: Register init is nil")
	}

	d, ok := drivers[driver]
	if !ok {
		panic("store: Unsupported driver" + driver)
	}

	if _, ok := dbManager.Load(d); ok {
		panic("store: Register called twice for init " + driver)
	}

	dbManager.Store(d, init)

	DBStore = init(driver+".db")
}
