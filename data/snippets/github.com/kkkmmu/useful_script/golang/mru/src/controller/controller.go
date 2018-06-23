package controller

import (
	"fmt"
	"log"
	"session"
)

type Controller struct {
	Sessions map[string]*session.Session
}

func New() *Controller {
	return &Controller{
		Sessions: make(map[string]*session.Session, 1),
	}
}

func (c *Controller) GetSessionByID(id string) *session.Session {
	if s, ok := c.Sessions[id]; ok {
		return s
	}

	return nil
}

func (c *Controller) GetSessionByUsernameAndPassword(name, pass string) *session.Session {
	id := session.GenerateSessionIDByUserNameAndPassword(name, pass)
	if s, ok := c.Sessions[id]; ok {
		return s
	}

	return nil
}

func (c *Controller) IsSessionExist(name, pass string) bool {
	id := session.GenerateSessionIDByUserNameAndPassword(name, pass)
	if _, ok := c.Sessions[id]; ok {
		return true
	}
	return false
}

func (c *Controller) AddSessionByUsernameAndPassword(name, pass string) (*session.Session, error) {
	id := session.GenerateSessionIDByUserNameAndPassword(name, pass)
	if _, ok := c.Sessions[id]; ok {
		return nil, fmt.Errorf("Session for user: %s already exist", name)
	}

	c.Sessions[id] = session.New(name, pass)

	log.Printf("Session DB: %#q", c.Sessions)
	log.Printf("Session DB: %#v", c.Sessions)

	return c.Sessions[id], nil
}

func (c *Controller) DelSessionByUsernameAndPassword(name, pass string) error {
	id := session.GenerateSessionIDByUserNameAndPassword(name, pass)
	if _, ok := c.Sessions[id]; ok {
		delete(c.Sessions, id)
		return nil
	}

	return fmt.Errorf("Session does not exist for user: %s", name)
}
