package test

import (
	"testing"
	"log"
)

func useBoth(db *DatabaseClient, redis *RedisClient) {
	log.Printf("[invoke] Database & Redis: %s & %s\n", db.Name, redis.Name)
}
func useRedis(redis *RedisClient) {
	log.Printf("[invoke] Redis: %s\n", redis.Name)
}
func useDatabase(db *DatabaseClient) {
	log.Printf("[invoke] Database: %s\n", db.Name)
}

func Test2(t *testing.T)  {

	injector := getInjector()
	injector.Invoke(useBoth)
	injector.Invoke(useRedis)
	injector.Invoke(useDatabase)
}
