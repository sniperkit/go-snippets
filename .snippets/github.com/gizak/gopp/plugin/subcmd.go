package plugin

import (
	"flag"
	"fmt"
	"os"
)

var PLUGIN cmdPlugin

type cmdPlugin struct {
	*flag.FlagSet
}

func (cmdPlugin) Subcmd() string {
	return "plugin"
}

func (cmdPlugin) Usage() string {
	fmt.Println("Plugin!!")
	return ""
}

func (cmdPlugin) Descrip() string {
	return "manage gopp plugins"
}

func (cp cmdPlugin) RunSubcmd(args []string) error {
	cp.FlagSet = flag.NewFlagSet("plugin", flag.ContinueOnError)
	//cp.FlagSet.Usage = func() { fmt.Print(cp.Usage()) }

	islist := cp.Bool("list", false, "List all loaded plugins")
	isupdate := cp.Bool("update", false, "Update and rebuild plugins")
	cp.Parse(args[1:])

	if *islist {
		fmt.Print("The installed plugins:\n\n")
		for _, p := range Plugins {
			fmt.Printf("\t%-12s%s\n", p.Name, p.Uri)
		}
	}

	if *isupdate {
		SetLogOutput(os.Stdout)
		if err := InstallWithConfig(); err != nil {
			return err
		}
	}

	return nil
}
