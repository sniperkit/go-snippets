package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// LogError calls log.Printf on error, and adds location in source code
func LogError(err error) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Printf("%s:%d : %s", file, line, err)
	} else {
		log.Println(err)
	}
}

// LogHttp log error and sends it to browser
func LogHttp(w http.ResponseWriter, err error) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Printf("%s:%d : %s", file, line, err)
	} else {
		log.Println(err)
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// LogFatal calls log.Fatalf on error, and adds location in source code
func LogFatal(err error) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Fatalf("%s:%d : %s", file, line, err)
	} else {
		log.Fatal(err)
	}
	
}

func ko(w http.ResponseWriter) {
	w.Write([]byte("ko"))
}

func writeFiles(w http.ResponseWriter, fs ...string) error {
	for _, f := range fs {
		b, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}
		w.Write(b)
	}
	return nil
}

func SetInfo(w http.ResponseWriter, msg string) {
	cookie := &http.Cookie {
		Name	:	"tags-info",
		Value	:	strings.Replace(msg, " ", "_", -1),
		Path	:	"/",
	}
	http.SetCookie(w, cookie)
}

func SetError(w http.ResponseWriter, err error) {
	SetInfo(w, "Error: "+err.Error())
}

func q(s string) string {
	return "'"+s+"'"
}
