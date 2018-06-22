package godb

import (
	"log"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
)

func TestMiniredis(t *testing.T) {
	s, err := miniredis.Run()

	if err != nil {
		log.Println(s, err)
	}
	// log.Println(s, err)

	log.Println("address ==> ", s.Addr())
}

func TestRedisClientWithMiniredis(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic("use mock redis server error")
	}

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
		DB:   0,
	})

	log.Println(client.Ping().Result())

}
