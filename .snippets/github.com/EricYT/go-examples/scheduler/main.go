package main

import (
	"context"
	"log"
	"sync"
	"time"
)

func main() {
	log.Println("main scheduler v1 go ...")
	const currence = 4
	sched := NewFIFOScheduler(currence)

	var wg sync.WaitGroup
	for index := 0; index < 10; index++ {
		wg.Add(1)
		job := func(i int) Job {
			return func(ctx context.Context) {
				defer wg.Done()
				log.Printf("Job %d running now %s", i, time.Now())
				//time.Sleep(time.Second * time.Duration(i))
				time.Sleep(time.Second * 2)
				//log.Printf("Job %d running done now %s", i, time.Now())
			}
		}(index)
		sched.Schedule(job)
	}

	log.Println("main schedule done")
	wg.Wait()

	log.Println("pendings ", sched.Pending())
	log.Println("finished ", sched.Finished())
	log.Println("schedued ", sched.Scheduled())
}
