package main

import (
	"telnet"
	//"util"
	"fmt"
	//"unicode"
	"io/ioutil"
	"strings"
)

func main() {
	c, err := telnet.New("admin", "Dasan123456", "10.55.192.213", "23")
	if err != nil {
		panic(err)
	}

	c.WriteLine("enable")
	n, err := c.ReadUntil("#")
	c.WriteLine("terminal length 0")
	n, err = c.ReadUntil("#")
	if err != nil {
		panic(err)
	}

	c.WriteLine("show runnin")
	n, err = c.ReadUntil("#")
	//c.WriteLine("exit")
	//n, err = c.ReadUntil("login")
	fmt.Println(string(n))

	if err != nil {
		panic(err)
	}
	c.WriteLine("configure terminal")
	n, err = c.ReadUntil("#")

	data, _ := ioutil.ReadFile("z.c")

	lines := strings.Split(string(data), "\n")
	for _, l := range lines {
		c.WriteLine(l)
		n, err = c.ReadUntil("#")
	}

}
