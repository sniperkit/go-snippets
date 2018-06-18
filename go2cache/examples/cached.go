package main

import (
	"github.com/jiangew/hancock/go2cache"
	"time"
	"fmt"
)

type moreData struct {
	text string
	ext  []byte
}

func main() {
	cache := go2cache.Cache("myCache")

	val := moreData{"JamesiWorks Test", []byte{}}
	cache.Add("someKey", &val, 5*time.Second)

	res, err := cache.Value("someKey")
	if err == nil {
		fmt.Println("Found value in cache:", res.Data().(*moreData).text)
	} else {
		fmt.Println("Error retrieving value from cache:", err)
	}

	time.Sleep(6 * time.Second)
	res, err = cache.Value("someKey")
	if err != nil {
		fmt.Println("Item is not cached (anymore).")
	}

	// Add another item that never expires
	cache.Add("someKey", &val, 0)

	// go2cache supports a few handy callbacks and loading mechanisms
	cache.SetAboutToDeleteItemCallback(func(e *go2cache.CacheItem) {
		fmt.Println("Deleting:", e.Key(), e.Data().(*moreData).text, e.CreatedOn())
	})

	cache.Delete("someKey")

	cache.Flush()
}
