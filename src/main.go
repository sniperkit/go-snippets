package main

import (
	"log"
	"os"
)

var homeDir string = os.Getenv("HOME")
var gfDir string = homeDir + "/gotfiles"

func main() {
	var LError = log.New(os.Stderr, "", 0)

	args := os.Args[1:]
	argNum := len(args)
	if argNum < 1 {
		LError.Fatalln("You must provide at least one command as an argument!")
	}
	cmd := args[0]

	if cmd == "init" {
		if argNum == 2 {
			initFiles(gfDir + "/" + args[1])
		} else {
			LError.Fatalln("You must provide the name of the repo to initialize!")
		}
	}

	if cmd == "add" {
		if argNum >= 2 {
			if exists(gfDir + "/" + args[1]) {
				for _, p := range args[2:] {
					addPath(gfDir+"/"+args[1], p)
				}
			} else {
				if exists(gfDir + "/.default") {
					for _, p := range args[1:] {
						addPath(gfDir+"/.default", p)
					}
				} else {
					LError.Fatalln("You must specify a default repo, or use use one as an argument")
				}
			}
		} else {
			LError.Fatalln("You must specify a file or directory to add to the repo")
		}
	}

	if cmd == "default" {
		if argNum == 2 {
			if exists(gfDir + "/" + args[1]) {
				if exists(gfDir + "/.default") {
					os.Remove(gfDir + "/.default")
				}
				if err := os.Symlink(gfDir+"/"+args[1], gfDir+"/.default"); err != nil {
					LError.Fatalln("Warning, could not properly create the default directory. Please try again")
				}
			}
		} else {
			LError.Fatalln("You must provide the name of the repo to set as default!")
		}
	}

	if cmd == "sync" {
		// rsync this because I'm too lazy to think of anything better
		if argNum == 2 {
			if exists(gfDir + "/" + args[1]) {
				if err := sync(gfDir + "/" + args[1]); err != nil {
					LError.Fatalln("Warning, could not sync your files. Please ensure that rsync is installed and functioning properly.")
				}
			} else {
				LError.Fatalln("You must specify the name of a repo which actually exists!")
			}
		} else {
			if exists(gfDir + "/.default") {
				if err := sync(gfDir + "/.default"); err != nil {
					LError.Fatalln("Warning, could not sync your files. Please ensure that rsync is installed and functioning properly.")
				}
			} else {
				LError.Fatalln("You must specify the name of a repo or set a default repo!")
			}
		}
	}

	if cmd == "deploy" {
		force := false
		if exists, ind := stringInSliceInd("-f", args); exists {
			args = append(args[:ind], args[ind+1:]...)
			argNum--
			force = true
		}
		// Deploys the entire default repo
		if argNum == 1 {
			if exists(gfDir + "/.default") {
				deploy(gfDir+"/.default/", force)
			} else {
				LError.Fatalln("Please specify a repo to use, or set a default repo!")
			}
		} else if argNum == 2 {
			// If a repo or a specific file/dir has been specified
			if exists(gfDir + "/" + args[1]) {
				deploy(gfDir+"/"+args[1]+"/", force)
			} else {
				if exists(gfDir + "/.default/" + args[1]) {
					deploy(gfDir+"/.default/"+args[1], force)
				} else {
					LError.Fatalln("Please ensure that the default repo is set and that the specified file/dir exists within it!")
				}
			}
		} else if argNum == 3 {
			if exists(gfDir+"/"+args[1]) && exists(gfDir+"/"+args[1]+"/"+args[2]) {
				deploy(gfDir+"/"+args[1]+"/"+args[2], force)
			} else {
				LError.Fatalln("Please ensure that the specified repo exists and that the specified file/dir exists within it!")
			}
		} else {
			LError.Fatalln("You must specify the name of a repo to deploy, or the specific file/directory and name!")
		}
	}
}
