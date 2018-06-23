package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

var chars []byte = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var charsLen int = len(chars)

func randomBody(size int) []byte {
	body := make([]byte, size)

	var index int
	for index < size {
		charIndex := rand.Intn(charsLen)
		body[index] = chars[charIndex]
		index++
	}
	return body
}

func calculateBodyMd5(data []byte) string {
	md5tmp := md5.Sum(data)
	//FIXME: the same as use hex.EncodeToString(md5tmp[:])
	var md5 string
	for _, c := range md5tmp {
		part := fmt.Sprintf("%x", c)
		md5 = md5 + part
	}
	return md5
}

func main() {
	log.Println("iris go")
	rand.Seed(time.Now().Unix())
	app := iris.New()

	app.Handle("POST", "/ping", func(ctx context.Context) {
		body := randomBody(4096)
		//log.Printf("data: %s\n", body)
		bodyMd5 := calculateBodyMd5(body)
		ctx.Header("X-Data-Md5", bodyMd5)
		ctx.Write(body)
	})
	app.Run(iris.Addr(":9880"))
}
