package godb

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/go-redis/redis"
)

/*
 redis is used as cluster key-value storage, focus on following concerns:
 1. read -- goroutine safe, transaction support. can't modified when read,
 2. write -- goroutine safe, transaction supoort. only one write is allowed, write success or roll back
 3. read all -- goroutine safe, transaction supoort. can't modified when read
 4. purge -- goroutine safe, transaction support. only one goroutines is allowed. purge success or roll back.

 *performance*
 1. read batch
 2. write batch

 Note: keys are not ordered in redis.
*/

// add data to redis for testing
func init() {
	var key, value string
	redisClient.FlushDB()
	for k := 0; k < 10; k++ {
		key = fmt.Sprintf("key%d", k+1)
		value = fmt.Sprintf("value%d", k+1)
		if err := addRedisKV(key, value); err != nil {
			panic("Failed to add KV to redis, err: " + err.Error())
		}
	}
}

func TestRedisPing(t *testing.T) {
	ping()
}
func TestGetRedisKV(t *testing.T) {
	log.Println(redisClient.PoolStats())
	log.Println(getRedisKV("key1"))
}

func TestGetRedisKVForOtherTypes(t *testing.T) {
	redisClient.Set("intKey1", 1, 0)
	redisClient.Set("intKey2", 1, 0)
	redisClient.Set("intKey3", 1, 0)

	log.Println(redisClient.Get("intKey1").Bytes())
}

func TestAddRedisKVWithTx(t *testing.T) {
	addRedisKVWithTx("txKey1", "txValue1")
	log.Println(getRedisKV("txKey1"))
}

func TestGetAllRedisKV(t *testing.T) {
	m, err := getAllRedisKV()
	if err != nil {
		t.Errorf("Failed to get all KVs from redis, err: %v", err)
	}
	for k, v := range m {
		log.Println("key ==> ", k, " value ==> ", v)
	}
}

func TestPurgeRedis(t *testing.T) {
	err := emptyRedisDB()
	if err != nil {
		t.Errorf("Failed to purge redis db")
	}

	value, err := getRedisKV("key1")
	if err == nil {
		t.Errorf("After purge, there is no <k, v>, value: %v", value)
	}
}

func TestSetRedis(t *testing.T) {
	// redis client support several data types
	redisClient.Set("hello", 1, 0)
	val := redisClient.Get("hello").Val()
	log.Println(strconv.Atoi(val))
}

func TestGetAllRedisKVWithTx(t *testing.T) {
	m, err := getAllRedisKVWithTx()
	if err != nil {
		log.Println(err)
	} else {
		for k, v := range m {
			log.Println(k, v)
		}
	}

}

func TestAddAndIncreaseRedisTx(t *testing.T) {
	redisClient.FlushDB()
	addAndIncreaseRedisTx()
}

func TestGetNonExistsKeyRedis(t *testing.T) {
	redisClient.FlushDB()
	val, err := redisClient.Get("NON-EXIST-INT").Int64()

	if err != nil && err == redis.Nil {
		log.Println("non exists")
	}
	log.Println(val, err)
}

func TestDelNonExistsKeyRedis(t *testing.T) {
	redisClient.FlushDB()

	log.Println(redisClient.Del("NON-EXIST-KEY"))
}

func TestRedisConnection(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	log.Println(client)
	_, err := client.Ping().Result()
	log.Println(err)
	client.Close()

	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	info, err := client.Ping().Result()
	if err != nil || !strings.EqualFold(info, "PONG") {
		log.Println(info, err)
	}

}

func TestRedisWatch(t *testing.T) {
	var incr func(string) error
	// Transactionally increments key using GET and SET commands.
	incr = func(key string) error {
		err := redisClient.Watch(func(tx *redis.Tx) error {
			n, err := tx.Get(key).Int64()
			log.Println("watch ====> ", n)
			if err != nil && err != redis.Nil {
				return err
			}

			_, err = tx.Pipelined(func(pipe redis.Pipeliner) error {
				pipe.Set(key, strconv.FormatInt(n+1, 10), 0)
				return nil
			})
			return err
		}, key)
		if err == redis.TxFailedErr {
			return incr(key)
		}
		return err
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			err := incr("key")
			if err != nil {
				log.Println(err)
			} else {
				log.Println("success ==> ", k)
			}

		}(i)
	}
	wg.Wait()
}

func TestRedisWatchWithMultipleClient(t *testing.T) {
	var incr func(*redis.Client, string) error
	// Transactionally increments key using GET and SET commands.
	incr = func(redisClient *redis.Client, key string) error {
		err := redisClient.Watch(func(tx *redis.Tx) error {
			n, err := tx.Get(key).Int64()
			log.Println("watch ====> ", n)
			if err != nil && err != redis.Nil {
				return err
			}

			_, err = tx.Pipelined(func(pipe redis.Pipeliner) error {
				pipe.Set(key, strconv.FormatInt(n+1, 10), 0)
				return nil
			})
			return err
		}, key)
		if err == redis.TxFailedErr {
			return incr(redisClient, key)
		}
		return err
	}

	// TEST
	var client *redis.Client
	client1 := newClientMust()
	client2 := newClientMust()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)

		if rand.Float32() < 0.5 {
			client = client1
			log.Println("client ==> client1")
		} else {
			client = client2
			log.Println("client ==> client2")
		}

		go func(k int) {
			defer wg.Done()
			err := incr(client, "key")
			if err != nil {
				log.Println(err)
			} else {
				log.Println("success ==> ", k)
			}

		}(i)
	}
	wg.Wait()
}

func TestSMembersMapForRedis(t *testing.T) {
	redisClient.SAdd("myset", "hello")
	redisClient.SAdd("myset", "world")

	res, _ := redisClient.SMembersMap("myset").Result()
	log.Println(res)
}
