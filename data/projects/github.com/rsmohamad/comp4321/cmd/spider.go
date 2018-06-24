package main

import (
	"flag"
	"fmt"
	"github.com/rsmohamad/comp4321/database"
	"github.com/rsmohamad/comp4321/webcrawler"
	"time"
)

func main() {
	index, _ := database.LoadIndexer("index.db")
	defer index.Close()

	start := flag.String("start", "http://www.cse.ust.hk/", "-start=<starting url>")
	numPages := flag.Int("pages", 300, "-pages=<number of pages>")
	aggressive := flag.Bool("a", false, "-a")
	flag.Parse()

	startCrawl := time.Now()
	obtained := webcrawler.Crawl(*start, *numPages, index, true, *aggressive)
	elapsed := time.Since(startCrawl)
	fmt.Printf("Indexing %d pages took %s\n", len(obtained), elapsed)
	fmt.Println("Updating term weights...")
	index.UpdateTermWeights()
	fmt.Println("Updating adj list...")
	index.UpdateAdjList()
	fmt.Println("Updating page rank...")
	index.UpdatePageRank()
}
