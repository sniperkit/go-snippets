// +build !appengine

package cache

import (
	"github.com/bradberger/context"
	"github.com/bradberger/gocache"
	"github.com/bradberger/gocache/cache"
	"github.com/bradberger/gocache/drivers/lru"
)

type CacheDriver int

const (
	CacheDriverLRU CacheDriver = iota
)

var (
	gc *gocache.Client

	Driver     CacheDriver
	Memory     = 50 * lru.Megabyte
	MaxEntries = 2 << 10
	Nodes      int
)

func New(ctx context.Context) cache.Cache {
	if gc != nil {
		return gc
	}

	gc = gocache.New()
	switch Driver {
	default:
		gc.AddNode("node-01", lru.NewBasic(Memory, MaxEntries))
	}

	return gc
}
