package mailer

import (
	"github.com/fatih/color"
	"github.com/maxpowel/dislet"
	"gopkg.in/gomail.v2"

)

type Config struct {
	Hostname string `default:"valor"`
	Port int `default:"PUERTOOO"`
	Username string
	Password string
}

func NewMailer(hostname, username, password string, port int) *gomail.Dialer{
	d := gomail.NewDialer(hostname, port, username, password)
	//d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return d
}

func Bootstrap(k *dislet.Kernel) {
	mapping := k.Config.Mapping
	mapping["mailer"] = &Config{}

	var baz dislet.OnKernelReady = func(k *dislet.Kernel){
		color.Green("Booting mailer")
		conf := k.Config.Mapping["mailer"].(*Config)
		k.Container.RegisterType("mailer", NewMailer, conf.Hostname, conf.Username, conf.Password, conf.Port)

	}
	k.Subscribe(baz)
}
