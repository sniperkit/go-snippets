package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"plugin"
	"strings"
	"time"

	"github.com/andrestc/go-plugserver/plughttp"
)

type fileEntry struct {
	pattern string
	modTime time.Time
}

var fileIndex map[string]fileEntry

func main() {
	fileIndex = make(map[string]fileEntry)
	mux := plughttp.NewServeMux()
	log.Println("Starting server on port :8080")
	go watchHandlers(mux, "./handlers")
	http.ListenAndServe(":8080", mux)
}

func watchHandlers(mux *plughttp.ServeMux, pluginsDir string) {
	loadHandlers(mux, pluginsDir)
	for {
		select {
		case <-time.After(5 * time.Second):
			log.Println("Refreshing handlers...")
			loadHandlers(mux, pluginsDir)
			removeOldHandlers(mux, pluginsDir)
			log.Println("Done refreshing handlers...")
		}
	}
}

func loadHandlers(mux *plughttp.ServeMux, pluginsDir string) {
	dir, err := os.Open(pluginsDir)
	if err != nil {
		panic(fmt.Sprintf("unable to open plugins dir: %s", err))
	}
	files, err := dir.Readdir(0)
	if err != nil {
		log.Printf("failed to read plugins dir content: %s", err)
	}
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".so") {
			continue
		}
		if old, ok := fileIndex[f.Name()]; ok {
			if !f.ModTime().After(old.modTime) {
				log.Printf("Skipping %q. Already scanned.", f.Name())
				continue
			}
		}
		p, err := plugin.Open(filepath.Join(dir.Name(), f.Name()))
		if err != nil {
			log.Printf("failed to read plugin file: %s", err)
			continue
		}
		handler, err := p.Lookup("Handler")
		if err != nil {
			log.Printf("failed to lookup handler in plugin %s: %s", f.Name(), err)
			continue
		}
		fun, ok := handler.(func() (string, func(w http.ResponseWriter, r *http.Request)))
		if !ok {
			log.Printf("invalid Handler function on plugin %q, has type %T", f.Name(), handler)
			continue
		}
		pattern, handleFunc := fun()
		log.Printf("Registering handler for pattern %q from file %q", pattern, f.Name())
		mux.HandleFunc(pattern, handleFunc)
		fileIndex[f.Name()] = fileEntry{pattern: pattern, modTime: f.ModTime()}
	}
}

func removeOldHandlers(mux *plughttp.ServeMux, pluginsDir string) {
	checkedMap := make(map[string]bool)
	dir, err := os.Open(pluginsDir)
	if err != nil {
		panic(fmt.Sprintf("unable to open plugins dir: %s", err))
	}
	files, err := dir.Readdir(0)
	if err != nil {
		log.Printf("failed to read plugins dir content: %s", err)
	}
	for _, f := range files {
		checkedMap[f.Name()] = true
	}
	for k, v := range fileIndex {
		if !checkedMap[k] {
			log.Printf("Removing pattern %q previously added from file %q", v.pattern, k)
			mux.DeRegister(v.pattern)
			delete(fileIndex, k)
		}
	}
}
