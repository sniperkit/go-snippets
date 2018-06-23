package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
	"strings"
	"strconv"

	"github.com/fatih/color"
	"github.com/urfave/cli"
	"github.com/atotto/clipboard"

	"./diceware_wordlist"
)

func throwdice() string {
	var die [5]string
	for i := 0; i < 5; i++ {
		die[i] = strconv.Itoa(rand.Intn(6) + 1)
	}
	return diceware_wordlist.DiceWare[strings.Join(die[:], "")]
}

func pwgen(count int, sep string) string {
	pswds := make([]string, count)
	for i := 0; i < count; i++ {
		pswds[i] = throwdice()
	}
	return strings.Join(pswds[:], sep)
}


func main() {
	rand.Seed(time.Now().UnixNano())

	var count int
	var sep string
	var clip bool

	app := cli.NewApp()
	app.Name = "xkpswd"
	app.Version = "0.1.2"
	app.Usage = fmt.Sprintf("Generate diceware-xkcd-ish passwords")
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Prashant Sinha",
			Email: "prashant@noop.pw",
		},
	}

	app.Flags = []cli.Flag {
		cli.IntFlag{
			Name: "count, c",
			Value: 5,
			Usage: "specify length of phrase",
			Destination: &count,
		},
		cli.StringFlag{
			Name: "sep, s",
			Value: "-",
			Usage: "specify phrase separator",
			Destination: &sep,
		},
		cli.BoolFlag{
			Name: "noclip, x",
			Usage: "do not copy to clipboard",
			Destination: &clip,
		},
	}

	app.Action = func(c *cli.Context) error {
		password := pwgen(count, sep)
		if !clip {
			clipboard.WriteAll(password)
		}
		color.Set(color.FgCyan)
		fmt.Printf("xkpswd <%d>: ", count)
		color.Set(color.FgGreen, color.Bold)
		fmt.Println(password)
		color.Unset()
		return nil
	}

	defer func() {
		r := recover()
		if r != nil {
			if count < 0 {
				fmt.Println(
					color.CyanString("xkpswd"),
					color.RedString("`count` should be greater than zero."))
			} else {
				panic(r)
			}
		}
	}()

	app.Run(os.Args)
}
