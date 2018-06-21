package main

import (
	"assertion"
	"context"
	"fmt"
	"rut"
)

func main() {
	dev, err := rut.New(&rut.RUT{
		Name:     "DUT1",
		Device:   "V5",
		IP:       "10.71.20.102",
		Port:     "23",
		Username: "admin",
		Hostname: "SWITCH",
		Password: "",
	})

	if err != nil {
		panic(err)
	}

	dev.Init()

	ctx := context.Background()

	ruts := rut.NewDB()
	ruts.AddDUT("0", dev)

	a := assertion.New("0", "normal", "show running-config", "interface")
	res, ok := a.Do(ctx, ruts)
	fmt.Println(res, ok)

	a = assertion.New("0", "normal", "show running-config", "!!interface")
	res, ok = a.Do(ctx, ruts)
	fmt.Println(res, ok)
}
