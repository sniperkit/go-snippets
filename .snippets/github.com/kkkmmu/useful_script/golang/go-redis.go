//
// author xiongyongshun
// project go_redis
// version 1.0
// created 16/10/6 03:49
//
package main

import (
	"fmt"
	"gopkg.in/redis.v5"
	"sync"
	"time"
)

func main() {
	client := createClient()
	defer client.Close()

	stringOperation(client)
	listOperation(client)
	setOperation(client)
	hashOperation(client)

	connectPool(client)

}

// 创建 redis 客户端
func createClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 5,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	return client
}

// String 操作
func stringOperation(client *redis.Client) {
	// 第三个参数是过期时间, 如果是0, 则表示没有过期时间.
	err := client.Set("name", "xys", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get("name").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("name", val)

	// 这里设置过期时间.
	err = client.Set("age", "20", 1*time.Second).Err()
	if err != nil {
		panic(err)
	}

	client.Incr("age") // 自增
	client.Incr("age") // 自增
	client.Decr("age") // 自减

	val, err = client.Get("age").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("age", val) // age 的值为21

	// 因为 key "age" 的过期时间是一秒钟, 因此当一秒后, 此 key 会自动被删除了.
	time.Sleep(1 * time.Second)
	val, err = client.Get("age").Result()
	if err != nil {
		// 因为 key "age" 已经过期了, 因此会有一个 redis: nil 的错误.
		fmt.Printf("error: %v\n", err)
	}
	fmt.Println("age", val)
}

// list 操作
func listOperation(client *redis.Client) {
	client.RPush("fruit", "apple")               //在名称为 fruit 的list尾添加一个值为value的元素
	client.LPush("fruit", "banana")              //在名称为 fruit 的list头添加一个值为value的 元素
	length, err := client.LLen("fruit").Result() //返回名称为 fruit 的list的长度
	if err != nil {
		panic(err)
	}
	fmt.Println("length: ", length) // 长度为2

	value, err := client.LPop("fruit").Result() //返回并删除名称为 fruit 的list中的首元素
	if err != nil {
		panic(err)
	}
	fmt.Println("fruit: ", value)

	value, err = client.RPop("fruit").Result() // 返回并删除名称为 fruit 的list中的尾元素
	if err != nil {
		panic(err)
	}
	fmt.Println("fruit: ", value)
}

// set 操作
func setOperation(client *redis.Client) {
	client.SAdd("blacklist", "Obama")     // 向 blacklist 中添加元素
	client.SAdd("blacklist", "Hillary")   // 再次添加
	client.SAdd("blacklist", "the Elder") // 添加新元素

	client.SAdd("whitelist", "the Elder") // 向 whitelist 添加元素

	// 判断元素是否在集合中
	isMember, err := client.SIsMember("blacklist", "Bush").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Is Bush in blacklist: ", isMember)

	// 求交集, 即既在黑名单中, 又在白名单中的元素
	names, err := client.SInter("blacklist", "whitelist").Result()
	if err != nil {
		panic(err)
	}
	// 获取到的元素是 "the Elder"
	fmt.Println("Inter result: ", names)

	// 获取指定集合的所有元素
	all, err := client.SMembers("blacklist").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("All member: ", all)
}

// hash 操作
func hashOperation(client *redis.Client) {
	client.HSet("user_xys", "name", "xys") // 向名称为 user_xys 的 hash 中添加元素 name
	client.HSet("user_xys", "age", "18")   // 向名称为 user_xys 的 hash 中添加元素 age

	// 批量地向名称为 user_test 的 hash 中添加元素 name 和 age
	client.HMSet("user_test", map[string]string{"name": "test", "age": "20"})
	// 批量获取名为 user_test 的 hash 中的指定字段的值.
	fields, err := client.HMGet("user_test", "name", "age").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("fields in user_test: ", fields)

	// 获取名为 user_xys 的 hash 中的字段个数
	length, err := client.HLen("user_xys").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("field count in user_xys: ", length) // 字段个数为2

	// 删除名为 user_test 的 age 字段
	client.HDel("user_test", "age")
	age, err := client.HGet("user_test", "age").Result()
	if err != nil {
		fmt.Printf("Get user_test age error: %v\n", err)
	} else {
		fmt.Println("user_test age is: ", age) // 字段个数为2
	}
}

// redis.v4 的连接池管理
func connectPool(client *redis.Client) {
	wg := sync.WaitGroup{}
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()

			for j := 0; j < 100; j++ {
				client.Set(fmt.Sprintf("name%d", j), fmt.Sprintf("xys%d", j), 0).Err()
				client.Get(fmt.Sprintf("name%d", j)).Result()
			}

			fmt.Printf("PoolStats, TotalConns: %d, FreeConns: %d\n", client.PoolStats().TotalConns, client.PoolStats().FreeConns)
		}()
	}

	wg.Wait()
}
