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

	"github.com/google/go-github/github"

	"github.com/TechCatsLab/Andariel/models"
	gitClient "github.com/fengyfei/nuts/github/client"
)

// GetOwnerByID 调用 GitHub API 获取作者信息
func GetOwnerByID(ownerID int, client *gitClient.GHClient) (*models.MDUser, *github.Response, error) {
	owner, resp, err := client.Client.Users.GetByID(context.Background(), ownerID)
	if err != nil {
		if resp != nil {
			return nil, resp, err
		}

		return nil, nil, err
	}

	user := &models.MDUser{
		Login:             owner.Login,
		ID:                owner.ID,
		HTMLURL:           owner.HTMLURL,
		Name:              owner.Name,
		Email:             owner.Email,
		PublicRepos:       owner.PublicRepos,
		PublicGists:       owner.PublicGists,
		Followers:         owner.Followers,
		Following:         owner.Following,
		CreatedAt:         owner.CreatedAt,
		UpdatedAt:         owner.UpdatedAt,
		SuspendedAt:       owner.SuspendedAt,
		Type:              owner.Type,
		TotalPrivateRepos: owner.TotalPrivateRepos,
		OwnedPrivateRepos: owner.OwnedPrivateRepos,
		PrivateGists:      owner.PrivateGists,
	}

	return user, resp, nil
}
