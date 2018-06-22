package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/ryanuber/columnize"
	"github.com/tylerfowle/dtags/db"
	"github.com/tylerfowle/dtags/install"
)

var database *db.Database

// Entrypoint for dtags application. Initializes database instance and
// processed passed command to perform the expected action.
func main() {
	var err error
	if database, err = db.Init(); err != nil {
		panic(err)
	}
	defer database.Instance.Close()

	if len(os.Args) < 2 || ((os.Args[1] == "add" || os.Args[1] == "del") && len(os.Args) < 3) {
		showHelp()
	}

	switch os.Args[1] {
	case "add":
		addNewTag(os.Args[2:])
		break
	case "del":
		if err := database.DeleteKey(strings.ToLower(os.Args[2])); err != nil {
			fmt.Println("Failed.")
		} else {
			fmt.Printf("Successfully deleted tag [%s]\n", os.Args[2])
		}
		break
	case "list":
		for _, v := range database.GetCurrentTags() {
			fmt.Println(v)
		}
		break
	case "tags", "completion":
		for _, tag := range database.GetTags() {
			fmt.Println(tag)
		}
		break
	case "ls":
		printBoth()
		break
	case "install":
		_, err := install.WriteFile()
		if err != nil {
			log.Fatal("Installation failed.", err)
		}
		break
	default:
		printPath(os.Args[1:])
	}
}

// Show the help message to a user when not enough arguments are passed
// to the command.
func showHelp() {
	fmt.Println("usage: dtags command [params]")
	fmt.Println("       dtags <tag>")
	fmt.Println()
	fmt.Println("       add <tag> [path] - add a new tag at path; defaults to current directory")
	fmt.Println("       del <tag> - delete the provided tag from storage")
	fmt.Println("       list - list all tags on the current directory")
	fmt.Println("       ls - list all tags currently stored and associated paths")

	os.Exit(64)
}

// Add a new tag to the database prompting user for permission if the tag
// already exists. Uses the current working directory if no path is provided
// as an argument when calling the command.
func addNewTag(args []string) {
	k := strings.ToLower(args[0])
	v := database.CurrentDirectory

	if len(args[0:]) > 1 {
		v = args[1]
	}

	if database.Exists(k) {
		fmt.Printf("Overwrite existing tag? [%s] (y/n)", k)
		if confirmation() == false {
			return
		}
	}

	if err := database.AddKey(k, v); err != nil {
		fmt.Printf("Failed.\n")
	} else {
		fmt.Printf("Successfully added tag [%s] to path [%s]\n", k, v)
	}

}

// Print all of the tags currently stored in the database along with the
// folder path that belongs to that tag.
func printBoth() {
	var unformattedlist []string
	for tag, path := range database.All() {
		unformattedlist = append(unformattedlist, fmt.Sprintf("%s|%s\n", tag, path))
	}
	sort.Strings(unformattedlist)
	formattedList := columnize.SimpleFormat(unformattedlist)
	fmt.Println(formattedList)
}

// Print the path for a specific tag that is passed in. If no tag is found, an
// empty response is returned.
func printPath(args []string) {
	cwd := database.GetValue(args[0])
	if cwd == "" {
		fmt.Println("")
		os.Exit(1)
	}

	fmt.Fprint(os.Stdout, cwd)
	os.Exit(1)
}

// Confirm with the user that they do wish to continue with a destructive operation
// or would prefer to cancel and not continue.
func confirmation() bool {
	var response string

	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}

	response = strings.ToLower(strings.TrimSpace(response))

	switch response {
	case "y", "yes":
		fmt.Printf("Overwriting tag [%s]\n", os.Args[1])
		return true
	case "n", "no":
		fmt.Println("Overwrite cancelled.")
		return false
	default:
		fmt.Println("(yes/no) required:")
		return confirmation()
	}
}
