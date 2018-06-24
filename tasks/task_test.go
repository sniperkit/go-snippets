package tasks

import (
	"fmt"
	// "github.com/ipfs/go-datastore"
	// "testing"
)

type ExampleTask struct {
}

func NewExampleTask() Taskable {
	return &ExampleTask{}
}

func (e ExampleTask) Valid() error {
	return nil
}

func (e ExampleTask) Do(updates chan Progress) {
	updates <- Progress{
		Done: true,
	}
}

// TODO - finish
// func TestTaskStorage(t *testing.T) {
// 	// defer resetTestData(store, "tasks")

// 	store := datastore.NewMapDatastore()

// 	task := &Task{
// 		Title: "test",
// 		Type:  "example.task",
// 		Params: map[string]interface{}{
// 			"a": "b",
// 		},
// 	}
// 	if err := task.Save(store); err != nil {
// 		t.Error(err.Error())
// 		return
// 	}

// 	if err := task.Run(store); err != nil {
// 		t.Error(err.Error())
// 		return
// 	}
// 	if task.Request == nil {
// 		t.Errorf("task run didn't set request datestamp")
// 		return
// 	}

// 	if err := task.Cancel(store); err != nil {
// 		t.Error(err.Error())
// 		return
// 	}
// 	if task.Fail == nil {
// 		t.Errorf("task cancel didn't set fail datestamp")
// 		return
// 	}

// 	if err := task.Run(store); err != nil {
// 		t.Error(err.Error())
// 		return
// 	}
// 	if task.Request == nil {
// 		t.Errorf("task run didn't set request datestamp")
// 		return
// 	}

// 	if err := task.Errored(store, "test failure message"); err != nil {
// 		t.Error(err.Error())
// 		return
// 	}
// 	if task.Fail == nil {
// 		t.Errorf("task run didn't set fail datestamp")
// 		return
// 	}

// 	if err := task.Run(store); err != nil {
// 		t.Error(err.Error())
// 		return
// 	}
// 	if task.Request == nil {
// 		t.Errorf("task run didn't set request datestamp")
// 		return
// 	}

// 	if err := task.Succeeded(store, "test_success_url", "test_success_hash"); err != nil {
// 		t.Error(err.Error())
// 		return
// 	}
// 	if task.Success == nil {
// 		t.Errorf("task run didn't set success datestamp")
// 		return
// 	}

// 	task2 := &Task{Id: task.Id}
// 	if err := task2.Read(store); err != nil {
// 		t.Error(err.Error())
// 		return
// 	}

// 	if err := CompareTasks(task, task2); err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	if err := task.Delete(store); err != nil {
// 		t.Error(err.Error())
// 		return
// 	}
// }

func CompareTasks(a, b *Task) error {
	if a.Id != b.Id {
		return fmt.Errorf("Id mismatch: %s != %s", a.Id, b.Id)
	}
	if !a.Created.Equal(b.Created) {
		return fmt.Errorf("Created mismatch: %s != %s", a.Created, b.Created)
	}
	if !a.Updated.Equal(b.Updated) {
		return fmt.Errorf("Updated mismatch: %s != %s", a.Updated, b.Updated)
	}
	return nil
}
