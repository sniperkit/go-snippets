package redis

import "github.com/garyburd/redigo/redis"

type (
	// Redis 用于描述一个redis的数据结构
	Redis struct {
		pool *redis.Pool
	}
)

// CreateInstance 用于创建一个实例
func CreateInstance() *Redis {

}
