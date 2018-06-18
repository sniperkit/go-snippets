package main

import (
	"github.com/hashicorp/consul/api"
	"context"
	"net/http"
	"log"
)

func main() {
	// creating consul client
	client, err := api.NewClient(&api.Config{Address: "127.0.0.1:8500"})

	// creating lock instance
	opts := &api.LockOptions{
		Key:        "webhook_receiver/1",
		Value:      []byte("set by sender 1"),
		SessionTTL: "10s",
		SessionOpts: &api.SessionEntry{
			Checks:   []string{"check1", "check2"},
			Behavior: "release",
		},
	}
	lock, err := client.LockOpts(opts)

	// creating a lock with all opts set to default except entry name
	//lock, err := client.LockKey("webhook_receiver/1")

	// acquiring lock
	stopCh := make(chan struct{})
	lockCh, err := lock.Lock(stopCh)
	if err != nil {
		panic(err)
	}

	cancelCtx, cancelRequest := context.WithCancel(context.Background())
	req, _ := http.NewRequest("GET", "https://example.com/webhook", nil)
	req = req.WithContext(cancelCtx)

	go func() {
		http.DefaultClient.Do(req)
		select {
		case <-cancelCtx.Done():
			log.Println("request canceled")
		default:
			log.Println("request done")

			err = lock.Unlock()
			if err != nil {
				log.Println("lock already unlocked")
			}
		}
	}()

	go func() {
		<-lockCh
		cancelRequest()
	}()

}
