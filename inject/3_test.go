package test

import (
	"testing"
	"github.com/codegangsta/inject"
	"log"
)

func getInjector3() inject.Injector {
	injector := inject.New()
	injector.Map(&AdminDatabaseClient{"Admin"})
	injector.Map(&AnonDatabaseClient{"Anonymous"})
	return injector
}

func useAdmin(db *AdminDatabaseClient) {
	log.Printf("[invoke] Admin: %s\n", db.Name)
}
func useAnon(db *AnonDatabaseClient) {
	log.Printf("[invoke] Anon: %s\n", db.Name)
}

func Test3(t *testing.T) {
	injector := getInjector3()
	injector.Invoke(useBoth)
	injector.Invoke(useAnon)
	injector.Invoke(useAdmin)
}