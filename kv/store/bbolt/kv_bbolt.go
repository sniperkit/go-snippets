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

package store

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/coreos/bbolt"

	"github.com/fengyfei/nuts/kv"
)

type dbBaseBbolt struct {
	db *bolt.DB
}

func init() {
	kv.Register("bbolt", NewBbolt)
}

func NewBbolt(name string) kv.Store {
	db := new(dbBaseBbolt)
	db.init(name)
	return db
}

func (bb *dbBaseBbolt) DB(name string) {
	bb.init(name)
}

func (bb *dbBaseBbolt) init(path string) error {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		fmt.Printf("[open] init db error: %v", err)
		return err
	}
	bb.db = db

	tx, err := bb.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	return tx.Commit()
}

func (bb *dbBaseBbolt) Put(bucket string, key interface{}, value interface{}) error {
	return bb.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		if err := b.Put(convertToByte(key), convertToByte(value)); err != nil {
			return err
		}

		return nil
	})
}

func (bb *dbBaseBbolt) Get(bucket string, key interface{}) ([]byte, error) {
	tx, err := bb.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	b := tx.Bucket([]byte(bucket))
	v := b.Get(convertToByte(key))
	if v == nil {
		err := errors.New("get no record")
		return nil, err
	}

	return v, nil
}

func (bb *dbBaseBbolt) Delete(bucket string, key interface{}) error {
	return bb.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		if err := b.Delete(convertToByte(key)); err != nil {
			return err
		}

		return nil
	})
}

func convertToByte(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}
