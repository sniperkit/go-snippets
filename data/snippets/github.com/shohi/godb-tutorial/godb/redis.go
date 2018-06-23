package godb

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/go-redis/redis"
	"github.com/shohi/godb-tutorial/config"
)

var redisClient *redis.Client

func init() {
	log.Println("init in redis...")

	redisClient = newClientMust()

}

func newClientMust() *redis.Client {
	bs, err := ioutil.ReadFile("../config/config.json")
	if err != nil {
		panic("can't load config.json, err: " + err.Error())
	}
	err = json.Unmarshal(bs, &config.RedisConfig)

	if err != nil {
		panic("can't unmarshal config to RedisConfig")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisConfig.Address,
		Password: config.RedisConfig.Password,

		// DB will be changed to [0, 15]
		DB: config.RedisConfig.DB, // actually, this will be used separately, e.g. 0: for active vault, 1: for archive vault, 2: for recon vault

	})

	return client
}

func ping() {
	log.Println(redisClient.Ping())
}

func addRedisKV(key string, value string) error {
	return redisClient.Set(key, value, 0).Err()
}

func addRedisKVWithTx(key string, value string) error {
	// redisClient.TxPipeline(func())
	pipe := redisClient.TxPipeline()
	pipe.Set(key, value, 0)

	// pipe.Discard()
	cmds, err := pipe.Exec()
	log.Println(cmds)
	log.Println("val ==> ", cmds[0].String(), cmds[0].Err())

	return err
}

func deleteRedisKV(key string) error {
	return redisClient.Del(key).Err()
}

func getRedisKV(key string) (string, error) {
	bs, err := redisClient.Get(key).Bytes()
	if err != nil {
		return "", err
	}
	return string(bs), err
}

func getAllRedisKV() (map[string]string, error) {
	m := make(map[string]string)
	keys, err := redisClient.Keys("*").Result()
	if err != nil {
		return m, err
	}

	for _, key := range keys {
		value, err := getRedisKV(key)
		if err != nil {
			return m, err
		}
		m[key] = value
	}

	return m, err
}

func emptyRedisDB() error {
	return redisClient.FlushDB().Err()
}

func getAllRedisKVWithTx() (map[string]string, error) {

	m := make(map[string]string)
	pipe := redisClient.TxPipeline()

	keys, err := pipe.Keys("*").Result()
	log.Println(keys, err)

	pipe.Set("keytx", "hello2", 0)
	// pipe.Exec()
	pipe.Discard()

	if err != nil {
		pipe.Discard()
		return m, err
	}

	for _, key := range keys {
		value, err := pipe.Get(key).Result()
		if err != nil {
			pipe.Discard()
			return m, err
		}
		m[key] = value
	}
	_, err = pipe.Exec()

	return m, err
}

func addAndIncreaseRedisTx() error {

	//
	var err error
	pipe := redisClient.TxPipeline()

	// update count
	pipe.Set("keyIncTx", "valIncTx", 0)
	// update verison
	pipe.Incr("keyIncInt")

	pipe.Decr("keyDecrInt")

	cmds, err := pipe.Exec()
	log.Println(cmds)

	// err = pipe.Discard()
	// log.Printf("discard after exec, err: %v\n", err)
	return err
}
