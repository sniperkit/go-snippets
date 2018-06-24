package main

import (
	"bufio"
	"fmt"
	"github.com/rsmohamad/comp4321/retrieval"
	"os"
)

func main() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter search term: ")
		query, _ := reader.ReadString('\n')

		se := retrieval.NewSearchEngine("index.db")
		defer se.Close()

		results := se.RetrieveVSpace(query)
		for _, doc := range results {
			fmt.Println(doc.Title, doc.Score)
		}
	}
}
