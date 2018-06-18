// Package cache provides a unified interface to in memory cache. Under the App Engine environment,
// it uses App Engine's memcache, while under standard environments it will provide access to any number
// of cache interfaces, right now it only supports an in-memory LRU, but in the future it will support
// Redis, Memcache, and even permanent key/value stores
package cache
