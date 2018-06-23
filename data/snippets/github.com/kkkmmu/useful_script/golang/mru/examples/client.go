package main

import (
	"command"
	"context"
	"fmt"
	"rut"
)

var CTX = context.Background()

func main() {

	dev, err := rut.New(&rut.RUT{
		Name:     "SWITCH",
		Device:   "V8",
		IP:       "10.42.60.218",
		Port:     "22",
		Username: "admin",
		Protocol: "telnet",
		Hostname: "V8500",
		SFU:      "A",
		Password: "krcho",
	})

	if err != nil {
		panic(err)
	}

	err = dev.Init()
	if err != nil {
		panic(err)
	}

	_, err = dev.RunCommand(CTX, &command.Command{
		Mode: "normal",
		CMD:  " config terminal",
	})
	if err != nil {
		fmt.Println(err)
	}

	data, err := dev.RunCommand(CTX, &command.Command{
		Mode: "config",
		CMD:  " show running-config",
	})

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))

	/*
			data, err = dev.RunCommand(CTX, &command.Command{
				Mode: "config",
				CMD:  " policy mulitcast_225 modify",
			})

			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(data))

			data, err = dev.RunCommand(CTX, &command.Command{
				Mode: "config-policy",
				CMD:  " include-flow multicast_224",
			})

			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(data))

			data, err = dev.RunCommand(CTX, &command.Command{
				Mode: "config-policy",
				CMD:  " include-flow multicast_225",
			})

			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(data))

			data, err = dev.RunCommand(CTX, &command.Command{
				Mode: "config-policy",
				CMD:  "apply",
			})

			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(data))

			data, err = dev.RunCommand(CTX, &command.Command{
				Mode: "config-policy",
				CMD:  " include-flow multicast_226",
			})

			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(data))

		data, err = dev.RunCommand(CTX, &command.Command{
			Mode: "config",
			CMD:  "do q sh -l",
		})

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(data))

		data, err = dev.RunCommand(CTX, &command.Command{
			Mode: "shell",
			CMD:  "scontrol -f /proc/switch/ASIC/ctrl dump table 0 fp_tcam 0 1",
		})

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(data))

	*/
}
