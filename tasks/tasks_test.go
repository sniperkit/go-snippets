package tasks

import (
	"github.com/ipfs/go-datastore"
	"testing"
)

func TestReadTasks(t *testing.T) {
	egs := []*Task{
		&Task{Title: "a", Type: "test"},
		&Task{Title: "b", Type: "test"},
		&Task{Title: "c", Type: "test"},
	}

	RegisterTaskdef("test", NewExampleTask)

	store := datastore.NewMapDatastore()
	for _, tsk := range egs {
		if err := tsk.Save(store); err != nil {
			t.Error(err.Error())
			return
		}
	}

	tasks, err := ReadTasks(store, "created DESC", 10, 0)
	if err != nil {
		t.Error(err)
		return
	}
	if len(tasks) != len(egs) {
		t.Errorf("task length mismatch: %d != %d", len(tasks), len(egs))
	}
}

// TODO - re-enable
// func TestGenerateAvailableTasks(t *testing.T) {
// 	defer resetTestData(appDB, "tasks")
// 	tasks, err := GenerateAvailableTasks(appDB)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	if len(tasks) == 0 {
// 		t.Error("no tasks returned")
// 	}
// }
