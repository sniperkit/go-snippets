package builder

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"syscall"
	"time"

	"github.com/iron-io/iron_go/mq"
)

type Image struct {
	Id     string `json:"id"`
	Tag    string `json:"tag"`
	Source string `json:"source"`
}

type Message struct {
	Body string
	Id   string
}

func UtsnameToByte(name [65]int8) []byte {
	b := make([]byte, len(name))
	i := 0
	for ; i < len(name); i++ {
		if name[i] == 0 {
			break
		}
		b[i] = byte(name[i])
	}
	return b[:i]
}

func GetMachineName() string {
	uts := &syscall.Utsname{}
	err := syscall.Uname(uts)
	if err != nil {
		fmt.Println(err)
	}

	return string(UtsnameToByte(uts.Machine))
}

func GetBuildRequests(queue *mq.Queue, wait time.Duration) <-chan Message {
	c := make(chan Message)

	go func() {
		for {
			msg, err := queue.Get()
			if err != nil {
				time.Sleep(wait)
			} else {
				fmt.Println("Request: " + msg.Body)
				c <- Message{msg.Body, msg.Id}
			}
		}
	}()
	return c
}

func GetDockerfile(imageId string) {
	u := "http://localhost:3000/api/v1/build/" + imageId + "/dockerfile"
	res, err := http.Get(u)
	if err != nil {
		log.Fatalf("error retrieving: %s", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error reading body: %s", err)
	}
	log.Println(string(body))
}

func Run() {
	machine := GetMachineName()
	queueName := string("builder-" + machine)

	queue := mq.New(queueName)

	request := GetBuildRequests(queue, 10*time.Second)

	for {
		select {
		case msg := <-request:
			log.Println(msg.Body)
			imageId := msg.Body
			go GetDockerfile(imageId)
			queue.DeleteMessage(msg.Id)
		}
	}
}
