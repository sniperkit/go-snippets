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
 *     Initial: 05/05/2017        Jia Chenhui
 */

package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type LabelServiceProvider struct {
}

var (
	GithubSession   *mgo.Session
	LabelService    *LabelServiceProvider
	LabelCollection *mgo.Collection
)

// PrepareGitLabel 连接数据库、创建索引
func PrepareGitLabel() {
	LabelCollection = GithubSession.DB("github").C("Label")

	LabelIndex := mgo.Index{
		Key:        []string{"Name"},
		Unique:     true,
		Background: true,
		Sparse:     true,
	}
	if err := LabelCollection.EnsureIndex(LabelIndex); err != nil {
		panic(err)
	}

	LabelService = &LabelServiceProvider{}
}

// MDLabel 标签数据结构
type MDLabel struct {
	LabelID bson.ObjectId `bson:"LabelID,omitempty" json:"id"`
	Name    string        `bson:"Name" json:"name"`
	Desc    string        `bson:"Desc" json:"desc"`
	Active  bool          `bson:"Active" json:"active"`
	Total   uint64        `bson:"Total" json:"total"`
}

// Activate 修改标签状态数据结构
type Activate struct {
	Name   string
	Active bool
}

// MDModifyLabel 修改标签内容数据结构
type MDModifyLabel struct {
	LabelID string
	Name    string
	Desc    string
	Active  bool
}

// Create 新建标签
func (tsp *LabelServiceProvider) Create(Label *MDLabel) error {
	l := MDLabel{
		LabelID: bson.NewObjectId(),
		Name:    Label.Name,
		Active:  Label.Active,
		Desc:    Label.Desc,
	}

	err := LabelCollection.Insert(&l)
	if err != nil {
		return err
	}

	return nil
}

// ListAll 获取所有标签
func (tsp *LabelServiceProvider) ListAll() ([]MDLabel, error) {
	var l []MDLabel

	err := LabelCollection.Find(nil).All(&l)
	if err != nil {
		return nil, err
	}

	return l, nil
}

// GetLabelInfo 获取单个标签内容
func (tsp *LabelServiceProvider) GetLabelInfo(labelID string) (MDLabel, error) {
	var l MDLabel

	err := LabelCollection.Find(bson.M{"LabelID": bson.ObjectIdHex(labelID)}).One(&l)

	return l, err
}

// Activate 修改标签状态
func (tsp *LabelServiceProvider) Activate(activate Activate) error {
	update := bson.M{"$set": bson.M{
		"Active": activate.Active,
	}}

	err := LabelCollection.Update(bson.M{"Name": activate.Name}, &update)
	if err != nil {
		return err
	}

	return nil
}

// Modify 修改标签内容
func (tsp *LabelServiceProvider) Modify(label MDModifyLabel) error {
	update := bson.M{"$set": bson.M{
		"Name":   label.Name,
		"Desc":   label.Desc,
		"Active": label.Active,
	}}

	err := LabelCollection.Update(bson.M{"LabelID": bson.ObjectIdHex(label.LabelID)}, &update)
	if err != nil {
		return err
	}

	return nil
}
