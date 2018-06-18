package test

import (
	"testing"
	"github.com/codegangsta/inject"
	"fmt"
	"github.com/e154/smart-home/api/log"
)

func getInjector() inject.Injector {
	injector := inject.New()
	injector.Map(&DatabaseClient{"Hello from DatabaseClient"})
	injector.Map(&RedisClient{"Hello from RedisClient"})
	return injector
}

func getContainer(container interface{}) error {
	injector := getInjector()
	return injector.Apply(container)
}

func Test1(t *testing.T) {

	c := DependencyContainer{}

	err := getContainer(&c)
	if err != nil {
		log.Error(err.Error())
		return
	}

	fmt.Printf("Database: %s\n", c.Db.Name)
	fmt.Printf("Redis: %s\n", c.Redis.Name)
}
