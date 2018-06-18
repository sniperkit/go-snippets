// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package index

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"

	elastic "gopkg.in/olivere/elastic.v5"
)

func ExampleBodyString() {

	// This works for TestExists
	// client := setupTestClientAndCreateIndexAndAddDocs(t, SetTraceLog(log.New(os.Stdout, "", 0)))

	// Obtain a client. You can also provide your own HTTP client here.

	// This is the default way
	// errorlog := log.New(os.Stdout, "APP ", log.LstdFlags)
	// client, err := elastic.NewClient(elastic.SetErrorLog(errorlog))

	// Do a trace log
	tracelog := log.New(os.Stdout, "", 0)
	client, err := elastic.NewClient(elastic.SetTraceLog(tracelog))

	if err != nil {
		// Handle error
		panic(err)
	}

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("twitter").Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("twitter").Do(context.Background())
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

	tweet1 := `{"user" : "olivere", "town" : "Pittsburgh", "message" : "It's a Good Life"}`
	put1, err := client.Index().
		Index("twitter").
		Type("tweet").
		Id("10").
		BodyString(tweet1).
		Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed tweet %s \n", put1.Id)

	tweet2 := `{"user" : "smith", "message" : "It's a Raggy Waltz"}`
	put2, err := client.Index().
		Index("twitter").
		Type("tweet").
		Id("20").
		BodyString(tweet2).
		Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed tweet %s \n", put2.Id)

	tweet3 := `{"user" : "cohn", "party" : "temporary"}`
	put3, err := client.Index().
		Index("twitter").
		Type("tweet").
		Id("30").
		BodyString(tweet3).
		Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed tweet %s \n", put3.Id)

	// Get tweet with specified ID
	get1, err := client.Get().
		Index("twitter").
		Type("tweet").
		Id("10").
		Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	if get1.Found {
		fmt.Printf("Got document %s \n", get1.Id)
	}

	// Get tweet with specified ID
	get2, err := client.Get().
		Index("twitter").
		Type("tweet").
		Id("20").
		Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	if get2.Found {
		fmt.Printf("Got document %s \n", get2.Id)
	}

	// Get tweet with specified ID
	get3, err := client.Get().
		Index("twitter").
		Type("tweet").
		Id("30").
		Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	if get3.Found {
		fmt.Printf("Got document %s \n", get3.Id)
	}

	// Flush to make sure the documents got written.
	_, err = client.Flush().Index("twitter").Do(context.Background())
	if err != nil {
		panic(err)
	}
}
