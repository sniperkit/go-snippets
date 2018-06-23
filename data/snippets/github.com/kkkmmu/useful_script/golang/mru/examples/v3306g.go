package main

import (
	"command"
	"context"
	"fmt"
	"rut"
	"time"
	//"util"
	"strings"
)

func main() {
	dev, err := rut.New(&rut.RUT{
		Name:     "SWITCH",
		Device:   "V5",
		IP:       "10.55.192.212",
		Port:     "23",
		Username: "admin",
		Hostname: "SWITCH",
		Password: "Dasan123456",
	})

	if err != nil {
		panic(err)
	}

	err = dev.Init()
	if err != nil {
		fmt.Println(err.Error())
	}

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

	/*
		data, err = dev.RunCommand(ctx, &command.Command{
			Mode: "shell",
			CMD:  " bcm cos bandwidth_show",
		})
		fmt.Println(data)
		util.SaveToFile("cos_bw.txt", []byte(data))

		data, err = dev.RunCommand(ctx, &command.Command{
			Mode: "shell",
			CMD:  " bcm fp show",
		})
		fmt.Println(data)
		util.SaveToFile("fp_show.txt", []byte(data))

		data, err = dev.RunCommand(ctx, &command.Command{
			Mode: "shell",
			CMD:  " bcm d chg fp_tcam",
		})
		fmt.Println(data)
		util.SaveToFile("fp_tcam.txt", []byte(data))

		data, err = dev.RunCommand(ctx, &command.Command{
			Mode: "shell",
			CMD:  " bcm d chg fp_policy_table",
		})
		fmt.Println(data)
		util.SaveToFile("fp_policy_table.txt", []byte(data))
	*/
	for k := 0; k < 9; k++ {
		for i := 0; i < 32; i++ {
			var result string
			_, err = dev.RunCommand(ctx, &command.Command{
				Mode: "shell",
				CMD:  fmt.Sprintf("bcm s RDBGC%d_SELECT %d", k, 1<<uint(i)),
			})

			_, err = dev.RunCommand(ctx, &command.Command{
				Mode: "shell",
				CMD:  fmt.Sprintf("bcm s RDBGC%d 0", k),
			})

			_, err = dev.RunCommand(ctx, &command.Command{
				Mode: "shell",
				CMD:  fmt.Sprintf("bcm s RDBGC%d 0", k),
			})
			time.Sleep(time.Millisecond * 50)
			data, err = dev.RunCommand(ctx, &command.Command{
				Mode: "shell",
				CMD:  fmt.Sprintf("bcm g chg RDBGC%d", k),
			})
			result += string(data)
			if strings.Contains(result, "0x") {

				fmt.Printf("Checking RDGBC%d bit %d\n", k, i)
				fmt.Println(result)
			}
		}
	}
}
