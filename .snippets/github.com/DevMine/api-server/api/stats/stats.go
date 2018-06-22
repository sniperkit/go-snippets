// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package stats handles /stats... routes.
package stats

import (
	"net/http"

	"github.com/DevMine/api-server/cache"
	"github.com/DevMine/api-server/srv/context"
	"github.com/DevMine/api-server/util/json"
)

// Index handles "/stats" route.
func Index(c *context.Context, w http.ResponseWriter, r *http.Request) {
	w.Write(json.MarshalIndentPanic(cache.GetStats()))
}
