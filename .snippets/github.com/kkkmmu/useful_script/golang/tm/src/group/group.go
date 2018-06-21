package group

import (
	"fmt"
	"member"
	"product"
	"task"
	"team"
	"util"
)

type Group struct {
	Name     string
	ID       string
	Teams    []*team.Team
	Tasks    []*task.Task
	Members  []*member.Member
	Products []*product.Product
}

func New(name string) *Group {
	return &Group{
		Name:     name,
		ID:       util.GenerateIDByString(name),
		Teams:    make([]*team.Team, 0, 1),
		Tasks:    make([]*task.Task, 0, 1),
		Members:  make([]*member.Member, 0, 1),
		Products: make([]*product.Product, 0, 1),
	}
}

func (g *Group) AddTeam(name string) (*team.Team, error) {
	if g.IsTeamExist(name) {
		return nil, fmt.Errorf("Team with name: %s already exist!", name)
	}

	newTeam := team.New(name)
	g.Teams = append(g.Teams, newTeam)

	return newTeam, nil
}

func (g *Group) IsTeamExist(name string) bool {
	for _, t := range g.Teams {
		if t.Name == name {
			return true
		}
	}

	return false
}

func (g *Group) GetTeamByName(name string) (*team.Team, error) {
	for _, t := range g.Teams {
		if t.Name == name {
			return t, nil
		}
	}

	return nil, fmt.Errorf("Cannot find team by name: %s", name)
}

func (g *Group) GetTeamByID(id string) (*team.Team, error) {
	for _, t := range g.Teams {
		if t.ID == id {
			return t, nil
		}
	}

	return nil, fmt.Errorf("Cannot find team by id: %s", id)
}

func (g *Group) GetAllTeams() []*team.Team {
	return g.Teams
}

func (g *Group) DelTeamByName(name string) error {
	var index int
	for i, t := range g.Teams {
		if t.Name == name {
			index = i
			break
		}
	}

	if index < len(g.Teams) {
		g.Teams = append(g.Teams[0:index], g.Teams[index+1:]...)
		return nil
	} else {

		return fmt.Errorf("Cannot Delete Team %s with error: Not exist", name)
	}
}

func (g *Group) DelTeamByID(id string) error {
	var index int
	for i, t := range g.Teams {
		if t.ID == id {
			index = i
			break
		}
	}

	if index < len(g.Teams) {
		g.Teams = append(g.Teams[0:index], g.Teams[index+1:]...)
		return nil
	} else {

		return fmt.Errorf("Cannot Delete Team %s with error: Not exist", id)
	}
}

func (g *Group) DelTeam(t *team.Team) error {
	return g.DelTeamByName(t.Name)
}

func (g *Group) GetTaskByName(name string) (*task.Task, error) {
	for _, t := range g.Tasks {
		if t.Name == name {
			return t, nil
		}
	}

	return nil, fmt.Errorf("Cannot find task with name: %s", name)
}

func (g *Group) GetTaskByID(id string) (*task.Task, error) {
	for _, t := range g.Tasks {
		if t.ID == id {
			return t, nil
		}
	}

	return nil, fmt.Errorf("Cannot find task with id: %s", id)
}

func (g *Group) GetAllTasks() []*task.Task {
	return g.Tasks
}

func (g *Group) DelTaskByName(name string) error {
	var index int

	for i, t := range g.Tasks {
		if t.Name == name {
			index = i
			break
		}
	}

	if index < len(g.Tasks) {
		g.Tasks = append(g.Tasks[:index], g.Tasks[index+1:]...)
		return nil
	}

	return fmt.Errorf("Cannot find task with name: %s", name)
}

func (g *Group) DelTaskByID(id string) error {
	var index int

	for i, t := range g.Tasks {
		if t.ID == id {
			index = i
			break
		}
	}

	if index < len(g.Tasks) {
		g.Tasks = append(g.Tasks[:index], g.Tasks[index+1:]...)
		return nil
	}

	return fmt.Errorf("Cannot find task with id: %s", id)
}

func (g *Group) IsMemberExist(name string) bool {
	for _, m := range g.Members {
		if m.Name == name {
			return true
		}
	}

	return false
}

func (g *Group) AddMember(name string) (*member.Member, error) {
	if g.IsMemberExist(name) {
		return nil, fmt.Errorf("Same member: %s already exist", name)
	}

	newMember := member.New(name)
	g.Members = append(g.Members, newMember)

	return newMember, nil
}

func (g *Group) GetMemberByName(name string) (*member.Member, error) {
	for _, m := range g.Members {
		if m.Name == name {
			return m, nil
		}
	}

	return nil, fmt.Errorf("Cannot find member with name: %s", name)
}

func (g *Group) GetMemberByID(id string) (*member.Member, error) {
	for _, m := range g.Members {
		if m.ID == id {
			return m, nil
		}
	}

	return nil, fmt.Errorf("Cannot find member with id: %s", id)
}

func (g *Group) GetAllMembers() []*member.Member {
	return g.Members
}

func (g *Group) DelMemberByName(name string) error {
	var index int
	for i, m := range g.Members {
		if m.Name == name {
			index = i
			break
		}
	}

	if index < len(g.Members) {
		g.Tasks = append(g.Tasks[:index], g.Tasks[index+1:]...)
		return nil
	}

	return fmt.Errorf("Cannot find task with name: %s", name)
}

func (g *Group) DelMemberByID(id string) error {
	var index int
	for i, m := range g.Members {
		if m.ID == id {
			index = i
			break
		}
	}

	if index < len(g.Members) {
		g.Tasks = append(g.Tasks[:index], g.Tasks[index+1:]...)
		return nil
	}

	return fmt.Errorf("Cannot find task with id: %s", id)
}

func (g *Group) IsProductExist(name string) bool {
	for _, p := range g.Products {
		if p.Name == name {
			return true
		}
	}

	return false
}

func (g *Group) AddProduct(name string) (*product.Product, error) {
	if g.IsProductExist(name) {
		return nil, fmt.Errorf("Product %s alread exist!", name)
	}

	newProduct := product.New(name)

	g.Products = append(g.Products, newProduct)

	return newProduct, nil
}

func (g *Group) GetProductByName(name string) (*product.Product, error) {
	for _, p := range g.Products {
		if p.Name == name {
			return p, nil
		}
	}

	return nil, fmt.Errorf("Cannot find product by name: %s", name)
}

func (g *Group) GetProductByID(id string) (*product.Product, error) {
	for _, p := range g.Products {
		if p.ID == id {
			return p, nil
		}
	}

	return nil, fmt.Errorf("Cannot find product with id: %s", id)
}

func (g *Group) GetAllProducts() []*product.Product {
	return g.Products
}

func (g *Group) DelProductByName(name string) error {
	var index int

	for i, p := range g.Products {
		if p.Name == name {
			index = i
			break
		}
	}

	if index < len(g.Products) {
		g.Products = append(g.Products[:index], g.Products[index+1:]...)
		return nil
	}

	return fmt.Errorf("Cannot find product by name: %s", name)
}

func (g *Group) DelProductByID(id string) error {
	var index int

	for i, p := range g.Products {
		if p.ID == id {
			index = i
			break
		}
	}

	if index < len(g.Products) {
		g.Products = append(g.Products[:index], g.Products[index+1:]...)
		return nil
	}

	return fmt.Errorf("Cannot find product by id: %s", id)
}

func (g *Group) Reset() {
	g.Teams = make([]*team.Team, 0, 1)
	g.Tasks = make([]*task.Task, 0, 1)
	g.Members = make([]*member.Member, 0, 1)
	g.Products = make([]*product.Product, 0, 1)
}
