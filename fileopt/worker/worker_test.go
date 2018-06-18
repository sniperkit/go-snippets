package worker

import (
	"log"
	"sync"
	"testing"

	tomb "gopkg.in/tomb.v1"
)

var maxIndex int = 10

type generator struct {
	tomb *tomb.Tomb

	jobCh chan Jobber
}

func NewGenerate() *generator {
	g := &generator{
		tomb:  new(tomb.Tomb),
		jobCh: make(chan Jobber, 5),
	}
	go func() {
		defer g.tomb.Done()
		g.tomb.Kill(g.run())
	}()
	return g
}

func (g *generator) Generate() <-chan Jobber {
	return g.jobCh
}

func (g *generator) Wait() error {
	return g.tomb.Wait()
}

func (g *generator) run() error {
	var index int
	for {
		if index >= 10 {
			close(g.jobCh)
			return nil
		}
		job := NewJob(index)
		g.jobCh <- job
		if err := job.Wait(); err != nil {
			return err
		}
		index++
	}
}

type job struct {
	index int
	err   error
	done  chan struct{}
}

func NewJob(i int) *job {
	j := &job{
		done:  make(chan struct{}),
		index: i,
	}
	return j
}

func (j *job) Execute() error {
	log.Printf("job(%d) execute ...\n", j.index)
	return nil
}

func (j *job) Done(err error) {
	j.err = err
	close(j.done)
}

func (j *job) Wait() error {
	<-j.done
	return j.err
}

func TestGenerator(t *testing.T) {
	g := NewGenerate()
	var wg sync.WaitGroup
	wg.Add(1)
	var index int
	go func() {
		defer wg.Done()
		for {
			select {
			case job, ok := <-g.Generate():
				if !ok {
					log.Println("Generator: woo, we are here")
					return
				}
				index++
				job.Done(job.Execute())
			}
		}
	}()
	wg.Wait()
	if index != 10 {
		t.Errorf("Generator should producte 10 index but only has %d\n", index)
	}
}

func TestWorker(t *testing.T) {
	g := NewGenerate()
	go func() { NewWorker(g) }()
	go func() { NewWorker(g) }()
	if err := g.Wait(); err != nil {
		t.Errorf("Generator: generate wrong error: %s", err)
	}
}
