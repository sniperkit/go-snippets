package main

import (
	"log"
	"os"
)

func initFiles(name string) {
	if !exists(gfDir) {
		if err := os.Mkdir(gfDir, 0755); err != nil {
			log.Fatalln("Error, could not properly initialize the repo. Please check permissions and try again.")
		}
	}
	if !exists(name) {
		if err := os.Mkdir(name, 0755); err != nil {
			log.Fatalln("Error, could not properly initialize the repo. Please check permissions and try again.")
		}
	}
	if !exists(gfDir + "/.default") {
		if err := os.Symlink(name, gfDir+"/.default"); err != nil {
			log.Fatalln("Error, could not properly set the default repo. Please check permissions and try again.")
		}
	}
}
