package main

import (
	"command"
	"fmt"
	"rut"
)

func main() {
	dev, err := rut.New(&rut.RUT{
		Name:     "DUT1",
		Device:   "V8",
		IP:       "10.71.20.198",
		Port:     "23",
		Username: "admin",
		Hostname: "SWITCH[A]",
		Password: "",
	})

	if err != nil {
		panic(err)
	}

	dev.Init()

	data, err := dev.RunCommand(&command.Command{
		Mode: "normal",
		CMD:  "configure terminal",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(data)

	data, err = dev.RunCommand(&command.Command{
		Mode: "config",
		CMD:  "show running-config",
	})
}
