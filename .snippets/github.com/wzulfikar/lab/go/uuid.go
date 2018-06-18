package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/fatih/color"
	"github.com/rs/xid"
	"github.com/satori/go.uuid"
	"github.com/teris-io/shortid"
)

func main() {
	terisshortid()
	rs_xid()
	satoriuuid()
	randSeqGenerator()
}

func terisshortid() {
	defer timer("github.com/teris-io/shortid")()
	fmt.Println("\n=== Generator used: github.com/teris-io/shortid ===\n")

	// initialize generator
	sid, err := shortid.New(1, shortid.DefaultABC, 2342)
	if err != nil {
		log.Fatal(err)
	}

	// // generate the short id
	var id string
	for i := 0; i < 5; i++ {
		id, err = sid.Generate()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("shortid %d: %s\n", i+1, id)
	}
}

func satoriuuid() {
	defer timer("github.com/satori/go.uuid")()
	fmt.Println("\n=== Generator used: github.com/satori/go.uuid ===\n")

	for i := 0; i < 5; i++ {
		fmt.Printf("uuid %d: %s\n", i+1, uuid.NewV4())
	}

	// Parsing UUID from string input
	u2, err := uuid.FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	if err != nil {
		fmt.Println("Something went wrong: %s", err)
	}
	fmt.Println("\nTest parsing uuid from string")
	fmt.Println("Successfully parsed:", u2)
}

func rs_xid() {
	defer timer("github.com/rs/xid")()
	fmt.Println("\n=== Generator used: github.com/rs/xid ===\n")

	for i := 0; i < 5; i++ {
		fmt.Printf("xid %d: %s\n", i+1, xid.New())
	}
}

func randSeqGenerator() {
	defer timer("randSeq")()
	fmt.Println("\n=== Generator used: randSeq ===\n")

	// seed before generating randSeq
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 5; i++ {
		fmt.Printf("xid %d: %s\n", i+1, randSeq(11))
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func timer(taskName string) func() {
	start := time.Now()
	return func() {
		color.Cyan("\n=== " + taskName + " done in " + time.Since(start).String() + " ===\n")
	}
}
