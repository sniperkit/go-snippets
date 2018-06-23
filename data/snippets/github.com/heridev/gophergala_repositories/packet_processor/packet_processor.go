package main

import "log"
import "os"
import "net/http"
import "fmt"

func hello(lPtr *log.Logger) {
	lPtr.Print("Hello from a go routine")
}

func main() {
	log := log.New(os.Stdout, "logger:", log.Lshortfile)
	for i := 0; i < 100; i++ {
		go hello(log)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello\n")
	})
	log.Fatal(http.ListenAndServe(":8080", nil))

}
