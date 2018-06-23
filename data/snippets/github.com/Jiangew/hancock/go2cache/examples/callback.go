package main

import (
	"github.com/jiangew/hancock/go2cache"
	"fmt"
	"time"
)

func main() {
	cache := go2cache.Cache("myCache")

	// This callback will be triggered every time a new item gets added to the cache
	cache.SetAddedItemCallback(func(entry *go2cache.CacheItem) {
		fmt.Println("Added:", entry.Key(), entry.Data(), entry.CreatedOn())
	})

	// This callback will be triggered every time an item is about to be removed from the cache
	cache.SetAboutToDeleteItemCallback(func(entry *go2cache.CacheItem) {
		fmt.Println("Deleting:", entry.Key(), entry.Data(), entry.CreatedOn())
	})

	// Caching a new item will execute the AddedItem callback
	cache.Add("someKey", "JamesiWorks test Local Cache", 0)

	res, err := cache.Value("someKey")
	if err == nil {
		fmt.Println("Found value in cache:", res.Data())
	} else {
		fmt.Println("Error retrieving value from cache:", err)
	}

	// Deleting the item will execute the AboutToDeleteItem callback
	cache.Delete("someKey")

	// Caching a new item that expires in 3 seconds
	res = cache.Add("anotherKey", "This is another JamesiWorks test Local Cache", 3*time.Second)

	// This callback will be triggered when the item is about to expire
	res.SetAboutToExpireCallback(func(key interface{}) {
		fmt.Println("About to expire:", key.(string))
	})

	time.Sleep(5 * time.Second)
}
