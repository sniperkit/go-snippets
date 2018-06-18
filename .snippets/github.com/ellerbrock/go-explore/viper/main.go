package main

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func main() {
	setDefaults()
	setConfig()
	watchConfig()
	confEnv()
	setEnv()

	// get the default values
	fmt.Println("DEV Folder:", viper.Get("devDir"))
	fmt.Println("PROD Folder:", viper.Get("prodDir"))

	// get value from config
	fmt.Println("Firstname:", viper.Get("firstName"))
	fmt.Println("Lastname:", viper.Get("lastName"))
	fmt.Println("Age:", viper.Get("age"))

	// get environment variable
	fmt.Println("Repository:", getEnv("repo"))
}

func setDefaults() {
	viper.SetDefault("devDir", "dev")
	viper.SetDefault("prodDir", "prod")
}

func setConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	checkError(err)
}

func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
}

func confEnv() {
	viper.SetEnvPrefix("myapp")
	viper.BindEnv("repo")
	viper.BindEnv("type")
}

func setEnv() {
	os.Setenv("MYAPP_REPO", "https://github.com/ellerbrock")
	os.Setenv("MYAPP_TYPE", "git")
}

func getEnv(ev string) string {
	return viper.GetString(ev)
}

func checkError(err error) {
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
