package gocache

import (
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/allegro/bigcache"
)

var bigCache *bigcache.BigCache

func initNormal(lifeWindowSec uint32) {
	entriesInWindow := 5
	maxEntrySize := 10
	cache, err := bigcache.NewBigCache(bigcache.Config{
		Shards:             256,
		LifeWindow:         time.Duration(lifeWindowSec) * time.Second,
		CleanWindow:        1 * time.Second,
		MaxEntriesInWindow: entriesInWindow,
		MaxEntrySize:       maxEntrySize,
		Verbose:            true,
	})

	if err != nil {
		log.Println(cache, err)
	}

	bigCache = cache
}

func initZeroEvication() {
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(0))

	if err != nil {
		log.Println(cache, err)
	}

	bigCache = cache
}

func TestBigCacheForEvication(t *testing.T) {
	initNormal(2)
	bigCache.Set("hello", []byte("world"))
	log.Println(bigCache.Len())
	log.Println(bigCache.Len())

	for k := 0; k < 100; k++ {
		bigCache.Set("key"+strconv.Itoa(k), []byte("value"+strconv.Itoa(k)))
		log.Println(k, " ===> ", bigCache.Len())
	}
}

func TestBigCacheGetAndSetNil(t *testing.T) {
	initNormal(2)
	var bs []byte
	err := bigCache.Set("hello", bs)
	log.Println(err)
	log.Println(bigCache.Len())

	log.Println(bigCache.Get("hello"))

	// if key does not exist, NotFound error will occur
	log.Println(bigCache.Get("hello2"))
}

func TestBigCacheWithZeroEvication(t *testing.T) {
	initZeroEvication()

	bigCache.Set("hello", []byte("hello"))
	time.Sleep(2 * time.Second)

	bigCache.Get("hello")
	log.Printf("hits ==> %d", bigCache.Stats().Hits)
	log.Printf("miss ==> %d", bigCache.Stats().Misses)

	initNormal(1)

	bigCache.Set("hello", []byte("hello"))
	time.Sleep(2 * time.Second)
	// evication is triggered by `Set` method
	bigCache.Set("aa", []byte("aa"))
	bigCache.Get("hello")

	log.Printf("hits ==> %d", bigCache.Stats().Hits)
	log.Printf("miss ==> %d", bigCache.Stats().Misses)
}

func TestBigCacheWithNonExistKey(t *testing.T) {
	initNormal(2)

	data, err := bigCache.Get("Hello")
	log.Printf("data: %v, err: %v", data, err)
}
