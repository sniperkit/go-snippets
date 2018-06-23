# GO Worker pool

Provides a generic GoRoutine worker pool for implementing parallel GoRoutine jobs.
Creates a number of worker routines which picks Jobs to be executed from a JobQueue in parallel.
The Workers executes the algorithm through an implementation of the Task interface and also accepts
a return channel to return task results to clients. Once the return value is sent back to the return channel the return channel output is closed.

Supports merging of multiple job outputs into a single output channel (check sample Sum Task test).

## How to use

1. Implement the task interface:

```
// Anonymous Sum Task
type TestSumTask struct{}

func (e *TestSumTask) Run(payload Payload) JobResult {
	//log.Printf("Task Sum %v", payload.Params)

	x, _ := strconv.ParseInt(payload.Params["x"], 10, 0)
	y, _ := strconv.ParseInt(payload.Params["y"], 10, 0)

	return NewJobResult(x + y, nil)
}

func NewTestSumTask() TestSumTask {
	return TestSumTask{}
}
```

2. Create Worker Pool and JobQueue (Worker Pool is created once per app, avoiding the risk of spanning multiple pools)

```
	// create pool
	p := NewWorkerPool(4)
	jobQueue := make(chan Job)
	p.Run(jobQueue)
```

3. Create Jobs

```
	// let's create a test job
	ret := NewJobResultChannel()
	m := make(map[string]string)
	x, y := 1, 2
	m["x"] = strconv.Itoa(x)
	m["y"] = strconv.Itoa(y)
	task := NewTestSumTask()
	work := NewJob(&task, m, ret)

	// ... create multiple tasks
```

4. Push jobs into the queue

```
    // push 1 single job
	jobQueue <- work

	// ... or push multiple jobs
	jobQueue <- work2
	jobQueue <- work3

```

5. *Optional If necessary merge multiple job outputs

```
	// Consume the merged output from all jobs
	sum := int64(0)
	for n := range Merge(ret, ret2, ret3, ret4) {
		result := n.Value
		sum += result.(int64)
	}
```

6. Once application is done using the pool Stop it.

```
	// try to close pool
	p.Stop()
```

#### Any contributions or suggestions are welcome!
