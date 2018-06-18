package test

import (
	"testing"
	"github.com/tsaikd/inject"
)

func (r ObjectFactory) NewAdminDatabaseClient() *AdminDatabaseClient {
	return &AdminDatabaseClient{"Administrator"}
}
func (r ObjectFactory) NewAnonDatabaseClient() *AnonDatabaseClient {
	return &AnonDatabaseClient{"Anonymous"}
}

func getInjector5() inject.Injector {
	of := ObjectFactory{}
	injector := inject.New()
	injector.Provide(of.NewAdminDatabaseClient)
	injector.Provide(of.NewAnonDatabaseClient)
	return injector
}

func Test5(t *testing.T) {
	injector := getInjector5()
	injector.Provide(nil)
}