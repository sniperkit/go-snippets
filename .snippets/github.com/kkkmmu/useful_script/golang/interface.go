package main

import (
	"log"
	"strings"
)

type Handler interface {
	Raw() string
	Handle(string) string
}

type Upper struct {
	raw string
}

func handle(h Handler) string {
	return h.Handle(h.Raw())
}

func (u *Upper) Handle(raw string) string {
	return strings.ToUpper(raw)
}

func (u *Upper) Raw() string {
	return u.raw
}

type Lower struct {
	raw string
}

func (l *Lower) Handle(raw string) string {
	return strings.ToLower(raw)
}

func (l *Lower) Raw() string {
	return l.raw
}

func main() {
	u := &Upper{
		raw: "Hello 123hll",
	}

	log.Println(handle(u))

	l := &Lower{
		raw: "Dasan Networks",
	}

	log.Println(handle(l))
}
