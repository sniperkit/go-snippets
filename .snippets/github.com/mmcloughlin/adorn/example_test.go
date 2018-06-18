package adorn_test

import (
	"fmt"
	"strings"

	"github.com/mmcloughlin/adorn"
)

func ExampleConfig_FuncTypeName() {
	cfg := adorn.Config{
		TypeName: "Greeter",
	}
	fmt.Println(cfg.FuncTypeName())
	// Output: GreeterFunc
}

func ExampleConfig_ArgumentsUnnamed() {
	cfg := adorn.Config{
		ArgumentTypes: []string{"int", "int", "string", "string", "string"},
	}
	fmt.Println(cfg.ArgumentsUnnamed())
	// Output: int, int, string, string, string
}

func ExampleConfig_ArgumentNames() {
	cfg := adorn.Config{
		ArgumentTypes: []string{"int", "string", "string"},
	}
	fmt.Println(cfg.ArgumentNames())
	// Output: [a0 a1 a2]
}

func ExampleConfig_ArgumentTypesDeduped() {
	cfg := adorn.Config{
		ArgumentTypes: []string{"int", "int", "string", "string", "string"},
	}
	fmt.Println(strings.Join(cfg.ArgumentTypesDeduped(), ","))
	// Output: ,int,,,string
}

func ExampleConfig_ReturnSignature_single() {
	cfg := adorn.Config{
		ReturnTypes: []string{"int"},
	}
	fmt.Println(cfg.ReturnSignature())
	// Output: int
}

func ExampleConfig_ReturnSignature_multiple() {
	cfg := adorn.Config{
		ReturnTypes: []string{"string", "error"},
	}
	fmt.Println(cfg.ReturnSignature())
	// Output: (string, error)
}
