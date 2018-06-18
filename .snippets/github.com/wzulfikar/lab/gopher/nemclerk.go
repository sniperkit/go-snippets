package main

import (
	"fmt"
	"log"

	nemclient "github.com/wzulfikar/go-nem-client"

	prettyjson "github.com/hokaccha/go-prettyjson"
)

func main() {
	c, _ := nemclient.NewClient("http://23.228.67.85:7890")

	tx, err := c.GetAllTransactions("TC6Z5CY3TX4V3UCTWFZLCSEH6WOSIZ4PU2RURAHN", "", "")
	// tx, err := c.GetUnconfirmedTransactions("TDJDSDTS2UD2TLEZEO5DEG6BABJ64M6FTZTHQI6E")
	if err != nil {
		log.Fatal(err)
	}

	s, _ := prettyjson.Marshal(tx)
	fmt.Println(string(s))
}
