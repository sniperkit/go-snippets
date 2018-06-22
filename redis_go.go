package main

import (
	"github.com/garyburd/redigo/redis"
	//"time"
	//"fmt"
	//"reflect"
	"fmt"
)

func main() {
	c, err := redis.Dial("tcp", "192.168.1.32:6379",
		redis.DialDatabase(2))
	if err != nil {
		panic(err)
	}
	defer c.Close()

	var p2 struct {
		Title  string `redis:"title"`
		Author string `redis:"author"`
		Body   string `redis:"body"`
	}

	for _, id := range []string{"id1", "id2"} {

		v, err := redis.Values(c.Do("HGETALL", id))
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := redis.ScanStruct(v, &p2); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%+v\n", p2)

	}
}