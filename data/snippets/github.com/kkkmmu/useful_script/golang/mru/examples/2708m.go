package main

import (
	"cline"
	"command"
	"configuration"
	"fmt"
)

func main() {
	c, err := cline.NewCli(&configuration.Configuration{
		Device:         "V2708M",
		IP:             "10.71.20.177",
		Port:           "23",
		Username:       "admin",
		Password:       "",
		EnablePrompt:   ">",
		LoginPrompt:    "login",
		PasswordPrompt: "Password",
		Prompt:         "#",
		Name:           "DUT1",
		ModeDB: map[string]string{
			"login":    "login",
			"password": "Passowrd:",
			"enable":   "SWITCH>",
			//	"normal":        "#",
			"config":        "(config)",
			"config-vlan":   "(config-vlan)",
			"config-if":     "(config-if[",
			"config-dhcp":   "(config-dhcp[",
			"config-router": "(config-router)",
			"shell":         "*SWITCH",
			"bcmshell":      "BCM.0>",
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
		Mode: "normal",
		CMD:  "configure terminal",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(data))

	data, err = c.RunCommand(&command.Command{
		Mode: "config",
		CMD:  "show running-config",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(data))

	data, err = c.RunCommand(&command.Command{
		Mode: "config",
		CMD:  "do q sh -l",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(data))

	data, err = c.RunCommand(&command.Command{
		Mode: "shell",
		CMD:  "scontrol -f /proc/switch/ASIC/ctrl dump table 0 trunk_group 0 2",
		End:  "#",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(data))

	data, err = c.RunCommand(&command.Command{
		Mode: "shell",
		CMD:  "bcm.user.proxy\n\n",
		End:  "BCM.0>",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(data))

	data, err = c.RunCommand(&command.Command{
		Mode: "bcmshell",
		CMD:  "mod trunk_group 0 1 rtag=7",
		End:  "BCM.0>",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(data))

	data, err = c.RunCommand(&command.Command{
		Mode: "bcmshell",
		CMD:  "g RTAG7_HASH_CONTROL_3",
		End:  "BCM.0>",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(data))

	data, err = c.RunCommand(&command.Command{
		Mode: "bcmshell",
		CMD:  "g RTAG7_IPV4_TCP_UDP_HASH_FIELD_BMAP_1",
		End:  "BCM.0>",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(data))

	data, err = c.RunCommand(&command.Command{
		Mode: "bcmshell",
		CMD:  "g RTAG7_IPV4_TCP_UDP_HASH_FIELD_BMAP_2",
		End:  "BCM.0>",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(data))

	data, err = c.RunCommand(&command.Command{
		Mode: "bcmshell",
		CMD:  "g RTAG7_HASH_FIELD_BMAP_1",
		End:  "BCM.0>",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(data))

	data, err = c.RunCommand(&command.Command{
		Mode: "bcmshell",
		CMD:  "exit",
		End:  "#",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(data))

	data, err = c.RunCommand(&command.Command{
		Mode: "shell",
		CMD:  "scontrol -f /proc/switch/ASIC/ctrl dump table 0 trunk_group 0 2",
		End:  "#",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(data))
}
