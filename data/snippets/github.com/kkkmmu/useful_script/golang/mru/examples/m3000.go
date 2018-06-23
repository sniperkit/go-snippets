package main

import (
	"command"
	"context"
	"fmt"
	"rut"
	"util"
)

func main() {
	dev, err := rut.New(&rut.RUT{
		Name:     "M3000_210",
		Device:   "V5",
		IP:       "10.55.192.210",
		Port:     "23",
		Username: "admin",
		Hostname: "M3000_210",
		Password: "",
	})

	if err != nil {
		panic(err)
	}

	dev.Init()

	ctx := context.Background()

	data, err := dev.RunCommand(ctx, &command.Command{
		Mode: "normal",
		CMD:  "configure terminal",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(data)

	data, err = dev.RunCommand(ctx, &command.Command{
		Mode: "config",
		CMD:  "show running-config",
	})
	fmt.Println(data)

	data, err = dev.RunCommand(ctx, &command.Command{
		Mode: "config",
		CMD:  "do q sh -l",
	})
	fmt.Println(data)

	data, err = dev.RunCommand(ctx, &command.Command{
		Mode: "shell",
		CMD:  "bcm techsupport cos",
	})
	fmt.Println(data)
	util.SaveToFile("cos.txt", []byte(data))

	data, err = dev.RunCommand(ctx, &command.Command{
		Mode: "shell",
		CMD:  "bcm cos bandwidth_show",
	})
	fmt.Println(data)
	util.SaveToFile("cos_b3.txt", []byte(data))
}
