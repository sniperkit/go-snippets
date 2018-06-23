// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"time"
)

// GhUser represents a GitHub user.
type GhUser struct {
	ID                 *int64            `json:"id"`
	UserID             *int64            `json:"-"`
	GithubID           *int64            `json:"github_id"`
	Login              *string           `json:"login"`
	Bio                *string           `json:"bio"`
	Blog               *string           `json:"blog"`
	Company            *string           `json:"company"`
	Email              *string           `json:"email"`
	Hireable           *bool             `json:"hireable"`
	Location           *string           `json:"location"`
	AvatarURL          *string           `json:"avatar_url"`
	HTMLURL            *string           `json:"html_url"`
	FollowersCount     *int64            `json:"followers_count"`
	FollowingCount     *int64            `json:"following_count"`
	CollaboratorsCount *int64            `json:"collaborators_count"`
	CreatedAt          *time.Time        `json:"created_at"`
	UpdatedAt          *time.Time        `json:"updated_at"`
	GhOrganizations    []*GhOrganization `json:"gh_organizations,omitempty"`
}
