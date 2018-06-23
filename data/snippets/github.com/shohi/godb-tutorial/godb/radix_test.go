package godb

import (
	"fmt"
	"log"
	"testing"
)

func init() {
	var key, value string
	for k := 0; k < 10; k++ {
		key = fmt.Sprintf("rkey%d", k+1)
		value = fmt.Sprintf("rvalue%d", k+1)
		if err := radixClient.Cmd("SET", key, value).Err; err != nil {
			panic("Failed to add KV to redis, err: " + err.Error())
		}
	}

}

func TestRadixInit(t *testing.T) {
	log.Println("init radix...")
	log.Println(radixClient.Cmd("GET", "rkey1"))
}

func TestRadixTx(t *testing.T) {

	m := make(map[string]string)
	// start transaction
	err := radixClient.Cmd("MULTI").Err
	if err != nil {
		log.Println(err)
		return
	}

	for k := 0; k < 11; k++ {
		key := fmt.Sprintf("rkey%d", k+1)
		bs, err := radixClient.Cmd("GET", key).Str()
		if err != nil {
			return
		}
		m[key] = string(bs)
	}

	radixClient.Cmd("SET", "world", "world")
	radixClient.Cmd("EXEC")

	/*
		radixClient.Cmd("Hello")
		ereply := radixClient.Cmd("EXEC")
		if ereply.Err != nil {
			log.Println(ereply)
		} else if ereply.IsType(redis.Nil) {
			log.Println(redis.Nil)
		}
	*/

	//
	log.Println("output........")
	for k, v := range m {
		log.Printf("key=%s, value=%s\n", k, v)
	}
}
