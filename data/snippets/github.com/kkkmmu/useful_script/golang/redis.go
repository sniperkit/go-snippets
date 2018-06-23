package main

import (
	"gopkg.in/redis.v5"
	"log"
	"os"
)

var redisClient *redis.Client

func main() {
	result, err := redisClient.HGetAll("SPIDER:RESULT:CACHE").Result()
	if err != nil {
		panic(err)
	}

	file, _ := os.Create("test.txt")
	for k, v := range result {
		log.Println(k, v)
		file.WriteString(k + "\n")
	}
}

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
