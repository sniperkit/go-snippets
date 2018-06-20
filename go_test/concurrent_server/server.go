// concurrent server which creates a new go routine for each of the accepted connection of the client.
// Server returns a random number for each input from the string.
package main

import (
	"math/rand"
	"fmt"
	"os"
	"time"
	"net"
	"bufio"
	"strings"
	"strconv"
)


func handleconnection(c net.Conn) {
	fmt.Printf("Serving connection: %s\n", c.RemoteAddr().String())
	for {
		netdata, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		temp := strings.TrimSpace(string(netdata))
		if temp == "STOP" {
			break
		}
		result := strconv.Itoa(rand.Intn(100)) + "\n"
		c.Write([]byte(string(result)))
	}
	c.Close()
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("No arguments provided. Error.")
		return
	}
	PORT := ":" + string(arguments[1])
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println("Error while listening to the given port.\n")
		return
	}
	defer l.Close()
	rand.Seed(time.Now().Unix())

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error while accepting a connection from the client")
			return
		}
		go handleconnection(c)
	}
}