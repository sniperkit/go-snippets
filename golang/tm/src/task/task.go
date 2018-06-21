package task

import (
	"fmt"
	"util"
)

const (
	TaskOpen = iota
	TaskAssigned
	TaskScheduled
	TaskFinished
	TaskClosed
)

const (
	TaskPriorityHigh = iota
	TaskPriorityMiddle
	TaskPriorityLow
)

type Task struct {
	ID       string
	Name     string
	Ditto    string
	State    int
	Priority int
}

var TaskStateToStr map[int]string = map[int]string{
	TaskOpen:      "Open",
	TaskAssigned:  "Assigned",
	TaskScheduled: "Scheduled",
	TaskFinished:  "Finished",
	TaskClosed:    "Closed",
}

var TaskPriorityToStr map[int]string = map[int]string{
	TaskPriorityHigh:   "High",
	TaskPriorityMiddle: "Middle",
	TaskPriorityLow:    "Low",
}

func New(name string) *Task {
	return &Task{
		Name:     name,
		ID:       util.GenerateIDByString(name),
		Ditto:    "",
		State:    TaskOpen,
		Priority: TaskPriorityMiddle,
	}
}

func (t *Task) SetState(new int) {
	t.State = new
}

func (t *Task) GetState() int {
	return t.State
}

func (t *Task) SetName(new string) error {
	t.Name = new
	return nil
}

func (t *Task) GetName() string {
	return t.Name
}

func (t *Task) GetDittoID() string {
	return t.Ditto
}

func (t *Task) SetDittoID(new string) error {
	t.ID = new
	return nil
}

func (t *Task) SetPriority(new int) {
	t.Priority = new
}

func (t *Task) String() string {
	p := fmt.Sprintf("Task: %22s, ID %70s, State: %s, Priority: %s", t.Name, t.ID, TaskStateToStr[t.State], TaskPriorityToStr[t.Priority])
	return p
}
