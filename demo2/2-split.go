package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"time"
)

func main() {
	var split []string

	for i := 0; i < 10; i++ {
		split = append(split, fmt.Sprintf("%d", i))
	}
	log.Println(split)
	split = split[len(split):]
	log.Println(split)
	m := md5.Sum([]byte(fmt.Sprintf("%d", time.Now().Unix())))
	log.Println(m)
	log.Println(fmt.Sprintf("%x", m))
}
