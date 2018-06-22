package main

import (
	"davebalmain.com/noshgobble/src/nutdb"
	"fmt"
	"github.com/blevesearch/bleve"
	"log"
	"strconv"
)

func indexFoodDb() {
	foodMapping := bleve.NewDocumentMapping()
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = "en"
	englishTextFieldMapping.Store = false
	foodMapping.AddFieldMapping(englishTextFieldMapping)

	mapping := bleve.NewIndexMapping()
	mapping.DefaultMapping = foodMapping

	index, err := bleve.New(IndexDirectory, mapping)
	if err != nil {
		log.Fatalf("Error creating Bleve Index %v\n", err)
	}
	batch := bleve.NewBatch()

	for i, f := range nutdb.FoodDb {
		if i%500 == 0 {
			index.Batch(batch)
			batch = bleve.NewBatch()
			fmt.Printf("Indexed %d foods\n", i)
		}
		batch.Index(strconv.FormatInt(int64(f.Id), 10), f)
	}
	index.Batch(batch)
}
