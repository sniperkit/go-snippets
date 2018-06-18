// +build appengine

package cache

import (
    "github.com/bradberger/context"
    "github.com/bradberger/gocache/cache"
    "github.com/bradberger/gocache/codec"
    "github.com/bradberger/gocache/drivers/appengine/memcache"
)

var (
    Codec = codec.Gob    
)

func New(ctx context.Context) cache.Cache {
    memcache.Codec = Codec
    return memcache.New(ctx)
}
