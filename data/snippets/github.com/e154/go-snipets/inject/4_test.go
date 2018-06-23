package test

import (
	"testing"
	"fmt"
	"github.com/codegangsta/inject"
)


func (r *ObjectFactory) GetDatabase(name string) (*DatabaseClient, error) {
	switch name {
	case "admin":
		return &DatabaseClient{"Administrator"}, nil
	case "anon":
		return &DatabaseClient{"Anonymous"}, nil
	}
	return nil, fmt.Errorf("Unknown database definition: '%s'", name)
}

func getInjector4() inject.Injector {
	injector := inject.New()
	injector.Map(&ObjectFactory{})
	return injector
}

func useObjectFactory(of *ObjectFactory) {
	names := []string{"admin", "anon", "missing"}
	for _, name := range names {
		db, err := of.GetDatabase(name)
		if err != nil {
			fmt.Printf("[db] got error: %s\n", err)
			continue
		}
		fmt.Printf("[db] got database '%s'\n", db.Name)
	}
}

func Test4(t *testing.T) {
	injector := getInjector4()
	injector.Invoke(useObjectFactory)
}