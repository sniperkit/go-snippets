// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"net/http"

	"github.com/DevMine/api-server/srv/context"
	"github.com/DevMine/api-server/util/json"
)

type api struct {
	Version int    `json:"version"`
	DocURL  string `json:"doc_url"`
}

// Index handles "/" route.
func Index(c *context.Context, w http.ResponseWriter, r *http.Request) {
	w.Write(json.MarshalIndentPanic(api{Version: Version, DocURL: DocURL}))
}
