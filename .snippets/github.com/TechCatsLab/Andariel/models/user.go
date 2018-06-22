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
 *     Initial: 17/04/2017        Jia Chenhui
 */

package models

import (
	"github.com/google/go-github/github"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// 对外服务接口
type GitUserServiceProvider struct {
}

var (
	GitUserService *GitUserServiceProvider
	GitUserCollection *mgo.Collection
)

// 连接、设置索引
func PrepareGitUser() {
	GitUserCollection = GithubSession.DB("github").C("Owner")
	userIndex := mgo.Index{
		Key:        []string{"Login"},
		Unique:     true,
		Background: true,
		Sparse:     true,
	}

	if err := GitUserCollection.EnsureIndex(userIndex); err != nil {
		panic(err)
	}

	GitUserService = &GitUserServiceProvider{}
}

// GitHub 用户数据结构
type MDUser struct {
	UserID            bson.ObjectId     `bson:"_id,omitempty" json:"id"`
	Login             *string           `bson:"Login,omitempty" json:"login"`
	ID                *int              `bson:"ID,omitempty" json:"userid"`
	HTMLURL           *string           `bson:"HTMLURL,omitempty" json:"htmlurl"`
	Name              *string           `bson:"Name,omitempty" json:"name"`
	Email             *string           `bson:"Email,omitempty" json:"email"`
	PublicRepos       *int              `bson:"PublicRepos,omitempty" json:"publicrepos"`
	PublicGists       *int              `bson:"PublicGists,omitempty" json:"publicgists"`
	Followers         *int              `bson:"Followers,omitempty" json:"followers"`
	Following         *int              `bson:"Following,omitempty" json:"following"`
	CreatedAt         *github.Timestamp `bson:"CreatedAt,omitempty" json:"created"`
	UpdatedAt         *github.Timestamp `bson:"UpdatedAt,omitempty" json:"updated"`
	SuspendedAt       *github.Timestamp `bson:"SuspendedAt,omitempty" json:"suspended"`
	Type              *string           `bson:"Type,omitempty" json:"type"`
	TotalPrivateRepos *int              `bson:"TotalPrivateRepos,omitempty" json:"totalprivaterepos"`
	OwnedPrivateRepos *int              `bson:"OwnedPrivateRepos,omitempty" json:"ownedprivaterepos"`
	PrivateGists      *int              `bson:"PrivateGists,omitempty" json:"privategists"`
}

// GetUserByID 查询作者信息
func (usp *GitUserServiceProvider) GetUserByID(userID *int) (*MDUser, error) {
	var u MDUser

	err := GitUserCollection.Find(bson.M{"ID": userID}).One(&u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// GetUserID 通过 login 获取作者在数据库中的 _id
func (usp *GitUserServiceProvider) GetUserID(login *string) (string, error) {
	var u MDUser

	err := GitUserCollection.Find(bson.M{"Login": login}).One(&u)
	if err != nil {
		return "", err
	}

	return u.UserID.Hex(), nil
}

// Create 存储作者信息，先查询数据库中是否有此用户信息, 若没有则创建，有则更新
func (usp *GitUserServiceProvider) Create(user *MDUser) (string, error) {
	userID, err := usp.GetUserID(user.Login)
	if err != nil {
		if err != mgo.ErrNotFound {
			return "", err
		}

		goto create
	} else {
		update := MDUser{
			Login:             user.Login,
			ID:                user.ID,
			HTMLURL:           user.HTMLURL,
			Name:              user.Name,
			Email:             user.Email,
			PublicRepos:       user.PublicRepos,
			PublicGists:       user.PublicGists,
			Followers:         user.Followers,
			Following:         user.Following,
			CreatedAt:         user.CreatedAt,
			UpdatedAt:         user.UpdatedAt,
			SuspendedAt:       user.SuspendedAt,
			Type:              user.Type,
			TotalPrivateRepos: user.TotalPrivateRepos,
			OwnedPrivateRepos: user.OwnedPrivateRepos,
			PrivateGists:      user.PrivateGists,
		}

		err = GitUserCollection.Update(bson.M{"_id": bson.ObjectIdHex(userID)}, &update)
		if err != nil {
			return "", err
		}

		return userID, nil
	}

create:
	create := MDUser{
		UserID:            bson.NewObjectId(),
		Login:             user.Login,
		ID:                user.ID,
		HTMLURL:           user.HTMLURL,
		Name:              user.Name,
		Email:             user.Email,
		PublicRepos:       user.PublicRepos,
		PublicGists:       user.PublicGists,
		Followers:         user.Followers,
		Following:         user.Following,
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
		SuspendedAt:       user.SuspendedAt,
		Type:              user.Type,
		TotalPrivateRepos: user.TotalPrivateRepos,
		OwnedPrivateRepos: user.OwnedPrivateRepos,
		PrivateGists:      user.PrivateGists,
	}

	err = GitUserCollection.Insert(&create)
	if err != nil {
		return "", err
	}

	return create.UserID.Hex(), nil
}
