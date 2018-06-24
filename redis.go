package main

import (
	"encoding/json"
	"fmt"
	"sync"
	// "github.com/datatogether/task-mgmt/tasks"
	"github.com/garyburd/redigo/redis"
	// "net"
	"time"
)

// Main redis connection
var rconn redis.Conn

// func connectRedis() (err error) {
// 	var netConn net.Conn

// 	if cfg.RedisUrl == "" {
// 		return fmt.Errorf("no redis url specified")
// 	}

// 	for i := 0; i <= 1000; i++ {
// 		netConn, err = net.Dial("tcp", cfg.RedisUrl)
// 		if err != nil {
// 			log.Infoln("error connecting to redis: %s", err.Error())
// 			time.Sleep(time.Second)
// 			continue
// 		}
// 		break
// 	}

// 	if netConn == nil {
// 		return fmt.Errorf("no net connection after 1000 tries")
// 	}

// 	log.Infoln("connected to redis")

// 	rconn = redis.NewConn(netConn, time.Second*20, time.Second*20)
// 	return SubscribeTaskProgress(rconn)
// }

func SubscribeTaskProgress() (err error) {
	var conn redis.Conn

	if cfg.RedisUrl == "" {
		return fmt.Errorf("no redis url specified")
	}

	for i := 0; i <= 1000; i++ {
		conn, err = redis.Dial("tcp", cfg.RedisUrl,
			redis.DialReadTimeout(10*time.Second),
			redis.DialWriteTimeout(0),
			redis.DialConnectTimeout(20*time.Second),
		)
		if err != nil {
			log.Infoln("error connecting to redis: %s", err.Error())
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}

	if conn == nil {
		return fmt.Errorf("couldn't connect to redis after 1000 tries")
	}

	defer conn.Close()
	var wg sync.WaitGroup
	wg.Add(2)

	log.Infoln("connected to redis")
	psc := redis.PubSubConn{Conn: conn}
	if err = psc.PSubscribe("tasks.*"); err != nil {
		return err
	}
	defer psc.PUnsubscribe()

	go func() {
		defer wg.Done()
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				// log.Infof("%s: message: %s\n", v.Channel, v.Data)

				// TODO - other types of messages will eventually come through
				// here...

				res := &ClientResponse{
					Type:      "TASK_PROGRESS",
					RequestId: "server",
					Schema:    "TASK",
					Data:      json.RawMessage(v.Data),
				}

				data, err := json.Marshal(res)
				if err != nil {
					log.Infoln(err.Error())
				} else {
					room.broadcast <- data
				}
			case redis.PMessage:
				// log.Infof("PMessage: %s %s %s\n", v.Pattern, v.Channel, v.Data)

				// TODO - other types of messages will eventually come through
				// here...

				res := &ClientResponse{
					Type:      "TASK_PROGRESS",
					RequestId: "server",
					Schema:    "TASK",
					Data:      json.RawMessage(v.Data),
				}

				data, err := json.Marshal(res)
				if err != nil {
					log.Infoln(err.Error())
				} else {
					room.broadcast <- data
				}
			case redis.Pong:
				// log.Infof("received pong")
			case redis.Subscription:
				log.Infof("%s: %s %d\n", v.Channel, v.Kind, v.Count)
			case error:
				log.Infoln("connection error:", conn.Err())
				log.Infoln("message error: %s", v.Error())
				return
			}
		}
	}()

	go func() {
		for {
			if err := psc.Ping("PING"); err != nil {
				log.Infoln("error sending ping")
				return
			}
			time.Sleep(time.Second * 8)
		}
	}()

	wg.Wait()
	return nil
}
