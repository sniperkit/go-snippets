package main

import (
	"net"
	"time"
	"fmt"
	"strings"
	"bufio"
	"log"
)

func echo( c net.Conn,shout string,delay time.Duration){
	fmt.Fprintln(c,"\t",strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c,"\t",shout)
	time.Sleep(delay)
	fmt.Fprintln(c,"\t",strings.ToLower(shout))
}



func handleConn2(c net.Conn){
	input := bufio.NewScanner(c)
	for input.Scan() {
		echo (c,input.Text(),1*time.Second)
	}

	c.Close()
}

func main() {
	listener ,err := net.Listen("tcp","localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn,err :=  listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn2(conn);
	}
}