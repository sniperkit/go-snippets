package main

import "github.com/gophergala/abbita/transcoder"

func main() {
	session := transcoder.NewSession("transcoder")
	server := transcoder.NewServer(session)
	server.Run()
}
