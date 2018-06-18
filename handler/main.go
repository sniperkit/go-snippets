package main

import (
	"fmt"
	"errors"
)

const (
	BOOL int = iota
	INT64
	FLOAT64
	STRING
)

type variable struct {
	typ_e int
	b *bool
	i *int64
	f *float64
	s *string
}

type server struct {
	h map[string]func(key, value string)
	variables map[string]*variable
	b *bool
}

func (s *server) Register(name string, h func(key, value string)) {
	if s.h[name] != nil {
		delete(s.h, name)
	}
	s.h[name] = h
}

func (s *server) AddBool(name string, b *bool) {
	s.b = b
}

func (s *server) SetBool(name string, value bool) {
	*s.b = value
}

func (s *server) GetBool(name string) bool {
	return *s.b
}

func (s *server) AddString(name string, str *string) error {
	if s.variables[name] != nil {
		return errors.New("variable name is used")
	}

	v := new(variable)
	v.typ_e = STRING
	v.s = str
	s.variables[name] = v

	return nil
}

func (s *server) GetString(name string) (string, error) {
	if s.variables[name] == nil {
		return "", errors.New("variable not found")
	}

	return *s.variables[name].s, nil
}

func (s *server) SetVariable(name string, value interface {}) error {
	if s.variables[name] == nil {
		return errors.New("variable not found")
	}

	switch s.variables[name].typ_e{
	case BOOL:
		*s.variables[name].b, _ = value.(bool)
	case INT64:
		*s.variables[name].i, _ = value.(int64)
	case FLOAT64:
		*s.variables[name].f, _ = value.(float64)
	case STRING:
		*s.variables[name].s, _ = value.(string)
	default:
		break
	}

	return nil
}

func (s *server) Run(name, key, value string) {
	if s.h[name] != nil {
		s.h[name](key, value)
	}
}

var Server = &server{
	h: make(map[string]func(key, value string)),
	variables: make(map[string]*variable),
}

func client(key, value string) {
	fmt.Printf("%s %s\n", key, value)
}

func client2(key, value string) {
	fmt.Printf("%s %s\n", key, value)
}

var (
	qwe bool = false
	qwe2 string = "qwe2"
)

func main(){
	Server.Register("client", client)
	Server.Register("client2", client2)

	Server.Register("client3", func(key, value string) {
		fmt.Printf("%s %s\n", key, value)
	})
	Server.Run("client", "key", "value")
	Server.Run("client2", "key1", "value2")
	Server.Run("client3", "key1", "value2")

	Server.AddBool("qwe", &qwe)
	Server.SetBool("qwe", false)

	fmt.Println(qwe)

	Server.AddString("qwe2", &qwe2)
	Server.SetVariable("qwe2", "qwe3")
	fmt.Println(qwe2)
}
