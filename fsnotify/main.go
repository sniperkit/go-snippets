package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	checkError(err)
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("Event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("File:", event.Name)
					content := readFile("./info.txt")
					fmt.Println("====================")
					fmt.Println(content)
					fmt.Println("====================")
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("./info.txt")
	checkError(err)
	<-done
}

func readFile(fileName string) string {
	content, err := ioutil.ReadFile(fileName)
	checkError(err)
	return string(content)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
