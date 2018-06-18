package scheduler

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"

	tomb "gopkg.in/tomb.v1"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type fakeJob struct {
	tomb *tomb.Tomb
	id   int
}

func (f *fakeJob) Run() error {
	//duration := time.Second * time.Duration(rand.Intn(3))
	duration := time.Second * time.Duration(2)
	select {
	case <-f.tomb.Dying():
		log.Printf("fake job: ready to die")
		return nil
	case <-time.After(duration):
		log.Printf("fake job %d done now: %s", f.id, time.Now())
		return nil
	}
}

func (f *fakeJob) Kill(err error) {
	f.tomb.Kill(err)
}

func (f *fakeJob) Wait() error {
	return f.tomb.Wait()
}

func (f *fakeJob) Done() {
	f.tomb.Done()
}

func waitAllJobsDone(t *testing.T, jobs []*fakeJob, timeout time.Duration) {
	log.Printf("wait jobs: %d done", len(jobs))
	var wg sync.WaitGroup
	wg.Add(len(jobs))
	for _, job := range jobs {
		go func(j *fakeJob) {
			defer wg.Done()
			j.Wait()
		}(job)
	}
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// wait
	select {
	case <-time.After(timeout):
		t.Errorf("wait all jobs done timeout")
	case <-done:
		// all done
	}
}

func TestSchdulerRunWithoutOverflow(t *testing.T) {
	log.Printf("\ntest scheduler run")
	sch := NewReactor(3, 100)
	var jobs []*fakeJob
	for i := 0; i < 20; i++ {
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
		err := sch.Schedule(jobWrapper)
		if err != nil {
			continue
		}
		jobs = append(jobs, job)
	}
	waitAllJobsDone(t, jobs, time.Second*30)
}

func TestSchdulerRunOverflow(t *testing.T) {
	log.Printf("\ntest scheduler run overflow")
	sch := NewReactor(3, 10)
	var jobs []*fakeJob
	for i := 0; i < 20; i++ {
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
		err := sch.Schedule(jobWrapper)
		if err != nil {
			log.Printf("job %d add error: %s", i, err)
			continue
		}
		jobs = append(jobs, job)
	}
	waitAllJobsDone(t, jobs, time.Second*30)
}

func TestSchedulerKill(t *testing.T) {
	log.Print("\ntest scheduler run kill")
	sch := NewReactor(3, 100)
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
		err := sch.Schedule(jobWrapper)
		if err != nil {
			log.Printf("job %d add error: %s", i, err)
			continue
		}
		jobs = append(jobs, job)
	}
	time.AfterFunc(time.Second*3, func() { sch.Kill() })
	waitAllJobsDone(t, jobs, time.Second*30)
}
