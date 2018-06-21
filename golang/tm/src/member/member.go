package member

import (
	"fmt"
	"task"
	"util"
)

type Member struct {
	Name     string
	ID       string
	Portrait string
	Tasks    []*task.Task
	Phone    string
}

func New(name string) *Member {
	return &Member{
		Name:     name,
		ID:       util.GenerateIDByString(name),
		Portrait: "",
		Tasks:    make([]*task.Task, 0, 1),
		Phone:    "",
	}
}

func (m *Member) SetID(id string) {
	m.ID = id
}

func (m *Member) SetPhoneNumber(phone string) {
	m.Phone = phone
}

func (m *Member) SetPortrait(link string) {
	m.Portrait = link
}

func (m *Member) Do(t *task.Task) {
	m.Tasks = append(m.Tasks, t)
}

func (m *Member) UnDo(t *task.Task) error {
	var index int
	for i, ts := range m.Tasks {
		if ts.ID == t.ID {
			index = i
			break
		}
	}

	if index < len(m.Tasks) {
		m.Tasks = append(m.Tasks[:index], m.Tasks[index+1:]...)
		return nil
	}

	return fmt.Errorf("Task %s is not assigned to: %s", t.Name, m.Name)
}

func (m *Member) String() string {
	p := fmt.Sprintf("Member: %20s, ID: %70s, Phone: %11s, Portrait: %20s\n", m.Name, m.ID, m.Phone, m.Portrait)
	p += fmt.Sprintf("    TaskList: \n")
	for _, t := range m.Tasks {
		p += fmt.Sprintf("            %s\n", t)
	}

	return p
}
