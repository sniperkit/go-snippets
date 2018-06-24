package database

import (
	"github.com/rsmohamad/comp4321/models"

	"github.com/boltdb/bolt"
	"sort"
	"strconv"
	"strings"
)

// Class for reading the database
// Reads the .db file in read-only mode.
type Viewer struct {
	db *bolt.DB
}

// Load a Viewer object from .db file
func LoadViewer(filename string) (*Viewer, error) {
	var viewer Viewer
	var err error
	viewer.db, err = bolt.Open(filename, 0666, &bolt.Options{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	return &viewer, nil
}

func (v *Viewer) containsKey(key string, table int) (present bool) {
	v.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(intToByte(table))
		val := b.Get([]byte(key))
		present = val != nil
		return nil
	})
	return
}

func (v *Viewer) stringToId(key string, table int) (rv []byte) {
	rv = nil
	v.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(intToByte(table))
		val := b.Get([]byte(key))
		if val != nil {
			rv = make([]byte, len(val))
			copy(rv, val)
		}
		return nil
	})
	return
}

func (v *Viewer) idToString(key []byte, table int) (rv string) {
	rv = ""
	v.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(intToByte(table))
		val := b.Get(key)
		if val != nil {
			temp := make([]byte, len(val))
			copy(temp, val)
			rv = string(temp)
		}
		return nil
	})
	return
}

func (v *Viewer) idToWord(key []byte) string {
	return v.idToString(key, WordIdToWord)
}

func (v *Viewer) wordToId(word string) []byte {
	return v.stringToId(word, WordToWordId)
}

func (v *Viewer) urlToId(url string) []byte {
	return v.stringToId(url, UrlToPageId)
}

// Check if a URL exist in the database
func (v *Viewer) ContainsUrl(url string) bool {
	return v.containsKey(url, UrlToPageId)
}

// Check if a word exist in the database
func (v *Viewer) ContainsWord(word string) bool {
	return v.containsKey(word, WordToWordId)
}

// Returns a list of page IDs containing a word
func (v *Viewer) GetContainingPages(word string) []uint64 {
	rv := make([]uint64, 0)
	set := make(map[uint64]bool)
	wordId := v.wordToId(word)
	if wordId == nil {
		return rv
	}
	v.db.View(func(tx *bolt.Tx) error {
		wordBucket := tx.Bucket(intToByte(InvertedTable))
		wordBucketTitle := tx.Bucket(intToByte(InvertedTableTitle))
		docBucket := wordBucket.Bucket(wordId)
		docBucketTitle := wordBucketTitle.Bucket(wordId)

		if docBucket != nil {
			docBucket.ForEach(func(k, v []byte) error {
				set[byteToUint64(k)] = true
				return nil
			})
		}

		if docBucketTitle != nil {
			docBucketTitle.ForEach(func(k, v []byte) error {
				set[byteToUint64(k)] = true
				return nil
			})
		}

		return nil
	})

	for key := range set {
		rv = append(rv, key)
	}

	sort.Slice(rv, func(i, j int) bool {
		return rv[i] < rv[j]
	})

	return rv
}

// Returns a document object from a pageId.
// Returns nil if the pageId does not exist
func (v *Viewer) GetDocument(pageId uint64) (document *models.Document) {
	pageIdByte := uint64ToByte(pageId)
	v.db.View(func(tx *bolt.Tx) error {
		documents := tx.Bucket(intToByte(PageInfo))
		docBytes := documents.Get(pageIdByte)

		// pageId does not exist if docBytes == nil
		if docBytes == nil {
			return nil
		}

		document = byteToDoc(docBytes)
		return nil
	})
	return
}

// Returns a list of documents from a list of IDs
func (v *Viewer) GetDocuments(pageIds []uint64) (documents []*models.Document) {
	documents = make([]*models.Document, 0)
	for _, pageId := range pageIds {
		document := v.GetDocument(pageId)
		documents = append(documents, document)
	}
	return
}

func (v *Viewer) GetParentLinks(pageId uint64) []string {
	pageIdByte := uint64ToByte(pageId)
	rv := make([]string, 0)
	v.db.View(func(tx *bolt.Tx) error {
		adjLists := tx.Bucket(intToByte(AdjList))
		idToUrl := tx.Bucket(intToByte(PageIdToUrl))

		parents := adjLists.Bucket(pageIdByte)
		if parents == nil {
			return nil
		}

		parents.ForEach(func(parentId, _ []byte) error {
			linkStr := string(idToUrl.Get(parentId))
			rv = append(rv, linkStr)
			return nil
		})
		return nil
	})

	return rv
}

func (v *Viewer) GetPageRank(pageId uint64) float64 {
	pageIdByte := uint64ToByte(pageId)
	rv := 0.0
	v.db.View(func(tx *bolt.Tx) error {
		pageranks := tx.Bucket(intToByte(PageRank))

		prBytes := pageranks.Get(pageIdByte)
		if prBytes == nil {
			return nil
		}

		rv = byteToFloat64(prBytes)
		return nil
	})

	return rv
}

// Return the positions of a word in a document.
// If the word does not exist in the inverted table, returns an empty slice.
func (v *Viewer) GetPositionIndices(docId uint64, word string, title bool) []uint64 {
	rv := make([]uint64, 0)

	tablename := intToByte(InvertedTable)
	if title {
		tablename = intToByte(InvertedTableTitle)
	}

	wordId := v.wordToId(word)
	if wordId == nil {
		return rv
	}

	v.db.View(func(tx *bolt.Tx) error {
		inv := tx.Bucket(tablename)
		docs := inv.Bucket(wordId)
		if docs == nil {
			return nil
		}

		indices := docs.Get(uint64ToByte(docId))
		if indices == nil {
			return nil
		}

		indicesArr := strings.Split(string(indices), ",")
		for _, indexStr := range indicesArr {
			index, _ := strconv.Atoi(indexStr)
			rv = append(rv, uint64(index))
		}
		return nil
	})

	return rv
}

// Iterate over all documents
func (v *Viewer) ForEachDocument(fn func(p *models.Document, i int)) {
	v.db.View(func(tx *bolt.Tx) error {
		// Get pages bucket
		documents := tx.Bucket(intToByte(PageInfo))
		count := 0

		// Iterate over all documents
		documents.ForEach(func(k, docBytes []byte) error {
			page := byteToDoc(docBytes)
			fn(page, count)
			count++
			return nil
		})
		return nil
	})
}

func (v *Viewer) GetMagnitude(docId uint64, title bool) (rv float64) {
	tablename := intToByte(PageMagnitude)
	if title {
		tablename = intToByte(TitleMagnitude)
	}

	v.db.View(func(tx *bolt.Tx) error {
		tw := tx.Bucket(tablename)
		val := tw.Get(uint64ToByte(docId))
		if val == nil {
			rv = 0
		} else {
			rv = byteToFloat64(val)
		}
		return nil
	})
	return
}

func (v *Viewer) GetTfIdf(docId uint64, word string, title bool) (rv float64) {
	wordId := v.wordToId(word)
	if wordId == nil {
		return 0
	}

	tablename := intToByte(TermWeights)
	if title {
		tablename = intToByte(TitleWeights)
	}

	v.db.View(func(tx *bolt.Tx) error {
		tw := tx.Bucket(tablename)
		words := tw.Bucket(uint64ToByte(docId))
		val := words.Get(wordId)

		if val == nil {
			rv = 0
		} else {
			rv = byteToFloat64(val)
		}
		return nil
	})

	return
}

func (v *Viewer) GetKeywords() []string {
	rv := make([]string, 0)
	v.db.View(func(tx *bolt.Tx) error {
		words := tx.Bucket(intToByte(WordToWordId))
		words.ForEach(func(word, _ []byte) error {
			rv = append(rv, string(word))
			return nil
		})
		return nil
	})
	return rv
}

func (v *Viewer) Close() {
	v.db.Close()
}
