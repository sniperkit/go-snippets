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
 *     Initial: 28/04/2017        Jia Chenhui
 */

package github

import (
	"context"
	"math"
	"time"

	"github.com/google/go-github/github"

	"github.com/TechCatsLab/Andariel/pkg/constants"
	"github.com/TechCatsLab/Andariel/pkg/utility"
	gitClient "github.com/fengyfei/nuts/github/client"
)

// SearchRepos 按条件从 github 搜索库，受 github API 限制，一次请求只能获取 1000 条记录
// GitHub API docs: https://developer.github.com/v3/search/#search-repositories
func searchRepos(client *gitClient.GHClient, query string, opt *github.SearchOptions) ([]github.Repository, *github.Response, string, error) {
	var (
		result []github.Repository
		repos  *github.RepositoriesSearchResult
		resp   *github.Response
		stopAt string
		err    error
	)

	page := 1
	maxPage := math.MaxInt32

	for page <= maxPage {
		opt.Page = page

		repos, resp, err = client.Client.Search.Repositories(context.Background(), query, opt)
		if err != nil {
			goto finish
		}

		maxPage = resp.LastPage
		result = append(result, repos.Repositories...)

		page++
	}

finish:
	stopAt = utility.SplitQuery(query)

	return result, resp, stopAt, err
}

// SearchReposByCreated 按创建时间及其它指定条件搜索库
// queries: 指定库的创建时间
// For example:
//     queries := []string{"\"2008-06-01 .. 2012-09-01\"", "\"2012-09-02 .. 2013-03-01\"", "\"2013-03-02 .. 2013-09-03\"", "\"2013-09-04 .. 2014-03-05\"", "\"2014-03-06 .. 2014-09-07\"", "\"2014-09-08 .. 2015-03-09\"", "\"2015-03-10 .. 2015-09-11\"", "\"2015-09-12 .. 2016-03-13\"", "\"2016-03-14 .. 2016-09-15\"", "\"2016-09-16 .. 2017-03-17\""}
//
// querySeg: 指定除创建时间之外的其它条件
// For example:
//     queryPart := constants.QueryLanguage + ":" + constants.LangLua + " " + constants.QueryCreated + ":"
//
// opt: 为搜索方法指定可选参数
// For example:
//     opt := &github.SearchOptions{
//         Sort:        constants.SortByStars,
//         Order:       constants.OrderByDesc,
//         ListOptions: github.ListOptions{PerPage: 100},
//     }
// GitHub API docs: https://developer.github.com/v3/search/#search-repositories
func SearchReposByCreated(client *gitClient.GHClient, queries []string, querySeg string, opt *github.SearchOptions) ([]github.Repository, *github.Response, string, error) {
	var (
		result, repos []github.Repository
		resp          *github.Response
		stopAt        string
		err           error
	)

	for _, q := range queries {
		query := querySeg + q

		repos, resp, stopAt, err = searchRepos(client, query, opt)
		if err != nil {
			goto finish
		}

		result = append(result, repos...)
	}

finish:
	return result, resp, stopAt, err
}

// SearchReposByStartTime 按指定创建时间、时间间隔及其它条件搜索库
// year、month、day: 从此创建时间开始搜索
// For example：
//     year = 2016 month = time.January day = 1
//     时间格式化只能使用 "2006-01-02 15:04:05" 进行，可将年月日和 时分秒拆开使用
//
// incremental: 以此时间增量搜索，如第一次搜索 1 月份的库，第二次搜索 2 月份的库
// For example:
//     interval = "month"
//
// querySeg: 指定除创建时间之外的其它条件
// For example:
//     queryPart := constants.QueryLanguage + ":" + constants.LangLua + " " + constants.QueryCreated + ":"
//
// opt: 为搜索方法指定可选参数
// For example:
//     opt := &github.SearchOptions{
//         Sort:        constants.SortByStars,
//         Order:       constants.OrderByDesc,
//         ListOptions: github.ListOptions{PerPage: 100},
//     }
// GitHub API docs: https://developer.github.com/v3/search/#search-repositories
func SearchReposByStartTime(client *gitClient.GHClient, year int, month time.Month, day int, incremental, querySeg string, opt *github.SearchOptions) ([]github.Repository, *github.Response, string, error) {
	var (
		result, repos []github.Repository
		resp          *github.Response
		stopAt        string
		err           error
	)

	date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	for date.Unix() < time.Now().Unix() {
		var dateFormat string

		switch incremental {
		case constants.Quarter:
			dateFormat = date.Format("2006-01-02") + " .. " + date.AddDate(0, 3, 0).Format("2006-01-02")
		case constants.Month:
			dateFormat = date.Format("2006-01-02") + " .. " + date.AddDate(0, 1, 0).Format("2006-01-02")
		case constants.Week:
			dateFormat = date.Format("2006-01-02") + " .. " + date.AddDate(0, 0, 6).Format("2006-01-02")
		case constants.Day:
			dateFormat = date.Format("2006-01-02") + " .. " + date.AddDate(0, 0, 0).Format("2006-01-02")
		default:
			dateFormat = date.Format("2006-01-02") + " .. " + date.AddDate(0, 1, 0).Format("2006-01-02")
		}

		query := querySeg + "\"" + dateFormat + "\""

		repos, resp, stopAt, err = searchRepos(client, query, opt)
		if err != nil {
			goto finish
		}

		result = append(result, repos...)

		// 防止触发 GitHub 的滥用检测机制，等待一秒
		time.Sleep(1 * time.Second)

		switch incremental {
		case constants.Quarter:
			date = date.AddDate(0, 3, 1)
		case constants.Month:
			date = date.AddDate(0, 1, 1)
		case constants.Week:
			date = date.AddDate(0, 0, 7)
		case constants.Day:
			date = date.AddDate(0, 0, 1)
		default:
			date = date.AddDate(0, 1, 1)
		}
	}

finish:
	return result, resp, stopAt, err
}
