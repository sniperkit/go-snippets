package machinery

import (
	"github.com/fatih/color"
	"github.com/RichardKnop/machinery/v1"
	"fmt"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/maxpowel/dislet"

)

type Config struct {
	Broker string `default:"valor"`
	ResultBackend string `default:"PUERTOOO"`
	DefaultQueue string
}



/*func NewConnection(dialect string, uri string) *gorm.DB {
	db, _ := gorm.Open(dialect, uri)
	fmt.Println("CREANDO CONEXION")
	return db
}*/


func Add(args ...int64) (int64, error) {
	sum := int64(0)
	for _, arg := range args {
		sum += arg
	}
	return sum, nil
}

func Bootstrap(k *dislet.Kernel) {
	//fmt.Println("DATABASE BOOT")
	mapping := k.Config.Mapping
	mapping["machinery"] = &Config{}
	var baz dislet.OnKernelReady = func(k *dislet.Kernel){
		color.Green("Booting machinery")
		mConfig := k.Config.Mapping["machinery"].(*Config)
		//fmt.Println(mConfig)

		var cnf = config.Config{
			Broker: mConfig.Broker,
			ResultBackend: mConfig.ResultBackend,
			DefaultQueue: mConfig.DefaultQueue,
			//Broker : "redis://localhost:6379/0",
			//ResultBackend: "redis://localhost:6379/0",
			//Broker:             "amqp://guest:guest@localhost:5672/",
			//ResultBackend:      "amqp://guest:guest@localhost:5672/",
			//DefaultQueue: "machinery_tasks",
			/*AMQP:               &config.AMQPConfig{
				Exchange:     "machinery_exchange",
				ExchangeType: "direct",
				BindingKey:   "machinery_task",
			},*/
		}
		//fmt.Println(cnf)

		server, err := machinery.NewServer(&cnf)
		if err != nil {
			panic(err)
		}

		runWorker := func (server *machinery.Server) {
			worker := server.NewWorker("machinery_worker")
			if err := worker.Launch(); err != nil {
				fmt.Println("Error worker")
				fmt.Println(err)
			}
		}

		go runWorker(server)

		//k.Container.RegisterType("database", NewConnection, conf.Dialect, conf.Uri)
		iny := func() *machinery.Server{
			return server
		}
		k.Container.RegisterType("machinery", iny)


	}
	k.Subscribe(baz)




}