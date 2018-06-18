// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

// BulkInsert illustrates how to bulk insert documents into Elasticsearch.
//
// It uses two goroutines to do so. The first creates a simple document
// and sends it to the second via a channel. The second goroutine collects
// those documents, creates a bulk request that is added to a Bulk service
// and committed to Elasticsearch after reaching a number of documents.
// The number of documents after which a commit happens can be specified
// via the "bulk-size" flag.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/5.0/docs-bulk.html
// for details on the Bulk API in Elasticsearch.
//
// Example
//
// Bulk index 100.000 documents into the index "warehouse", type "product",
// committing every set of 1.000 documents.
//
//     bulk_insert -index=warehouse -type=product -n=100000 -bulk-size=1000
//     bulkstring -index=warehouse -type=product -n=100 -bulk-size=10
//
package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"golang.org/x/sync/errgroup"
	"github.com/olivere/elastic"
)

func main() {
	var (
		index    = flag.String("index", "", "Elasticsearch index name")
		typ      = flag.String("type", "", "Elasticsearch type name")
		n        = flag.Int("n", 0, "Number of documents to bulk insert")
		bulkSize = flag.Int("bulk-size", 0, "Number of documents to collect before committing")
	)
	flag.Parse()
	log.SetFlags(0)
	rand.Seed(time.Now().UnixNano())

	if *index == "" {
		log.Fatal("missing index parameter")
	}
	if *typ == "" {
		log.Fatal("missing type parameter")
	}
	if *n <= 0 {
		log.Fatal("n must be a positive number")
	}
	if *bulkSize <= 0 {
		log.Fatal("bulk-size must be a positive number")
	}

	// Do a trace log
	tracelog := log.New(os.Stdout, "", 0)
	client, err := elastic.NewClient(elastic.SetTraceLog(tracelog))
	// Or with nothing...
	// client, err := elastic.NewClient()

	if err != nil {
		// Handle error
		log.Fatal(err)
	}

	// Setup a group of goroutines from the excellent errgroup package
	g, ctx := errgroup.WithContext(context.TODO())

	// The first goroutine will emit documents and send it to the second goroutine
	// via the docsc channel.
	// The second Goroutine will simply bulk insert the documents.
	docsc := make(chan string)

	// Goroutine to create documents
	g.Go(func() error {
		defer close(docsc)

		for i := 0; i < *n; i++ {

			// Construct the json string
			d := `{"user":"olivere","city":"santafe","age":56}`

			// Send over to 2nd goroutine, or cancel
			select {
			case docsc <- d:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return nil
	})

	// Second goroutine will consume the documents sent from the first and bulk insert into ES
	g.Go(func() error {
		bulk := client.Bulk().Index(*index).Type(*typ)
		count := 0
		for d := range docsc {
			// Enqueue the document
			countstr := strconv.Itoa(count)
			bulk.Add(elastic.NewBulkStringRequest().Id(countstr).SetSource(d))
			count = count + 1
			if bulk.NumberOfActions() >= *bulkSize {
				// Commit
				res, err := bulk.Do(ctx)
				if err != nil {
					return err
				}
				if res.Errors {
					// Look up the failed documents with res.Failed(), and e.g. recommit
					return errors.New("bulk commit failed")
				}
				// "bulk" is reset after Do, so you can reuse it
			}

			select {
			default:
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		// Commit the final batch before exiting
		if bulk.NumberOfActions() > 0 {
			_, err = bulk.Do(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})

	// Wait until all goroutines are finished
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
