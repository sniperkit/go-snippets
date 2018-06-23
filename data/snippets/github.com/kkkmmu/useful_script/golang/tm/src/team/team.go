package team

import (
	"fmt"
	"member"
	"product"
	"task"
	"util"
)

type Team struct {
	Name     string
	ID       string
	Members  []*member.Member
	Tasks    []*task.Task
	Products []*product.Product
}

func New(name string) *Team {
	return &Team{
		Name:     name,
		ID:       util.GenerateIDByString(name),
		Members:  make([]*member.Member, 0, 1),
		Tasks:    make([]*task.Task, 0, 1),
		Products: make([]*product.Product, 0, 1),
	}
}

func (t *Team) Join(m *member.Member) error {
	if t.IsMemberExist(m) {
		return fmt.Errorf("%s is already a member of team: %s", m.Name, t.Name)
	}

	t.Members = append(t.Members, m)

	return nil
}

func (t *Team) IsMemberExist(m *member.Member) bool {
	for _, mem := range t.Members {
		if mem.ID == m.ID {
			return true
		}
	}

	return false
}

func (t *Team) Leave(m *member.Member) error {
	var index int

	for i, mem := range t.Members {
		if mem.ID == m.ID {
			index = i
			break
		}
	}

	if index < len(t.Members) {
		t.Members = append(t.Members[:index], t.Members[index+1:]...)
		return nil
	}

	return fmt.Errorf("%s is not a member of team: %s", m.Name, t.Name)
}

func (t *Team) GetAllMembers() []*member.Member {
	return t.Members
}

func (t *Team) IsProductExist(p *product.Product) bool {
	for _, prod := range t.Products {
		if prod.ID == p.ID {
			return true
		}
	}

	return false
}

func (t *Team) Maintaince(p *product.Product) error {
	if t.IsProductExist(p) {
		return fmt.Errorf("Product %s already maintained by %s", p.Name, t.Name)
	}

	t.Products = append(t.Products, p)

	return nil
}

func (t *Team) UnMaintaince(p *product.Product) error {

	var index int

	for i, prod := range t.Products {
		if prod.ID == p.ID {
			index = i
			break
		}
	}

	if index < len(t.Products) {
		t.Products = append(t.Products[:index], t.Products[index+1:]...)
		return nil
	}

	return fmt.Errorf("Product %s is not maintained by %s", p.Name, t.Name)
}

func (t *Team) GetAllProducts() []*product.Product {
	return t.Products
}

func (t *Team) GetAllTasks() []*task.Task {
	return t.Tasks
}

func (t *Team) String() string {
	p := fmt.Sprintf("=========================================================\n")
	p += fmt.Sprintf("Team: %10s,  ID: %20s\n", t.Name, t.ID)

	p += fmt.Sprintf("Member List(%d): \n", len(t.Members))
	for _, m := range t.Members {
		p += fmt.Sprintf("             %10s\n", m.Name)
	}

	p += fmt.Sprintf("Product List(%d): \n", len(t.Products))
	for _, prod := range t.Products {
		p += fmt.Sprintf("             %10s\n", prod.Name)
	}

	p += fmt.Sprintf("Task List(%d): \n", len(t.Tasks))
	for _, ts := range t.Tasks {
		p += fmt.Sprintf("             %10s\n", ts.Name)
	}

	return p
}
