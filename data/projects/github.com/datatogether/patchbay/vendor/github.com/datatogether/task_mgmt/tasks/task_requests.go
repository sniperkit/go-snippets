package tasks

import (
	"github.com/ipfs/go-datastore"
)

// TaskRequests encapsulates all types of requests that can be made
// in relation to tasks, to be made available for RPC calls.
// TODO - should this internal state be moved into the package level
// via package-level setter funcs?
type TaskRequests struct {
	// url to amqp server for enqueuing tasks, only required
	// to fullfill requests, not submit them
	AmqpUrl string
	// Store to read / write tasks to only required
	// to fulfill requests, not submit them
	Store datastore.Datastore
}

// TasksEnqueueParams are for enqueing a task.
type TasksEnqueueParams struct {
	// Title of the task
	// Requesters should generate their own task title for now
	// tasks currently have no way of generating a sensible default title
	Title string
	// Type of task to perform
	Type string
	// User that initiated the request
	UserId string
	// Parameters to feed to the task
	Params map[string]interface{}
}

// Add a task to the queue for completion
func (r TaskRequests) Enqueue(params *TasksEnqueueParams, task *Task) (err error) {
	t := &Task{
		Title:  params.Title,
		Type:   params.Type,
		UserId: params.UserId,
		Params: params.Params,
	}

	if err := t.Enqueue(r.Store, r.AmqpUrl); err != nil {
		return err
	}

	*task = *t
	return nil
}

// Get a single Task, currently only lookup by ID is supported
type TasksGetParams struct {
	Id string
}

func (t TaskRequests) Get(args *TasksGetParams, res *Task) (err error) {
	tsk := &Task{
		Id: args.Id,
	}
	err = tsk.Read(t.Store)
	if err != nil {
		return err
	}

	*res = *tsk
	return nil
}

type TasksListParams struct {
	OrderBy string
	Limit   int
	Offset  int
}

func (t TaskRequests) List(args *TasksListParams, res *[]*Task) (err error) {
	ts, err := ReadTasks(t.Store, args.OrderBy, args.Limit, args.Offset)
	if err != nil {
		return err
	}
	*res = ts
	return nil
}
