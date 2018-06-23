package cline_test

import (
	"cline"
	"command"
	"configuration"
	"fmt"
	"testing"
)

func TestLogin(t *testing.T) {
	c, err := cline.NewCli(&configuration.Configuration{
		DeviceName:     "V8500",
		IP:             "10.71.20.198",
		Port:           "23",
		UserName:       "admin",
		Password:       "",
		EnablePrompt:   ">",
		LoginPrompt:    "login",
		PasswordPrompt: "Password",
		Prompt:         "#",
		ModeDB: map[string]string{
			"login":    "login",
			"password": "Passowrd:",
			"enable":   ">",
			//	"normal":        "#",
			"config":        "(config)",
			"config-vlan":   "(config-vlan)",
			"config-if":     "(config-if[",
			"config-dhcp":   "(config-dhcp[",
			"config-router": "(config-router)",
		},
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = c.Init()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func TestRunCommand(t *testing.T) {
	c, err := cli.NewCli(&configuration.Configuration{
		DeviceName:     "V8500",
		IP:             "10.71.20.198",
		Port:           "23",
		UserName:       "admin",
		Password:       "",
		EnablePrompt:   ">",
		LoginPrompt:    "login",
		PasswordPrompt: "Password",
		Prompt:         "#",
		ModeDB: map[string]string{
			"login":    "login",
			"password": "Passowrd:",
			"enable":   ">",
			//	"normal":        "#",
			"config":        "(config)",
			"config-vlan":   "(config-vlan)",
			"config-if":     "(config-if[",
			"config-dhcp":   "(config-dhcp[",
			"config-router": "(config-router)",
		},
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = c.Init()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	data, err := c.RunCommand(&command.Command{
		RequiredMode: "normal",
		CMD:          "configure terminal",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(data))

	data, err = c.RunCommand(&command.Command{
		RequiredMode: "config",
		CMD:          "show running-config",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(data))
}
