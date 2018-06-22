package config

type redisConfig struct {
	Address  string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// RedisConfig -- redis config
var RedisConfig redisConfig
