package conf

import (
	"log"

	"github.com/jinzhu/configor"
)

// Config contains the project configuration.
var Config = struct {
	AppName string `default:"demo"`
	Profile string `default:"dev"`
	DB      struct {
		Host string `default:"mysql"`
		User string `default:"root"`
	}
	KafkaBroker []string
}{}

func init() {
	err := configor.Load(&Config, "app.conf")
	if err != nil {
		log.Print(err)
	}
}
