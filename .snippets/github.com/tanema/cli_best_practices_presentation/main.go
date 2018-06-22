package main

import "github.com/spf13/cobra"

var (
	mainCmd = &cobra.Command{
		Use:   "meetup",
		Short: "A tool for managing your meetup ",
		Long: `Meetup is a tool for managing your meetup from the command line.
Meetup allows you to create, delete and edit events on your meetup`,
	}
	newCmd     = &cobra.Command{Use: "new", Short: "creates a new event"}
	configPath string
)

func init() {
	mainCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "/var/app_config.yml", "path to config.yml")
	mainCmd.AddCommand(newCmd)
}

func main() { mainCmd.Execute() }
