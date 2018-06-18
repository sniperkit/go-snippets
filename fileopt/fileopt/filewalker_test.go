package parallel

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"testing"
	"time"

	tomb "gopkg.in/tomb.v1"

	"github.com/EricYT/go-examples/fileopt/worker"
)

var tmux sync.Mutex
var jobs map[string]*job = make(map[string]*job)

func AddJob(j *job) {
	tmux.Lock()
	jobs[j.Path()] = j
	tmux.Unlock()
}

func RemoveJob(path string) {
	tmux.Lock()
	delete(jobs, path)
	tmux.Unlock()
}

func filterFun(path string, info os.FileInfo, err error) (worker.Jobber, error) {
	if info.IsDir() {
		log.Printf("parallel: path(%s) is directory. Skip it", path)
		return nil, filepath.SkipDir
	}

	match, _ := regexp.MatchString("foo.*", path)
	if match {
		log.Printf("parallel: path(%s) skip\n", path)
		j := NewJob(path)
		AddJob(j)
		return j, nil
	}

	return nil, nil
}

type job struct {
	tomb *tomb.Tomb
	path string
}

func NewJob(path string) *job {
	return &job{
		tomb: new(tomb.Tomb),
		path: path,
	}
}

func (j *job) Path() string { return j.path }

func (j *job) Execute() error {
	log.Printf("Job: path(%s) execute\n", j.path)
	//FIXME: Oh my god, dangerous
	return os.Remove(j.path)
}

func (j *job) Done(err error) {
	j.tomb.Kill(err)
	RemoveJob(j.path)
}

func (j *job) Wait() error {
	return j.tomb.Wait()
}

func TestFileGenerator(t *testing.T) {
	fg := NewFileGenerator(os.TempDir(), filterFun)
	worker.NewWorker(fg)
	if err := fg.Wait(); err != nil {
		t.Errorf("TestFileGenerator: file generator run error: %s", err)
	}
	for {
		if len(jobs) != 0 {
			log.Printf("TestFileGenerator: jobs are not done (%d)\n", len(jobs))
			time.Sleep(500 * time.Millisecond)
			continue
		}
		break
	}
}
