package mongo

import (
	"github.com/fatih/color"
	"github.com/maxpowel/dislet"
	"gopkg.in/mgo.v2"
)

type Config struct {
	Uri string `default:"PUERTOOO"`
	Database string
}


func NewConnection(uri string, database string) *mgo.Database {

	session, err := mgo.Dial(uri)
	if err != nil {
		panic(err)
	}
	db := session.DB(database)
	return db

}

/*type Person struct {
	Name string
	Phone string
}*/

func Bootstrap(k *dislet.Kernel) {
	mapping := k.Config.Mapping
	mapping["mongo"] = &Config{}

	var baz dislet.OnKernelReady = func(k *dislet.Kernel){
		color.Green("Booting mongo")
		conf := k.Config.Mapping["mongo"].(*Config)
		k.Container.RegisterType("mongo", NewConnection, conf.Uri, conf.Database)


		/*db := k.Container.MustGet("mongo").(*mgo.Database)
		col := db.C("nombres")
		err := col.Insert(&Person{"Ale", "+55 53 8116 9639"})
		if err != nil {
			fmt.Println(err)
		}*/
		//https://labix.org/mgo
	}
	k.Subscribe(baz)

}
