package main

import (
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/viper"
)

// Sample usage:
// VIPER_PATH=.env VIPER_FILE=yml VIPER_NAME=config go viperdemo.go
func main() {
	v := viper.New()
	v.AddConfigPath(os.Getenv("VIPER_PATH"))
	v.SetConfigFile(os.Getenv("VIPER_FILE"))
	v.SetConfigName(os.Getenv("VIPER_NAME"))
	v.ReadInConfig()

	spew.Dump(v.GetBool("disabled_modules.trello"))
}
