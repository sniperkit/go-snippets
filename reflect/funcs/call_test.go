package funcs

import (
	"log"
	"testing"
)

func SayHi() {
	log.Println("say hi")
}

func Sum(i, j int) int {
	return i + j
}

func TestFunctionsDirector(t *testing.T) {
	fd := NewFuncsDirector()

	if err := fd.Registe(SayHi); err != nil {
		t.Error("function director registe SayHi should return nil")
	}

	if err := fd.Registe(SayHi); err == nil {
		t.Error("function director registe SayHi should return a error ErrorFuncDirectorAlreadyRegisted")
	}

	fd.Registe(Sum)

	if err := fd.Registe(3); err == nil {
		t.Errorf("function director registe type is wrong should not return nil")
	}

	// use full name spcify the package function was contained
	_, err := fd.Call("github.com/EricYT/go-examples/reflect/funcs.SayHi")
	if err != nil {
		t.Errorf("call function SayHi should return nil but got %s", err)
	}

	_, err = fd.Call("github.com/EricYT/go-examples/reflect/funcs.Sum", 1, 2, 4)
	switch err {
	case ErrorFuncDirectorParamsNotMatch:
	default:
		t.Errorf("call function Sum with arguments 1, 2, 4 should return ErrorFuncDirectorParamsNotMatch")
	}

	res, err := fd.Call("github.com/EricYT/go-examples/reflect/funcs.Sum", 1, 2)
	if err != nil {
		t.Errorf("call function Sum should return nil but got: %s", err)
	}
	if len(res) != 1 {
		t.Errorf("call function Sum should return just one argument but got %d", len(res))
	}
	value := int(res[0].Int())
	if value != 3 {
		t.Errorf("call function Sum 1 + 2 should equal 3 but got %d", value)
	}
}
