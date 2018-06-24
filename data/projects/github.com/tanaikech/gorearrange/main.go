// Package main (gorearrange.go) :
// This file is included all commands and options.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	rearrange "github.com/tanaikech/go-rearrange"
	"github.com/urfave/cli"
)

// const :
const (
	appname = "gorearrange"
)

// main : Get command-line arguments.
func main() {
	app := cli.NewApp()
	app.Name = appname
	app.Author = "tanaike [ https://github.com/tanaikech/gorearrange ] "
	app.Email = "tanaike@hotmail.com"
	app.Usage = "Manually rearrange text data."
	app.Version = "1.0.2"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "step, s",
			Value: 5,
			Usage: "Number of steps for PageUp, PageDown.",
		},
		cli.StringFlag{
			Name:  "inputfile, i",
			Value: "",
			Usage: "Input text file that you want to rearrange.",
		},
		cli.StringFlag{
			Name:  "outputfile, o",
			Value: "",
			Usage: "Output rearranged text file.",
		},
		cli.BoolFlag{
			Name:  "selectmode, select",
			Usage: "Use as select mode. In this case, users only select a value.",
		},
		cli.BoolFlag{
			Name:  "indexmode, index",
			Usage: "After rearranging source data, the change of index for the source data is output.",
		},
		cli.BoolFlag{
			Name:  "history, hi",
			Usage: "Show history of selected data.",
		},
	}
	app.Action = gorearrange
	app.Run(os.Args)
}

// getValues : Get values from stdin and file.
func getValues(filepath string) []string {
	var data []string
	if terminal.IsTerminal(int(syscall.Stdin)) {
		if len(filepath) == 0 {
			fmt.Printf("Please input a text data file.\n\n $ cat filename | %s\n\n or \n\n $ %s -f filename\n\n Help is \"%s --help\"\n", appname, appname, appname)
			os.Exit(0)
		} else {
			f, err := os.Open(filepath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Script '%s' is not found.\n", filepath)
				os.Exit(1)
			}
			defer f.Close()
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				if scanner.Text() == "end" {
					break
				}
				data = append(data, scanner.Text())
			}
			if scanner.Err() != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", scanner.Err())
				os.Exit(1)
			}
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if scanner.Text() == "end" {
				break
			}
			data = append(data, scanner.Text())
		}
		if scanner.Err() != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", scanner.Err())
			os.Exit(1)
		}
	}
	if len(data) == 0 {
		fmt.Fprintf(os.Stderr, "Error: No data. Please check help.\n\n $ %s help\n\n", appname)
		os.Exit(1)
	}
	return data
}

// gorearrange : Handler for gorearrange
func gorearrange(c *cli.Context) error {
	result, history, err := rearrange.Do(getValues(c.String("inputfile")), c.Int("step"), c.Bool("selectmode"), c.Bool("indexmode"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outdata := strings.Join(result, "\n") + "\n"
	if !c.Bool("selectmode") {
		if len(c.String("outputfile")) == 0 {
			w := bufio.NewWriter(os.Stdout)
			w.WriteString(outdata)
			w.Flush()
		} else {
			ioutil.WriteFile(c.String("outputfile"), []byte(outdata), os.ModePerm)
		}
		if c.Bool("history") {
			r, _ := json.Marshal(history)
			fmt.Printf("\n%+v\n", string(r))
		}
	} else {
		w := bufio.NewWriter(os.Stdout)
		w.WriteString(history[len(history)-1].Value + "\n")
		w.Flush()
	}
	return nil
}
