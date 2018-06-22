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
 *     Initial: 2017/06/18        Liu Jiachang
 */

package task

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MgoQueueEngine struct {
	Sess	*mgo.Session
}

func (this *MgoQueueEngine) FetchTasks(n int) ([]Task, error) {
	var ta []Task

	c := this.Sess.DB(MDbName).C(MDColl)
	err := c.Find(bson.M{"status": 1}).Limit(n).All(&ta)

	if err != nil {
		return ta, err
	}

	return ta, nil
}

func (this *MgoQueueEngine) DeleteTask(id interface{}) error {
	c := this.Sess.DB(MDbName).C(MDColl)

	return c.RemoveId(id)
}

func (this *MgoQueueEngine) Activate(id interface{}, status int16) error {
	c := this.Sess.DB(MDbName).C(MDColl)

	return c.UpdateId(id, bson.M{"$set": bson.M{"status": status}})
}
