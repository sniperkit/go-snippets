package go2cache

import (
	"sync"
	"time"
)

// CacheItem is an individual cache item.
type CacheItem struct {
	sync.RWMutex

	// map的key，必须是可比较类型
	key interface{}

	// 缓存值可以是任何类型
	data interface{}

	// 每条缓存的生命周期
	lifeSpan time.Duration

	// 缓存项创建时间戳
	createdOn time.Time

	// 缓存项上次被访问时间戳
	accessedOn time.Time

	// 每条缓存的访问记数
	accessCount int64

	// 缓存项被删除时的回调函数「删除之前执行」
	aboutToExpire func(key interface{})
}

// NewCacheItem returns a newly created CacheItem.
func NewCacheItem(key interface{}, data interface{}, lifeSpan time.Duration) *CacheItem {
	t := time.Now()
	return &CacheItem{
		key:           key,
		data:          data,
		lifeSpan:      lifeSpan,
		aboutToExpire: nil,
		createdOn:     t,
		accessedOn:    t,
		accessCount:   0,
	}
}

// KeepAlive marks an item to be kept for another expireDuration period.
func (item *CacheItem) KeepAlive() {
	item.Lock()
	defer item.Unlock()
	item.accessedOn = time.Now()
	item.accessCount++
}

// LifeSpan returns this item's expiration duration.
func (item *CacheItem) LifeSpan() time.Duration {
	return item.lifeSpan
}

// CreatedOn returns when this item was added to the cache.
func (item *CacheItem) CreatedOn() time.Time {
	return item.createdOn
}

// AccessedOn returns when this item was last accessed.
func (item *CacheItem) AccessedOn() time.Time {
	item.RLock()
	defer item.RUnlock()
	return item.accessedOn
}

// AccessCount returns how often this item has been accessed.
func (item *CacheItem) AccessCount() int64 {
	item.RLock()
	defer item.RUnlock()
	return item.accessCount
}

// Key returns the key of this cached item.
func (item *CacheItem) Key() interface{} {
	return item.key
}

// Data returns the value of this cached item.
func (item *CacheItem) Data() interface{} {
	return item.data
}

// SetAboutToExpireCallback configures a callback, which will be called right before the item is about to be removed from the cache.
func (item *CacheItem) SetAboutToExpireCallback(f func(interface{})) {
	item.Lock()
	defer item.Unlock()
	item.aboutToExpire = f
}
