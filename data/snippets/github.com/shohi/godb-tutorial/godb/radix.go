package godb

import (
	"log"

	"github.com/mediocregopher/radix.v2/redis"
)

var radixClient *redis.Client

func init() {
	var err error
	radixClient, err = redis.Dial("tcp", "localhost:6379")
	log.Println(radixClient, err)
}
