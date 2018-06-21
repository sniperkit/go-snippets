package store

import (
	"runtime"
	"time"

	"github.com/garyburd/redigo/redis"
)

type (
	// Store 用于描述一个存储结构
	Store struct {
		pool *redis.Pool
	}
)

func createRedis(url string) *Store {
	return &Store{
		pool: &redis.Pool{
			MaxIdle:     runtime.NumCPU(),
			MaxActive:   200,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", url)
				if err != nil {
					return nil, err
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if time.Since(t) < time.Minute {
					return nil
				}
				_, err := c.Do("PING")
				return err
			},
		},
	}
}

// LPush yo
func (s *Store) LPush() {

}
