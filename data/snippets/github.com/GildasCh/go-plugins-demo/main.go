package main

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"path"
	"plugin"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/net/websocket"
)

type Transformer struct {
	name     string
	t        func(image.Image) image.Image
	priority int // priority can be 0 or 1
}

func loadPlugin(path string) (*Transformer, error) {
	fmt.Println("Loading", path)
	p, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}

	transformSymbol, err := p.Lookup("Transform")
	if err != nil {
		return nil, err
	}
	transform, ok := transformSymbol.(func(image.Image) image.Image)
	if !ok {
		return nil, errors.New(fmt.Sprintf("%s: Transform not a func(image.Image) image.Image", path))
	}

	var priority int
	prioritySymbol, err := p.Lookup("Priority")
	if err != nil {
		fmt.Println("Error looking up Priority:", err)
	} else {
		var ok bool
		priority, ok = prioritySymbol.(int)
		if !ok {
			fmt.Println("Error casting Priority to int")
		}
	}
	return &Transformer{path, transform, priority}, nil
}

func getPlugins(pluginDir string) (map[int][]*Transformer, error) {
	pDir, err := os.Open(pluginDir)
	if err != nil {
		return nil, err
	}
	pFiles, err := pDir.Readdir(0)
	if err != nil {
		return nil, err
	}

	ret := make(map[int][]*Transformer)
	for _, pFile := range pFiles {
		if pFile.IsDir() {
			continue
		}

		p, err := loadPlugin(path.Join(pluginDir, pFile.Name()))
		if err != nil {
			log.Printf("Failed to load pluging %s", p)
		} else {
			ret[p.priority] = append(ret[p.priority], p)
		}
	}
	return ret, nil
}

var processed image.Image
var listeners []chan bool

func processImage(events <-chan fsnotify.Event) {
	for {
		fi, _ := os.Open("lca.jpg")
		i, _ := jpeg.Decode(fi)

		plugins, err := getPlugins(pluginDir)
		if err != nil {
			panic(err)
		}

		for _, p := range plugins[1] {
			fmt.Println("Apply", p.name)
			i = p.t(i)
		}
		for _, p := range plugins[0] {
			fmt.Println("Apply", p.name)
			i = p.t(i)
		}
		processed = i
		for _, l := range listeners {
			l <- true
		}

		// Wait for a change in folder
		<-events
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	jpeg.Encode(w, processed, nil)
}

func wsHandler(ws *websocket.Conn) {
	l := make(chan bool)
	listeners = append(listeners, l) // Should have mutex
	for {
		fmt.Fprintf(ws, "echp")
		<-l
	}
}

var pluginDir string

func main() {
	pluginDir = os.Args[1]

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go processImage(watcher.Events)
	err = watcher.Add(pluginDir)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/a.jpg", handler)
	http.Handle("/ws", websocket.Handler(wsHandler))
	http.ListenAndServe(":8080", nil)
}
