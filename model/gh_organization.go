// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"time"
)

// GhOrganization represents a GitHub organization.
type GhOrganization struct {
	ID                 *int64     `json:"id"`
	GithubID           *int64     `json:"github_id"`
	Login              *string    `json:"login"`
	AvatarURL          *string    `json:"avatar_url"`
	HTMLURL            *string    `json:"html_url"`
	Name               *string    `json:"name"`
	Company            *string    `json:"company"`
	Blog               *string    `json:"blog"`
	Location           *string    `json:"location"`
	Email              *string    `json:"email"`
	CollaboratorsCount *int64     `json:"collaborators_count"`
	CreatedAt          *time.Time `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at"`
}
