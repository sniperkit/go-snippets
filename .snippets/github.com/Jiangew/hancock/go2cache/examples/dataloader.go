package main

import (
	"github.com/jiangew/hancock/go2cache"
	"strconv"
	"fmt"
)

func main() {
	cache := go2cache.Cache("myCache")

	// The data loader gets called automatically
	// whenever something tries to retrieve a non-existing key from the cache.
	cache.SetDataLoader(func(key interface{}, args ...interface{}) *go2cache.CacheItem {
		val := "This is a test with key " + key.(string)

		item := go2cache.NewCacheItem(key, val, 0)
		return item
	})

	// Let's retrieve a few auto-generated items from the cache
	for i := 0; i < 10; i++ {
		res, err := cache.Value("someKey_1" + strconv.Itoa(i))
		if err == nil {
			fmt.Println("Found value in cache:", res.Data())
		} else {
			fmt.Println("Error retrieving value from cache:", err)
		}
	}
}
