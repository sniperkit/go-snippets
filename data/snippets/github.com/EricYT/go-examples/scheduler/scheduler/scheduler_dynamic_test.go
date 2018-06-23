package scheduler

import (
	"context"
	"log"
	"testing"
	"time"

	tomb "gopkg.in/tomb.v1"
)

func generalReactorRun() (*dynamicReactor, []*fakeJob) {
	reactor := NewDynamicReactor(5, 3, 2)
	var jobs []*fakeJob
	for i := 0; i < 10; i++ {
		job := &fakeJob{id: i, tomb: new(tomb.Tomb)}
		jobFunc := func(ctx context.Context) {
			select {
			case <-ctx.Done():
				return
			default:
			}
			defer job.Done()
			job.Kill(job.Run())
		}
		jobWrapper := NewJobWrapper(jobFunc, func(err error) { defer job.Done(); job.Kill(err) })
		go func(j JobWrapper) {
			err := reactor.Schedule(j)
			if err != nil {
				log.Println(err)
				j.Interrupt(err)
				return
			}
			log.Println("job go")
		}(jobWrapper)
		jobs = append(jobs, job)
	}
	return reactor, jobs
}

func TestSchedulerDynamicRun(t *testing.T) {
	log.Println("TestSchedulerDynamicRun go")
	reactor, jobs := generalReactorRun()
	defer reactor.Kill()
	waitAllJobsDone(t, jobs, time.Second*30)
	log.Println("TestSchedulerDynamicRun over")
}

func TestSchedulerDynamicRunDone(t *testing.T) {
	log.Println("TestSchedulerDynamicRunDone go")
	reactor, jobs := generalReactorRun()
	defer reactor.Kill()
	waitAllJobsDone(t, jobs, time.Second*30)
	if !reactor.Idle() {
		t.Error("all jobs already done, reactor should be idle")
	}
	log.Println("TestSchedulerDynamicRunDone over")
}

func TestSchedulerDynamicRunKill(t *testing.T) {
	reactor, _ := generalReactorRun()

	done := make(chan struct{})
	go func() {
		defer close(done)
		reactor.Kill()
	}()

	select {
	case <-done:
		log.Printf("TestSchedulerDynamicRunKill done")
	case <-time.After(time.Second * 1):
		t.Errorf("TestSchedulerDynamicRunKill wait workers done error")
	}
}
