package redis

import (
	"github.com/fatih/color"
	"github.com/maxpowel/dislet"
    "github.com/go-redis/redis"

)

type Config struct {
	Uri string `default:"PUERTOOO"`
	Database int
}


func NewConnection(uri string, database int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     uri,
		Password: "", // no password set
		DB:       database,  // use default DB
	})

	return client
}


func Bootstrap(k *dislet.Kernel) {
	//fmt.Println("DATABASE BOOT")
	mapping := k.Config.Mapping
	mapping["redis"] = &Config{}

	var baz dislet.OnKernelReady = func(k *dislet.Kernel){
		color.Green("Booting redis")
		conf := k.Config.Mapping["redis"].(*Config)


		k.Container.RegisterType("redis", NewConnection, conf.Uri, conf.Database)
	}
	k.Subscribe(baz)

}
