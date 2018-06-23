package rut_test

import (
	"assert"
	"command"
	"configuration"
	"fmt"
	"routine"
	"rut"
	"testing"
)

func TestNewDut(t *testing.T) {
	base := "SWITCH[A]"
	_, err := rut.NewRUT(&config.Config{
		DeviceName:     "V8500",
		IP:             "10.71.20.198",
		Port:           "23",
		UserName:       "admin",
		Password:       "",
		EnablePrompt:   ">",
		LoginPrompt:    "login",
		PasswordPrompt: "Password",
		BasePrompt:     base,
		Prompt:         "#",
		ModeDB: map[string]string{
			"login":         "login",
			"password":      "Passowrd:",
			"enable":        base + ">",
			"normal":        base + "#",
			"config":        base + "(config)",
			"config-vlan":   base + "(config-vlan)",
			"config-if":     base + "(config-if[",
			"config-dhcp":   base + "(config-dhcp[",
			"config-router": base + "(config-router)",
		},
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func TestDutRunCommand(t *testing.T) {
	base := "SWITCH[A]"
	d, err := dut.NewRUT(&config.Config{
		DeviceName:     "V8500",
		IP:             "10.71.20.198",
		Port:           "23",
		UserName:       "admin",
		Password:       "",
		EnablePrompt:   ">",
		LoginPrompt:    "login",
		PasswordPrompt: "Password",
		BasePrompt:     base,
		Prompt:         "#",
		ModeDB: map[string]string{
			"login":         "login",
			"password":      "Passowrd:",
			"enable":        base + ">",
			"normal":        base + "#",
			"config":        base + "(config)",
			"config-vlan":   base + "(config-vlan)",
			"config-if":     base + "(config-if[",
			"config-dhcp":   base + "(config-dhcp[",
			"config-router": base + "(config-router)",
		},
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	data, err := d.RunCommand(&command.Command{
		RequiredMode: "normal",
		CMD:          "configure terminal",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(data))

	data, err = d.RunCommand(&command.Command{
		RequiredMode: "config",
		CMD:          "show running-config",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

func TestDutRunRoutine(t *testing.T) {
	base := "SWITCH[A]"
	d, err := dut.NewRUT(&config.Config{
		DeviceName:     "V8500",
		IP:             "10.71.20.198",
		Port:           "23",
		UserName:       "admin",
		Password:       "",
		EnablePrompt:   ">",
		LoginPrompt:    "login",
		PasswordPrompt: "Password",
		BasePrompt:     base,
		Prompt:         "#",
		ModeDB: map[string]string{
			"login":         "login",
			"password":      "Passowrd:",
			"enable":        base + ">",
			"normal":        base + "#",
			"config":        base + "(config)",
			"config-vlan":   base + "(config-vlan)",
			"config-if":     base + "(config-if[",
			"config-dhcp":   base + "(config-dhcp[",
			"config-router": base + "(config-router)",
		},
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	routines := []*routine.Routine{
		{
			Name: "Create VLAN Interface",
			CMD: []*command.Command{
				{
					RequiredMode: "normal",
					CMD:          "configure terminal",
				},
				{
					RequiredMode: "config",
					CMD:          "vlan database",
				},
				{
					RequiredMode: "config-vlan",
					CMD:          "vlan 1234",
				},
				{
					RequiredMode: "config-vlan",
					CMD:          "exit",
				},
				{
					RequiredMode: "config",
					CMD:          "interface vlan 1234",
				},
				{
					RequiredMode: "config-if",
					CMD:          "ip address 123.1.1.1/24",
				},
				{
					RequiredMode: "config-if",
					CMD:          "no shutdown",
				},
				{
					RequiredMode: "config-if",
					CMD:          "exit",
				},
				{
					RequiredMode: "config",
					CMD:          "interface giga 10/6",
				},
				{
					RequiredMode: "config-if",
					CMD:          "switchport access vlan 1234",
				},
				{
					RequiredMode: "config-if",
					CMD:          "exit",
				},
				{ //Go back to normal mode
					RequiredMode: "config",
					CMD:          "exit",
				},
			},
			Assert: []*assert.Assert{
				{
					Expected: "1234",
					CMD:      &command.Command{RequiredMode: "normal", CMD: "show vlan"},
					Done:     false,
				},
				{
					Expected: `br1234[[:space:]]+up[[:space:]]+up[[:space:]]+123\.1\.1\.1`,
					CMD:      &command.Command{RequiredMode: "normal", CMD: "show ip interface brief"},
					Done:     false,
				},
				{
					Expected: `UP,BROADCAST,RUNNING,MULTICAST`,
					CMD:      &command.Command{RequiredMode: "normal", CMD: "show interface vlan 100"},
					Done:     false,
				},
			},
		},
		{
			Name: "Destroy VLAN Interface",
			CMD: []*command.Command{
				{
					RequiredMode: "normal",
					CMD:          "configure terminal",
				},
				{
					RequiredMode: "config",
					CMD:          "interface vlan 1234",
				},
				{
					RequiredMode: "config-if",
					CMD:          "shutdown",
				},
				{
					RequiredMode: "config-if",
					CMD:          "exit",
				},
				{
					RequiredMode: "config",
					CMD:          "no interface vlan 1234",
				},
				{
					RequiredMode: "config",
					CMD:          "vlan database",
				},
				{
					RequiredMode: "config-vlan",
					CMD:          "no vlan 1234",
				},
				{
					RequiredMode: "config-vlan",
					CMD:          "exit",
				},
				{
					RequiredMode: "config",
					CMD:          "exit",
				},
			},
			Assert: []*assert.Assert{
				{
					UnExpected: "1234",
					CMD:        &command.Command{RequiredMode: "normal", CMD: "show vlan"},
					Done:       false,
				},
				{
					UnExpected: `br1234[[:space:]]+up[[:space:]]+up[[:space:]]+123\.1\.1\.1`,
					CMD:        &command.Command{RequiredMode: "normal", CMD: "show ip interface brief"},
					Done:       false,
				},
				{
					UnExpected: `UP,BROADCAST,RUNNING,MULTICAST`,
					CMD:        &command.Command{RequiredMode: "normal", CMD: "show interface vlan 100"},
					Done:       false,
				},
			},
		},
	}

	for _, r := range routines {
		err := d.RunRoutine(r)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func TestDutRunTask(t *testing.T) {
	_, err := rut.NewRUT(&configuraiton.Configuration{
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
}
