package tasks

import (
	"fmt"
	"github.com/ipfs/go-datastore"
)

// taskdefs is an internal registry of all types of tasks that can be performed.
// in order for a task to be managed, it must first be added by calling RegisterTaskdef
// with a function that produces new instances of taskable for marshalling params
var taskdefs = map[string]NewTaskFunc{}

// RegisterTaskdef registers a task type, must be called before a task can be used.
func RegisterTaskdef(name string, f NewTaskFunc) {
	taskdefs[name] = f
}

// NewTaskable generates a new Taskable instance from the registered
// types
func NewTaskable(name string) (Taskable, error) {
	if taskdefs[name] == nil {
		return nil, fmt.Errorf("unknown task type: %s", name)
	}

	return taskdefs[name](), nil
}

// Taskable anything that fits on a task queue, it is a type of "work"
// that can be performed. Lots of things
type Taskable interface {
	// are these task params valid? return error if not
	// this func will be called before adding the task to
	// the queue, and won't be added on failure.
	Valid() error
	// Do the task, returning incremental progress updates
	// it's expected that the func will send either
	// p.Done == true or p.Error != nil at least once
	// to signal that the task is either done or errored
	Do(updates chan Progress)
}

// NewTaskFunc is a function that creates new task instances
// task-orchestrators use NewTaskFunc to create new Tasks, and then attempt
// to json.Unmarshal params into the task definition
type NewTaskFunc func() Taskable

// SqlDbTaskable is a task that has a method for assigning a datastore to the task.
// If your task needs access to a datastore, implement DatastoreTaskable, task-orchestrators
// will detect this method and call it to set the datastore before calling Taskable.Do
type DatastoreTaskable interface {
	Taskable
	SetDatastore(ds datastore.Datastore)
}
