package main

import (
	"encoding/json"
	"fmt"
	"github.com/datatogether/task_mgmt/tasks"
	"github.com/garyburd/redigo/redis"
)

// Main redis connection
var rpool *redis.Pool

var ErrNoRedisConn = fmt.Errorf("No connection to redis could be found")

func connectRedis() (err error) {
	if cfg.RedisUrl == "" {
		return fmt.Errorf("no redis url specified")
	}

	rpool = redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", cfg.RedisUrl)
		return c, err
	}, 3)

	return nil
}

func PublishTaskProgress(pool *redis.Pool, t *tasks.Task) error {
	if pool == nil {
		return ErrNoRedisConn
	}

	c := pool.Get()
	defer c.Close()

	data, err := json.Marshal(t)
	if err != nil {
		return err
	}

	// log.Infof("publishing on channel %s: %s", t.PubSubChannelName(), string(data))
	_, err = c.Do("PUBLISH", t.PubSubChannelName(), data)
	return err
}
