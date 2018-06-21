package product

import (
	"fmt"
	"task"
	"util"
)

type Product struct {
	Name  string
	ID    string
	Tasks []*task.Task
}

func New(name string) *Product {
	return &Product{
		Name:  name,
		ID:    util.GenerateIDByString(name),
		Tasks: make([]*task.Task, 0, 1),
	}
}

func (p *Product) AddTask(name string) (*task.Task, error) {
	newTask := task.New(name)
	p.Tasks = append(p.Tasks, newTask)

	return newTask, nil
}

func (p *Product) GetTaskByName(name string) (*task.Task, error) {
	for _, t := range p.Tasks {
		if t.Name == name {
			return t, nil
		}
	}

	return nil, fmt.Errorf("Cannot find task by name: %s!", name)
}

func (p *Product) GetTaskByID(id string) (*task.Task, error) {
	for _, t := range p.Tasks {
		if t.ID == id {
			return t, nil
		}
	}

	return nil, fmt.Errorf("Cannot find task by id: %s!", id)
}

func (p *Product) GetAllTasks() []*task.Task {
	return p.Tasks
}

func (p *Product) DelTaskByName(name string) error {
	var index int

	for i, t := range p.Tasks {
		if t.Name == name {
			index = i
			break
		}
	}

	if index < len(p.Tasks) {
		p.Tasks = append(p.Tasks[:index], p.Tasks[index+1:]...)
		return nil
	}

	return fmt.Errorf("Cannot find task with name: %s", name)
}

func (p *Product) DelTaskByID(id string) error {
	var index int

	for i, t := range p.Tasks {
		if t.ID == id {
			index = i
			break
		}
	}

	if index < len(p.Tasks) {
		p.Tasks = append(p.Tasks[:index], p.Tasks[index+1:]...)
		return nil
	}

	return fmt.Errorf("Cannot find task with id: %s", id)
}

func (p *Product) DelAllTasks() {
	p.Tasks = nil
}

func (p *Product) String() string {
	sp := fmt.Sprintf("Product: %19s, ID: %70s\n", p.Name, p.ID)
	for _, t := range p.Tasks {
		sp += fmt.Sprintf("%s\n", t)
	}

	return sp
}
