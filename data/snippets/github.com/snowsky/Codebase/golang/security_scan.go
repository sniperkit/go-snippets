package main

// 不管使用AWS还是vSphere，后台都会使用docker来进行扫描
import (
	"fmt"
	"log"
	"net/http"
)

type Scan struct {
	Provider map[string]string // AWS or vSphere
	Command string or path
}

func main() {
	fmt.Println("Starting Security Scan...")
	http.HandleFunc("/v1/scan", newScan)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Scan) newScan(rw http.ResponseWriter, r *http.Request) error {
	fmt.Println("You are in newScan...")
}
