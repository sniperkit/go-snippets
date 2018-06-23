package basic

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

func loadTOML() error {
	viper.SetConfigName("email") // name of config file (without extension)
	viper.AddConfigPath("./setting")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

func createMail(cfg map[string]interface{}) *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From", cfg["account"].(string))
	m.SetHeader("To", cfg["to"].(string))
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")
	// m.Attach("/home/Alex/lolcat.jpg")

	return m
}

func TestSendMail(t *testing.T) {
	err := loadTOML()
	if err != nil {
		fmt.Println(err)
	}
	// create mail
	cfg := viper.GetStringMap("GMAIL")
	m := createMail(cfg)

	// send mail
	server := cfg["server"].(string)
	port := int(cfg["port"].(int64))
	user := cfg["account"].(string)
	password := cfg["password"].(string)

	d := gomail.NewPlainDialer(server, port, user, password)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
