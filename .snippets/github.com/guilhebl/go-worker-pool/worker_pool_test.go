package job

import (
	"strconv"
	"testing"
	"strings"
	"unicode"
	"math/rand"
)

// Anonymous Task
type TestTask struct{}

func (e *TestTask) Run(payload Payload) JobResult {
	return NewJobResult("job done!", nil)
}

func NewTestTask() TestTask {
	return TestTask{}
}

// Anonymous Sum Task
type TestSumTask struct{}

func (e *TestSumTask) Run(payload Payload) JobResult {
	x, _ := strconv.ParseInt(payload.Params["x"], 10, 0)
	y, _ := strconv.ParseInt(payload.Params["y"], 10, 0)
	return NewJobResult(x+y, nil)
}

func NewTestSumTask() TestSumTask {
	return TestSumTask{}
}

// Creates a sample sum Job task which aims to sum X and Y as int
func NewSumJobTask(x, y int) Job {
	ret := NewJobResultChannel()
	m := make(map[string]string)
	m["x"] = strconv.Itoa(x)
	m["y"] = strconv.Itoa(y)
	task := NewTestSumTask()
	work := NewJob(&task, m, ret)
	return work
}

// Creates a sample task that sums all the Rune int values of a string
func NewSumRunesOfWordTask() SumRuneTask {
	return SumRuneTask{}
}

func NewSumRunesOfWordJobTask(word string) Job {
	ret := NewJobResultChannel()
	m := make(map[string]string)
	m["word"] = word
	task := NewSumRunesOfWordTask()
	work := NewJob(&task, m, ret)
	return work
}

// Anonymous SumRune Task
type SumRuneTask struct{}

func spaceMap(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func (e *SumRuneTask) Run(payload Payload) JobResult {
	word := payload.Params["word"]
	word = spaceMap(word)

	res := int64(0)
	for _, r := range word {
		num := int64(r)
		res += num
	}
	return NewJobResult(res, nil)
}

func testClosePool(t *testing.T, p *WorkerPool) {
	// try to close pool
	ok := p.Stop()

	if ok {
		t.Error("WorkerPool not closed")
	}
}

// just creates a pool but doesn't run it tests open and close
func TestCreatePool(t *testing.T) {

	p := NewWorkerPool(4)

	if &p == nil {
		t.Error("Error while creating Pool")
	}

	if p.WorkerPool == nil {
		t.Error("Error while creating Pool")
	}

	if p.maxWorkers != 4 {
		t.Error("Error while creating Pool")
	}

	if p.Workers == nil {
		t.Error("Error while creating Pool")
	}

	testClosePool(t, &p)
}

// tries to run a single worker job pool
func TestRunSingleJob(t *testing.T) {

	// create pool
	p := NewWorkerPool(1)
	jobQueue := make(chan Job)
	p.Run(jobQueue)

	// check running pool num workers created
	if len(p.Workers) != 1 {
		t.Error("Wrong number of Workers")
	}

	// let's create a test job
	ret := NewJobResultChannel()
	//defer close(ret)
	m := make(map[string]string)
	task := NewTestTask()
	work := NewJob(&task, m, ret)

	// Push the job onto the queue.
	jobQueue <- work

	// wait for response from Job
	resp := <-ret

	if resp.Error != nil || resp.Value.(string) != "job done!" {
		t.Error("Error while getting result of Job")
	}

	// try to close pool
	testClosePool(t, &p)
}

// tries to run multiple jobs
func TestRunMultipleJobs(t *testing.T) {

	const workerLimit = 4

	// create pool
	p := NewWorkerPool(workerLimit)
	jobQueue := make(chan Job)
	p.Run(jobQueue)

	// check running pool num workers created
	if len(p.Workers) != workerLimit {
		t.Error("Wrong number of Workers")
	}

	// let's create a test job
	job1 := NewSumJobTask(1,2)

	// let's create a test job 2
	job2 := NewSumJobTask(3,4 )

	// let's create a test job 3
	job3 := NewSumJobTask(5,6 )

	// let's create a test job 4
	job4 := NewSumJobTask(7, 8)

	// Push each job onto the queue.
	jobQueue <- job1
	jobQueue <- job2
	jobQueue <- job3
	jobQueue <- job4

	// Consume the merged output from all jobs and check matching sum result
	sum := int64(0)
	for n := range Merge(job1.ReturnChannel, job2.ReturnChannel, job3.ReturnChannel, job4.ReturnChannel) {
		result := n.Value
		sum += result.(int64)
	}

	if sum != 36 {
		t.Error("error while summing all results from merged Jobs")
	}

	// try to close pool
	testClosePool(t, &p)
}

// tries to run workerLimit = 64000 multiple jobs
func TestRunThousandMultipleJob(t *testing.T) {

	const workerLimit = 64000

	// create pool
	p := NewWorkerPool(workerLimit)
	jobQueue := make(chan Job)
	p.Run(jobQueue)

	// check running pool num workers created
	if len(p.Workers) != workerLimit {
		t.Error("Wrong number of Workers")
	}

	// let's create the test jobs
	ret := make([]<-chan JobResult, 0)

	for i := 0; i < workerLimit; i++ {
		job := NewSumJobTask(5, 7)

		// Push the job onto the queue.
		jobQueue <- job
		ret = append(ret, job.ReturnChannel)
	}

	// Consume the merged output from all jobs and check matching sum result
	sum := int64(0)
	for n := range Merge(ret...) {
		result := n.Value
		sum += result.(int64)
	}

	if sum != (12 * workerLimit) {
		t.Error("error while summing all results from merged Jobs")
	}

	// try to close pool
	testClosePool(t, &p)
}

// tries to run multiple jobs SumRune
func TestRunMultipleSumRuneJobs(t *testing.T) {

	const workerLimit = 4

	// create pool
	p := NewWorkerPool(workerLimit)
	jobQueue := make(chan Job)
	p.Run(jobQueue)

	// check running pool num workers created
	if len(p.Workers) != workerLimit {
		t.Error("Wrong number of Workers")
	}

	// let's create a test job
	job1 := NewSumRunesOfWordJobTask("a new era of information is arising with guilhebl software solutions ... 123456789")

	// let's create a test job 2
	job2 := NewSumRunesOfWordJobTask(" a b c d e f g h i j k l m n o p q r s t u v x z ")

	// let's create a test job 3
	job3 := NewSumRunesOfWordJobTask(" test data !@#$%^" )

	// let's create a test job 4
	job4 := NewSumRunesOfWordJobTask(" !@#!%$#%!#%!#%!#%!_#)%*_@$)*%@#+%(#%{PQWIOR{PQIRAF:ALSKFMZX>V M>A< BDKJAS{PQOIRQOPWIRQ QWIEQ{WPI!@_$)(!_)%(*" )

	// Push each job onto the queue.
	jobQueue <- job1
	jobQueue <- job2
	jobQueue <- job3
	jobQueue <- job4

	// Consume the merged output from all jobs and check matching multiplication value
	sum := int64(0)
	for n := range Merge(job1.ReturnChannel, job2.ReturnChannel, job3.ReturnChannel, job4.ReturnChannel) {
		result := n.Value
		sum += result.(int64)
	}

	if sum != 17311 {
		t.Error("error while summing all results from merged Jobs")
	}

	// try to close pool
	testClosePool(t, &p)
}

// in order to create sample words
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSequenceString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// tries to run workerLimit = 64000 multiple jobs - SumRune
func TestRunThousandMultipleSumRuneJob(t *testing.T) {

	const workerLimit = 64000

	// create pool
	p := NewWorkerPool(workerLimit)
	jobQueue := make(chan Job)
	p.Run(jobQueue)

	// let's create the test jobs
	ret := make([]<-chan JobResult, 0)

	for i := 0; i < workerLimit; i++ {
		job := NewSumRunesOfWordJobTask("ABCDEFGHIJKLMNOPQRSTUVXYZ123456789 ABCDEFGHIJKLMNOPQRSTUVXYZ123456789 ABCDEFGHIJKLMNOPQRSTUVXYZ123456789")

		// Push the job onto the queue.
		jobQueue <- job
		ret = append(ret, job.ReturnChannel)
	}

	// Consume the merged output from all jobs and check matching sum result
	sum := int64(0)
	for n := range Merge(ret...) {
		result := n.Value
		sum += result.(int64)
	}

	if sum != (461760000) {
		t.Error("error while summing all results from merged SumRune Jobs")
	}

	// try to close pool
	testClosePool(t, &p)
}

// tries to run workerLimit = 100000 multiple jobs - SumRune
func TestRunHundredThousandMultipleSumRuneJob(t *testing.T) {

	const workerLimit = 100000

	// create pool
	p := NewWorkerPool(workerLimit)
	jobQueue := make(chan Job)
	p.Run(jobQueue)

	// let's create the test jobs
	ret := make([]<-chan JobResult, 0)

	for i := 0; i < workerLimit; i++ {
		job := NewSumRunesOfWordJobTask("ABCDEFGHIJKLMNOPQRSTUVXYZ123456789 ABCDEFGHIJKLMNOPQRSTUVXYZ123456789 ABCDEFGHIJKLMNOPQRSTUVXYZ123456789")

		// Push the job onto the queue.
		jobQueue <- job
		ret = append(ret, job.ReturnChannel)
	}

	// Consume the merged output from all jobs and check matching sum result
	sum := int64(0)
	for n := range Merge(ret...) {
		result := n.Value
		sum += result.(int64)
	}

	if sum != (721500000) {
		t.Error("error while summing all results from merged SumRune Jobs")
	}

	// try to close pool
	testClosePool(t, &p)
}

// tries to run workerLimit = 1000000 multiple jobs - SumRune
func TestRunMillionMultipleSumRuneJob(t *testing.T) {

	const workerLimit = 1000000

	// create pool
	p := NewWorkerPool(workerLimit)
	jobQueue := make(chan Job)
	p.Run(jobQueue)

	// let's create the test jobs
	ret := make([]<-chan JobResult, 0)

	for i := 0; i < workerLimit; i++ {
		job := NewSumRunesOfWordJobTask("ABCDEFGHIJKLMNOPQRSTUVXYZ123456789 ABCDEFGHIJKLMNOPQRSTUVXYZ123456789 ABCDEFGHIJKLMNOPQRSTUVXYZ123456789")

		// Push the job onto the queue.
		jobQueue <- job
		ret = append(ret, job.ReturnChannel)
	}

	// Consume the merged output from all jobs and check matching sum result
	sum := int64(0)
	for n := range Merge(ret...) {
		result := n.Value
		sum += result.(int64)
	}

	if sum != (7215000000) {
		t.Error("error while summing all results from merged SumRune Jobs")
	}

	// try to close pool
	testClosePool(t, &p)
}
