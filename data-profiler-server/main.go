package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// Conf object holds the server configuration
var Conf *Configuration

// TEngine is the TaskEngine used to submit new tasks
var TEngine *TaskEngine

// Configuration dictates the schema of the yml conf file
type Configuration struct {
	Server struct {
		Listen string
		Dirs   struct {
			Templates string
			Static    string
		}
	}
	Database string
	Logfile  string
	Scripts  struct {
		MDS      string
		Analysis []string
		ML       map[string]string
	}
}

// LoadConfig loads the configuration file in memory
func LoadConfig(filename string) *Configuration {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return nil
	}
	conf := new(Configuration)
	err = yaml.Unmarshal(b, conf)
	if err != nil {
		log.Println(err)
		return nil
	}
	return conf
}

// Server entry point
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Please provide conf file\n")
		os.Exit(1)
	}

	Conf = LoadConfig(os.Args[1])
	if Conf == nil {
		os.Exit(1)
	}
	rand.Seed(int64(time.Now().Nanosecond()))
	setLogger(Conf.Logfile)
	TEngine = NewTaskEngine()

	fs := http.FileServer(http.Dir(Conf.Server.Dirs.Static))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", uiHandler)
	http.HandleFunc("/api/", restHandler)
	err := http.ListenAndServe(Conf.Server.Listen, nil)
	if err != nil {
		log.Println(err)
	}
}

func setLogger(logfile string) {
	if logfile != "" {
		f, er := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if er != nil {
			fmt.Println(er)
			os.Exit(1)
		} else {
			log.SetOutput(f)
		}
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

}
