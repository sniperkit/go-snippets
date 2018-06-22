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
 *     Initial: 2017/04/17        Yusan Kurban
 	   Update:  2017/04/26        Jia Chenhui
*/

package models

import (
	"github.com/google/go-github/github"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// 对外服务接口
type GitReposServiceProvider struct {
}

var(
	GitReposService *GitReposServiceProvider
	GitReposCollection *mgo.Collection
)

// 连接、设置索引
func PrepareGitRepos(colName string) {
	GitReposCollection = GithubSession.DB("github").C(colName)

	fullNameIndex := mgo.Index{
		Key:        []string{"FullName"},
		Unique:     true,
		Background: true,
		Sparse:     true,
	}
	if err := GitReposCollection.EnsureIndex(fullNameIndex); err != nil {
		panic(err)
	}

	starCountIndex := mgo.Index{
		Key:        []string{"StarCount"},
		Background: true,
		Sparse:     true,
	}
	if err := GitReposCollection.EnsureIndex(starCountIndex); err != nil {
		panic(err)
	}

	forkCountIndex := mgo.Index{
		Key:        []string{"ForkCount"},
		Background: true,
		Sparse:     true,
	}
	if err := GitReposCollection.EnsureIndex(forkCountIndex); err != nil {
		panic(err)
	}

	langCountIndex := mgo.Index{
		Key:        []string{"Language"},
		Background: true,
		Sparse:     true,
	}
	if err := GitReposCollection.EnsureIndex(langCountIndex); err != nil {
		panic(err)
	}

	GitReposService = &GitReposServiceProvider{}
}

type MDRepos struct {
	ID              bson.ObjectId     `bson:"_id,omitempty",json:"id"`
	RepoID          *int              `bson:"RepoID,omitempty" json:"repoid,omitempty"`
	Owner           *string           `bson:"Owner,omitempty" json:"-"`
	Name            *string           `bson:"Name,omitempty" json:"name"`
	FullName        *string           `bson:"FullName,omitempty" json:"fullname"`
	Description     *string           `bson:"Description,omitempty" json:"description"`
	DefaultBranch   *string           `bson:"DefaultBranch,omitempty" json:"defaultbranch"`
	Language        *string           `bson:"Language,omitempty" json:"language"`
	Created         *github.Timestamp `bson:"Created,omitempty" json:"created"`
	Updated         *github.Timestamp `bson:"Updated,omitempty" json:"updated"`
	Pushed          *github.Timestamp `bson:"Pushed,omitempty" json:"pushed"`
	HasWiki         *bool             `bson:"HasWiki,omitempty" json:"haswiki"`
	HasIssues       *bool             `bson:"HasIssues,omitempty" json:"hasissues"`
	HasDownloads    *bool             `bson:"HasDownloads,omitempty" json:"hasdownloads"`
	ForkCount       *int              `bson:"ForkCount,omitempty" json:"forkcount"`
	StarCount       *int              `bson:"StarCount,omitempty" json:"starcount"`
	WatchersCounts  *int              `bson:"WatchersCounts,omitempty" json:"watcherscounts"`
	OpenIssuesCount *int              `bson:"OpenIssuesCount,omitempty" json:"openissuescount"`
	Size            *int              `bson:"Size,omitempty" json:"size"`
}

// Create 存储库信息及作者在 MDUser 数据库中的 ID
func (rsp *GitReposServiceProvider) Create(repos *github.Repository, owner *string) error {
	r := MDRepos{
		RepoID:          repos.ID,
		Owner:           owner,
		Name:            repos.Name,
		FullName:        repos.FullName,
		Description:     repos.Description,
		DefaultBranch:   repos.DefaultBranch,
		Language:        repos.Language,
		Created:         repos.CreatedAt,
		Updated:         repos.UpdatedAt,
		Pushed:          repos.PushedAt,
		HasWiki:         repos.HasWiki,
		HasIssues:       repos.HasIssues,
		HasDownloads:    repos.HasDownloads,
		ForkCount:       repos.ForksCount,
		StarCount:       repos.StargazersCount,
		WatchersCounts:  repos.WatchersCount,
		OpenIssuesCount: repos.OpenIssuesCount,
		Size:            repos.Size,
	}

	_, err := GitReposCollection.Upsert(bson.M{"RepoID": repos.ID}, &r)
	if err != nil {
		return err
	}

	return nil
}
