package basic

import (
	"log"
	"testing"
)

func TestInterfaceForNil(t *testing.T) {

	tFunc := func(val interface{}) {
		log.Println(val == nil)
		log.Println(val)
	}

	tFunc(nil)
	tFunc("hello")
	tFunc(12)
	tFunc(false)
}

func TestInterfaceSwitchWithBreakReturn(t *testing.T) {
	tFunc := func(val interface{}) {
		switch t := val.(type) {
		case int:
			log.Println("int ==> ", t)
			return
		case string:
			log.Println("string ==> ", t)
			return
		case int32:
			log.Println("int32 ==> ", t)
			break
		}

		log.Println("tFunc end")
	}

	tFunc(32)
}

type AA interface {
	Hello()
}

type BB interface {
	Hello()
}

type CC interface {
	World()
}

type CCImp struct{}

func (c CCImp) World() {
	log.Println("CC Imp")
}

type AAImp struct{}

func (a AAImp) Hello() {
	log.Println("AA Imp")
}

type BBImp struct{}

func (b BBImp) Hello() {
	log.Println("BB Imp")
}

func TestInterfaceCompatibility(t *testing.T) {
	var aa AA = AAImp{}
	var bb BB = BBImp{}
	var cc CC = CCImp{}
	aa = bb
	aa.Hello()
	bb.Hello()

	// following assigment will cause syntax error
	// Error: CC does not implement AA
	// aa = cc
	cc.World()
}
