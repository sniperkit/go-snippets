package main

import (
	"github.com/gophergala/ket/server"
	"log"
	"path/filepath"
)

func main() {
	log.SetFlags(log.Flags() | log.Llongfile)
	config, err := server.LiveConfig("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	//root := "./data"
	//certFile := filepath.Join(root, "cert.pem")
	//keyFile := filepath.Join(root, "key.pem")

	//ca := server.NewCertAuthority(certFile, keyFile)
	//err = ca.Init()
	//if err != nil {
	//	log.Fatal(err)
	//}
	srv := &server.Server{
		Config: config,
		//CA:     ca,
	}
	err = srv.Init()
	if err != nil {
		log.Fatal(err)
	}
	errors := make(chan error)
	go func() {
		errors <- srv.Start(":4891")
	}()
	//go func() {
	//	errors <- srv.StartTLS(":4819", certFile, keyFile)
	//}()

	log.Fatal(<-errors)
}
