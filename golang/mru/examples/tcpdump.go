package main

import (
	"command"
	"context"
	"fmt"
	"os"
	"rut"
	"time"
)

var CTX = context.Background()

func main() {

	/*
		dev, err := rut.New(&rut.RUT{
			Name:     "SWITCH",
			Device:   "V8",
			IP:       "10.71.20.191",
			Port:     "23",
			Username: "admin",
			Protocol: "telnet",
			Hostname: "V8500_2",
			Password: "",
			SFU:      "A",
		})
	*/
	dev, err := rut.New(&rut.RUT{
		Name:     "SWITCH",
		Device:   "V5",
		IP:       "10.71.20.55",
		Port:     "23",
		Username: "admin",
		Protocol: "telnet",
		Hostname: "SWITCH",
		Password: "",
		SFU:      "A",
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
		CMD:  "do q sh -l",
	})

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))

	data, err = dev.RunCommand(CTX, &command.Command{
		Mode: "shell",
		CMD:  " tcpdump -w 2224g.pcap&",
	})

	time.Sleep(time.Second * 5)

	data, err = dev.RunCommand(CTX, &command.Command{
		Mode: "shell",
		CMD:  " killall tcpdump",
	})

	time.Sleep(time.Second * 5)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if err = dev.FTPPut("/etc/.config/L2_default.CFG", "10.71.1.161", "pi", "raspberry", pwd); err != nil {
		panic(err)
	}

	if err = dev.SCPPut("/etc/login.conf", "10.71.1.161", "pi", "raspberry", pwd); err != nil {
		panic(err)
	}

	dev.TCPDUMP("eth05", "", "1.pcap", "")

	if err = dev.SCPPut("/etc/1.pcap", "10.71.1.161", "pi", "raspberry", pwd); err != nil {
		panic(err)
	}
}
