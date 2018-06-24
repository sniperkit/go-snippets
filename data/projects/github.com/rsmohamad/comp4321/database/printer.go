package database

import (
	"fmt"
	"github.com/boltdb/bolt"
)

type Printer struct {
	db *bolt.DB
}

// Load a Viewer object from .db file
func LoadPrinter(filename string) (*Printer, error) {
	var printer Printer
	var err error
	printer.db, err = bolt.Open(filename, 0666, &bolt.Options{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	return &printer, nil
}

func (p *Printer) printAllIDs(tablename int) {
	i := 0
	p.db.View(func(tx *bolt.Tx) error {
		words := tx.Bucket(intToByte(tablename))
		words.ForEach(func(key, val []byte) error {
			i++
			fmt.Println(i, ")", key, string(val))
			return nil
		})
		return nil
	})
}

func (p *Printer) PrintAllPages() {
	p.printAllIDs(PageIdToUrl)
}

func (p *Printer) PrintAllWords() {
	p.printAllIDs(WordIdToWord)
}

func (p *Printer) PrintAdjList() {
	p.db.View(func(tx *bolt.Tx) error {
		idToUrl := tx.Bucket(intToByte(PageIdToUrl))
		adjList := tx.Bucket(intToByte(AdjList))
		adjList.ForEach(func(child, _ []byte) error {
			fmt.Println("Child:", child, string(idToUrl.Get(child)))
			parentList := adjList.Bucket(child)
			parentList.ForEach(func(parent, _ []byte) error {
				fmt.Println(parent, string(idToUrl.Get(parent)))
				return nil
			})
			fmt.Println("----------------------------------------------------------")
			return nil
		})
		return nil
	})
}

func (p *Printer) PrintPageRank() {
	p.db.View(func(tx *bolt.Tx) error {
		idToUrl := tx.Bucket(intToByte(PageIdToUrl))
		prBucket := tx.Bucket(intToByte(PageRank))
		prBucket.ForEach(func(docID, pageRank []byte) error {
			fmt.Println("Document:", string(idToUrl.Get(docID)))
			fmt.Println("PageRank: ", byteToFloat64(pageRank))
			return nil
		})
		return nil
	})
}

func (p *Printer) PrintForwardIndex(title bool) {
	tablename := intToByte(ForwardTable)
	if title {
		tablename = intToByte(ForwardTableTitle)
	}
	p.db.View(func(tx *bolt.Tx) error {
		fw := tx.Bucket(tablename)
		i := 0
		fw.ForEach(func(docID, val []byte) error {
			i++
			fmt.Println(i, ")", docID, val)
			return nil
		})
		return nil
	})
}

func (p *Printer) Close() {
	p.db.Close()
}
