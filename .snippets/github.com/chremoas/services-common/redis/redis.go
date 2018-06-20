package chremoas_redis

import (
	"fmt"
	"github.com/go-redis/redis"
)

type Client struct {
	Client *redis.Client
	Prefix string
}

var Nil = redis.Nil

func Init(addr string, password string, db int, prefix string) *Client {
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,  // use default DB
	})

	return &Client{Client: c, Prefix: prefix}
}

func (r Client) KeyName(key string) string {
	return fmt.Sprintf("%s:%s", r.Prefix, key)
}
