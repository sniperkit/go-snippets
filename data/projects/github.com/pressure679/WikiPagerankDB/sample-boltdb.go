// The Wikipedia articles can be downloaded at https://dumps.wikimedia.org/enwiki/latest/
// Package main provides ...

// TODO: Big picture from here is to make sorting algorithm for pagerank function, update SqlInsert, then the "installation/CreateDB" part of the program is finished.
// Then prettify the output of WriteTXT with a CSS just for the sake of it.
// Then when the program is to utilize the dijkstra's algorithm and have found a path, in the main function, use the pagerank part of the DB to load all articles with a max distance/depth of 7 from the MySQL DB (test it to utilize RAM and/or CPU best, e.g how many articles to load at a time), then discard the neighboring articles from the RAM if they do not have the base article as a neighbor.
// Then get the top 5 or 10 best pageranking neighboring articles from each article in the path (len(path) / 7 * 3, then 5, then 7 if RAM is clocked), then use WriteTXT function to write the summaries of each article (path and top 5 or 10 pageranking algorithms).

package main

import (
	"compress/bzip2"
	"fmt"
	"io"
	"io/ioutil"
	// "bufio"
	"bytes"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"github.com/boltdb/bolt"
	"github.com/dustin/go-wikiparse"
	"github.com/alixaxel/pagerank"
	"rosettacode/dijkstra" // See RosettaCode, I cannot take credit for this package.
	"encoding/gob"
)

type SuccessorPr map[string][]byte
type PageItems struct {
	// Sections, title indicated by key, item is content/text
	Sections map[string]string

	// The NodeID, used for grahping/treeing
	NodeID []byte

	// Weight/distance from root
	Weight uint8

	// Links from this article, used to collect them for the MySQL Db, after that the program will use them to utilize the Db for Dijkstra's algorithm and the Pagerank algorithm.
	Links string

	Pageranks map[uint8]SuccessorPr
}
type BoltDBContent struct {
	Sections map[string]string
	Links string
}
type BoltDBInput struct {
	Title []byte
	Links []byte
	Sections []byte
}

func main() {
	// createDb := flag.String("create-db", "", "Whether to create db - It has to be to start with, otherwise error will occur.")
	// updateDb := flag.String("update-db", "", "Whether or not to update db")
	// base := flag.Bool("base", false, "A base article, the article used to communicate with other users.")
	// link := flag.Bool("link", false, "A target article, an article and etc. to link by either getting the Bloom-circle from another user or to create it your own pc, after the Bloom-circle is created it will link the your base articles with the new.")
	// flag.Parse()

	defer recover()
	dumpFiles, err := GetFilesFromArticlesDir()
	if err != nil { panic(err) }
	// fmt.Println(dumpFiles)
	var articlesDb *bolt.DB
	var articlesTx *bolt.Tx
	var cnt uint16 = 0
	articlesDb, err = bolt.Open("/run/media/naamik/Data/wiki.boltdb", 0666, nil)
	if err != nil { panic(err) }
	// articlesDb.AllocSize = int(allocSize * 6)
	articlesTx, err = articlesDb.Begin(true)
	// enc := gob.Encoder(&string)
	if err != nil {
		fmt.Println("articlesDb.Begin() error")
		panic(err)
	}
	var size uint = 0
	var contentBuffer bytes.Buffer
	var titleBuffer bytes.Buffer
	contentEncoder := gob.NewEncoder(&contentBuffer)
	titleEncoder := gob.NewEncoder(&titleBuffer)
	for _, file := range dumpFiles {
		articlesTx, err = articlesDb.Begin(true)
		articles := make(map[string][]byte)
		
		fmt.Println("Reading", file)
		// fmt.Println(file)
		ioReader, allocSize, err := DecompressBZip(file)
		// articlesTx.DB().AllocSize = int(float64(allocSize) * 6)
		if err != nil {
			fmt.Println("ioReader, DecompressBZip error")
			panic(err)
		}
		wikiParser, err := wikiparse.NewParser(ioReader)
		if err != nil {
			fmt.Println("wikiParser, wikiparser.NewParser error")
			panic(err)
		}
		// offset := strings.Index(file, ".xml")
		for err == nil {
			page, err := wikiParser.Next()
			if err != nil {
				if strings.EqualFold(err.Error(), "EOF") {
					fmt.Println("wikiParser.Next() error == EOF")
					break
				} else {
					fmt.Println("wikiParser.Next() error != EOF")
					continue
				}
			}
			// articles[page.Title], err = ReadWikiXML(*page)
			// title, err := enc.Encode(page.Title)
			err = ReadWikiXML(*encoder, *page)
			if err != nil { panic(err) }
			size += uint(len(articles[page.Title]))
		}
		// articlesTx.AllocSize = size
		if err := BoltInsertArticles(articlesTx, title, content); err != nil {
			fmt.Println("BoltInsertArticles error")
			panic(err)
		}
		/* cnt++
			if cnt == 1000 {
				if err := BoltInsertArticles(articlesTx, articles, -1); err != nil {
					fmt.Println("BoltInsertArticles error")
					panic(err)
				}
				articles = nil
				articles = make(map[string]PageItems)
				cnt = 0
				articlesTx.Commit()
				articlesTx, err = articlesDb.Begin(true)
			}
			// fmt.Println("put in", page.Title)
		}
		if cnt > 0 {
			if err := BoltInsertArticles(articlesTx, articles, -1); err != nil {
				fmt.Println("BoltInsertArticles error")
				panic(err)
			}
			articles = nil
			articles = make(map[string]PageItems)
			cnt = 0
			articlesTx.Commit()
			articlesTx, err = articlesDb.Begin(true)
		} */
		articles = nil
		articles = make(map[string][]byte)
		articlesTx.Commit()
	}
	articlesDb.Close()
}

	// Read all names of Bzipped Wikimedia XML files from "articles" dir.
func GetFilesFromArticlesDir() (files []string, err error) {
	osFileInfo, err := ioutil.ReadDir("/run/media/naamik/Data/articles")
	if err != nil { return nil, err }
	for _, fileInfo := range osFileInfo {
		if !fileInfo.IsDir() {
			files = append(files, fileInfo.Name())
		}
	}
	return
}

// uses os.Open to make an io.Reader from bzip2.NewReader(os.File) to read wikipedia xml file
func DecompressBZip(file string) (ioReader io.Reader, fileSize int64, err error) {
	osFile, err := os.Open("/run/media/naamik/Data/articles/" + file)
	if err != nil { return nil, -1, err }
	fileStat, err := osFile.Stat()
	if err != nil { return nil, fileStat.Size(), err }
	ioReader = bzip2.NewReader(osFile)
	return ioReader, fileStat.Size(), nil
}

// Reads Wikipedia articles from a Wikimedia XML dump bzip file, return the Article with titles as map keys and PageItems (Links, Sections and Text) as items - Also add Section "See Also"
func ReadWikiXML(titleEncoder, contentEncoder gob.Encoder, page wikiparse.Page) (err error) {
	var pageItems BoltDBContent
	for i := 0; i < len(page.Revisions); i++ {
		// if text is not nil then add to articles text and sections to articles
		if page.Revisions[i].Text != "" {
			pageItems.Sections = make(map[string]string)
			pageItems.Sections, err = GetSections(page.Revisions[i].Text, page.Title)
			if err != nil { return err }
			links := wikiparse.FindLinks(page.Revisions[i].Text)
			for num, link := range links {
				pageItems.Links += link
				if num < len(pageItems.Links) { pageItems.Links += "-" }
			}

			if err != nil { return err }
			// pageItems.NodeID = []byte(strconv.Itoa(int(page.Revisions[i].ID)))
		}
		contentEncoder.Encode(BoltDBContent{pageItems.Sections, pageItems.Links})
		titleEncoder.Encode(page.Title)
	}
	return
}

// Gets sections from a wikipedia article, page is article content, title is article title
func GetSections(page, title string) (sections map[string]string, err error) {
	sections = make(map[string]string)
	// Make a regexp search object for section titles
	re, err := regexp.Compile("[=]{2,5}.{1,50}[=]{2,5}")
	if err != nil { return }
	index := re.FindAllStringIndex(page, -1)
	for i := 0; i < len(index)-1; i++ {
		if index[i] != nil {
			if i == 0 {
				sections["Summary"] = page[:index[i][1]-1] // Error: slice out of bounds at article "Southern Hemisphere"
			} else if i < len(index)-1 {
				sections[page[index[i][0]:index[i][1]]] = page[index[i][1]:index[i+1][0]] // Assume this will create an error like the "if i == 0" condition error (maybe not)
			} else {
				sections[page[index[i][0]:index[i][1]]] = page[index[i][1]:len(page)] // Assume this will create an error like the "if i == 0" condition error /maybe not)
			}
		}
	}
	return
}

/* func BoltInsertArticles(articlesTx *bolt.Tx, articles map[string]PageItems, allocSize int) (err error) {
	var links string
	for key, _ := range articles {
		if err != nil { break }
		bucket, err := articlesTx.CreateBucketIfNotExists([]byte(key))
		if err != nil { return err }
		for sectionKey, sectionText := range articles[key].Sections {
			if err := bucket.Put([]byte(sectionKey), []byte(sectionText)); err != nil {
				fmt.Println("BoltInsertArticles, sections input error:", err)
				continue
			}
		}
		for num, link := range articles[key].Links {
			links = links + link
			if num < len(articles[key].Links) { links = links + "-" }
		}
		if err := bucket.Put([]byte("Links"), []byte(links)); err != nil {
			fmt.Println("BoltInsertVertices link input error:", err)
			return err
		}
	}
	// articlesTx.Commit()
	return nil
} */
func BoltInsertArticles(articlesTx *bolt.Tx, articles []BoltDBInput) (err error) {
	var buffer bytes.Buffer
	var encoder gob.NewEncoder(&buffer)
	for _, article := range articles {
		if err != nil { break }
		bucket, err := articlesTx.CreateBucketIfNotExists(article.Title)
		if err != nil { return err }
		bucket.Put(articles[key])
	}
	// articlesTx.Commit()
	return nil
}
func BoltInsertVertices(graphTx *bolt.Tx, articles map[string]PageItems, allocSize int) (err error) {
	var links string
	for key, _ := range articles {
		graphBucket, err := graphTx.CreateBucketIfNotExists([]byte(key)) // Creating buckets for each article makes a bottleneck in the read/write operation.
		if err != nil { panic(err) }
		for num, link := range articles[key].Links {
			links = links + link
			if num < len(articles[key].Links) {
				links = links + "-"
			}
		}
		if !strings.EqualFold(links, "") { // debugging purposes
			if err := graphBucket.Put([]byte("Links"), []byte(links)); err != nil {
				fmt.Println("BoltInsertVertices link input error:", err)
				return err
			}
		}
		if err := graphBucket.Put([]byte("NodeID"), []byte(articles[key].NodeID)); err != nil {
			fmt.Println("BoltInsertVertices nodeid input error:", err)
			return err
		}
	}
	// graphTx.Commit()
	return nil
}
// Yeah this needs to be updates to fit a better db schema,,
/* func BoltInsertPagerank(graphTx *bolt.Tx, articles map[string]PageItems) (err error) {
	var items []byte
	b, err := graphTx.CreateBucketIfNotExists([]byte("Pageranks")) // Creating buckets for each article makes a bottleneck in the read/write operation.
	if err != nil { if strings.EqualFold(err.Error(), "bucket already exists") {
			
		} else {
			return err
		}
	}
	for key, _ := range articles {
		for distance, successors := range articles[key].Pageranks {
			for title, pagerank := range successors {
				items = []byte(title + "-" + string(strconv.Itoa(int(distance))) + "-" + string(pagerank))
			}
		}
		if err = b.Put([]byte(key), items); err != nil { return err }
	}
	return nil
} */

// Gets graph successors, with sections
func BoltGetArticles(articlesTx *bolt.Tx, requests []string) (articles map[string]PageItems, err error) {
	articles = make(map[string]PageItems)
	var tmp PageItems
	for _, request := range requests {
		cursor := articlesTx.Bucket([]byte(request)).Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			tmp.Sections[string(key)] = string(value)
		}
		articles[request] = tmp
	}
	return articles, nil
}
func BoltGetChildren(articlesTx *bolt.Tx, rootTitle string, root map[string]PageItems) (map[string]PageItems, error) {
	for _, rootLink := range root[rootTitle].Links {
		oneLinks := strings.Split(string(articlesTx.Bucket([]byte(rootLink)).Get([]byte("Links"))), "-")
		root[rootTitle].Pageranks[0] = SuccessorPr{rootLink: nil}
		for _, oneLink := range oneLinks {
			twoLinks := strings.Split(string(articlesTx.Bucket([]byte(oneLink)).Get([]byte("Links"))), "-")
			root[rootLink].Pageranks[1] = SuccessorPr{oneLink: nil}
			for _, twoLink := range twoLinks {
				threeLinks := strings.Split(string(articlesTx.Bucket([]byte(twoLink)).Get([]byte("Links"))), "-")
				root[rootLink].Pageranks[2] = SuccessorPr{twoLink: nil}
				for _, threeLink := range threeLinks {
					fourLinks := strings.Split(string(articlesTx.Bucket([]byte(threeLink)).Get([]byte("Links"))), "-")
					root[rootLink].Pageranks[3] = SuccessorPr{threeLink: nil}
					for _, fourLink := range fourLinks {
						fiveLinks := strings.Split(string(articlesTx.Bucket([]byte(fourLink)).Get([]byte("Links"))), "-")
						root[rootLink].Pageranks[4] = SuccessorPr{fourLink: nil}
						for _, fiveLink := range fiveLinks {
							sixLinks := strings.Split(string(articlesTx.Bucket([]byte(fiveLink)).Get([]byte("Links"))), "-")
							root[rootLink].Pageranks[5] = SuccessorPr{fiveLink: nil}
							for _, sixLink := range sixLinks {
								root[rootLink].Pageranks[6] = SuccessorPr{sixLink: nil}
							}
						}
					}
				}
			}
		}
	}
	return root, nil
}

// Gets graph successors, not sections though.
/* func BoltGetGraph(request string) (wikiGraph map[string]*graph.DirGraph, err error) {
	wikiGraph = make(map[string]*graph.DirGraph)
	graphDb, err := bolt.Open("/home/naamik/go/wikiproj/graph+pagerank.boltdb", 0666, nil)
	if err != nil { return nil, err }
	var nodeLink [6]uint64
	var err0 error
	if err = graphDb.View(func(tx *bolt.Tx) error {
		wikiGraph[request] = graph.NewDirected()
		oneLinks := strings.Split(string(tx.Bucket([]byte(request)).Get([]byte("Links"))), "-")
		for _, oneLink := range oneLinks {
			nodeLink[0], err0 = strconv.ParseUint(string((tx.Bucket([]byte(oneLink)).Get([]byte("NodeID")))), 10, 64)
			if err0 != nil { return err0 }
			err0 = wikiGraph[request].AddVertex(graph.VertexId(uint(nodeLink[0])))
			if err0 != nil { return err0 }
			twoLinks := strings.Split(string(tx.Bucket([]byte(oneLink)).Get([]byte("Links"))), "-")
			for _, twoLink := range twoLinks {
				nodeLink[1], err0 = strconv.ParseUint(string((tx.Bucket([]byte(twoLink)).Get([]byte("NodeID")))), 10, 64)
				if err0 != nil { return err0 }
				err0 = wikiGraph[request].AddVertex(graph.VertexId(uint(nodeLink[1])))
				if err0 != nil { return err0 }
				err0 = wikiGraph[request].AddEdge(graph.VertexId(uint(nodeLink[0])), graph.VertexId(uint(nodeLink[1])), 1)
				if err0 != nil { return err0 }
				threeLinks := strings.Split(string(tx.Bucket([]byte(twoLink)).Get([]byte("Links"))), "-")
				for _, threeLink := range threeLinks {
					nodeLink[2], err0 = strconv.ParseUint(string((tx.Bucket([]byte(threeLink)).Get([]byte("NodeID")))), 10, 64)
					if err0 != nil { return err0 }
					err0 = wikiGraph[request].AddVertex(graph.VertexId(uint(nodeLink[2])))
					if err0 != nil { return err0 }
					err0 = wikiGraph[request].AddEdge(graph.VertexId(uint(nodeLink[1])), graph.VertexId(uint(nodeLink[2])), 1)
					if err0 != nil { return err0 }
					fourLinks := strings.Split(string(tx.Bucket([]byte(threeLink)).Get([]byte("Links"))), "-")
					for _, fourLink := range fourLinks {
						nodeLink[3], err0 = strconv.ParseUint(string((tx.Bucket([]byte(fourLink)).Get([]byte("NodeID")))), 10, 64)
						if err0 != nil { return err0 }
						err0 = wikiGraph[request].AddVertex(graph.VertexId(uint(nodeLink[3])))
						if err0 != nil { return err0 }
						err0 = wikiGraph[request].AddEdge(graph.VertexId(uint(nodeLink[2])), graph.VertexId(uint(nodeLink[3])), 1)
						if err0 != nil { return err0 }
						fiveLinks := strings.Split(string(tx.Bucket([]byte(fourLink)).Get([]byte("Links"))), "-")
						for _, fiveLink := range fiveLinks {
							nodeLink[4], err0 = strconv.ParseUint(string((tx.Bucket([]byte(fiveLink)).Get([]byte("NodeID")))), 10, 64)
							if err0 != nil { return err0 }
							err0 = wikiGraph[request].AddVertex(graph.VertexId(uint(nodeLink[4])))
							if err0 != nil { return err0 }
							err0 = wikiGraph[request].AddEdge(graph.VertexId(nodeLink[3]), graph.VertexId(nodeLink[4]), 1)
							if err0 != nil { return err0 }
							sixLinks := strings.Split(string(tx.Bucket([]byte(fiveLink)).Get([]byte("Links"))), "-")
							for _, sixLink := range sixLinks {
								nodeLink[5], err0 = strconv.ParseUint(string((tx.Bucket([]byte(sixLink)).Get([]byte("NodeID")))), 10, 64)
								if err0 != nil { return err0 }
								err0 = wikiGraph[request].AddVertex(graph.VertexId(uint(nodeLink[5])))
								if err0 != nil { return err0 }
								err0 = wikiGraph[request].AddEdge(graph.VertexId(nodeLink[4]), graph.VertexId(nodeLink[5]), 1)
								if err0 != nil { return err0 }
							}
						}
					}
				}
			}
		}
		return nil
	}); err != nil { return nil, err }
	if err = graphDb.Close(); err != nil { return nil, err }
	return wikiGraph, nil
} */
func BoltGetPagerank(request string, articles map[string]PageItems) (map[string]PageItems, error) {
	graphDb, err := bolt.Open("graph+pagerank.boltdb", 0666, nil)
	if err != nil { return nil, err }
	/* graphTx, err := graphDb.Begin(false)
	if err != nil { return nil, err } */
	if err = graphDb.View(func(graphTx *bolt.Tx) error {
		bCursor := graphTx.Bucket([]byte(request)).Cursor()
		for bKey, bValue := bCursor.First(); bKey != nil; bKey, bValue = bCursor.Next() {
			if strings.EqualFold(string(bKey[:8]), "pagerank") {
				strSplit := strings.Split(string(bKey), "-")
				depth, err := strconv.Atoi(strSplit[2])
				if err != nil { return err }
				articles[request].Pageranks[uint8(depth)] = SuccessorPr{strSplit[1]: bValue}
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	if err = graphDb.Close(); err != nil { return nil, err }
	return articles, nil
}
func PagerankGraph(title string, children map[uint8]SuccessorPr) (map[uint8]SuccessorPr, error) {
	articlesDb, err := bolt.Open("articles.boltdb", 0666, nil)
	if err != nil { return nil, err }
	articlesTx, err := articlesDb.Begin(false)
	if err != nil { return nil, err }
	pagerankGraph := pagerank.NewGraph()
	articlesTitle := make(map[string]string)
	articlesDepth := make(map[string]uint8)
	for depth := uint8(1); depth < uint8(7); depth += 2 {
		for neighbor, _ := range children[depth-1] {
			byteNeighborNodeID, err := strconv.Atoi(string(articlesTx.Bucket([]byte(neighbor)).Get([]byte("NodeID"))))
			if err != nil { return nil, err }
			articlesTitle[string(byteNeighborNodeID)] = neighbor
			articlesDepth[string(byteNeighborNodeID)] = depth - 1
			for article, _ := range children[depth] {
				byteArticleNodeID, err := strconv.Atoi(string(articlesTx.Bucket([]byte(article)).Get([]byte("NodeID"))))
				if err != nil { return nil, err }
				articlesTitle[string(byteArticleNodeID)] = article
				articlesDepth[string(byteArticleNodeID)] = depth
				neighborNodeID, err := strconv.Atoi(string(byteArticleNodeID))
				if err != nil { return nil, err }
				articleNodeID, err := strconv.Atoi(string(byteNeighborNodeID))
				if err != nil { return nil, err }
				pagerankGraph.Link(uint32(articleNodeID), uint32(neighborNodeID), 1)
			}
		}
	}
	pr := make(map[uint8]SuccessorPr)
	pagerankGraph.Rank(0.85 /*put damping factor here or just settle with weighing the graph?*/, 0.000001 /*precision*/, func(node uint32, rank float64) {
		bufferRank := []byte(strconv.FormatFloat(rank, 10, 6, 64))
		pr[articlesDepth[articlesTitle[strconv.Itoa(int(node))]]] = SuccessorPr{articlesTitle[strconv.Itoa(int(node))]: bufferRank}
	})
	return pr, nil
}

// graphs the articles, preferred input is a graph with 7 a distance/number of node generations of 7.
// TODO
func Dijkstra(request dijkstra.Request) (path []string) {
	return dijkstra.Get(request)
}

type Pagerank float64
type SortedPageranks []Pagerank

func (a SortedPageranks) Len() int           { return len(a) }
func (a SortedPageranks) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortedPageranks) Less(i, j int) bool { return a[i] < a[j] }
func WriteTxt(articles map[string]PageItems) (err error) {
	var pageranks map[string][]Pagerank
	for articleName, _ := range articles {
		file, err := os.Create(articleName + ".org")
		if err != nil { return err }
		defer file.Close()
		indexFile, err := os.Create("index-" + articleName + ".org")
		if err != nil { return err }
		defer indexFile.Close()
		file.WriteString("* " + articleName)
		file.WriteString("** Sections")
		for sectionName, sectionText := range articles[articleName].Sections {
			file.WriteString("*** " + sectionName)
			file.WriteString("    " + sectionText)
		}
		for depth, _ := range articles[articleName].Pageranks {
			indexFile.WriteString("* " + strconv.Itoa(int(depth)))
			pagerankIndex := make(map[float64]string)
			for article, pagerank := range articles[articleName].Pageranks[depth] {
				float64Pr, err := strconv.ParseFloat(string(pagerank), 64)
				if err != nil { return err }
				pageranks[articleName] = append(pageranks[articleName], Pagerank(float64Pr))
				pagerankIndex[float64Pr] = article
			}
			sort.Sort(sort.Reverse(SortedPageranks(pageranks[articleName])))
			for pagerank, article := range pagerankIndex {
				indexFile.WriteString("** " + article + " - " + strconv.FormatFloat(pagerank, 'f', 6, 64))
			}
		}
	}
	return nil
}

// Words: textcat/n-gram, bayesian+tfidf, pos-tagger(advanced-logic's freeling?)
// TODO: make a basic chatbot which uses Eliza's principle of chatting, the basic outrule for chatting uses Bloom's Taxonomical levels to present data/chat with a user.
// The packages to be used here are word2sentence etc., elasticsearch, and it should utilize the information written to a file in emacs-org format/html-format.
// From here the chatbot is to be written, just to make it J.A.R.V.I.S-like. the chatbot will use ELIZA like chat functionality, e.g explain how one article relates to another by using Bloom's Taxonomical levels (here word2sentence etc. and Elasticsearch will be used).
// The chatbot will load the summaries of the written (from WriteTXT) articles that are in the path of article A, B, C... and their top 5 or 10 pageranking articles.
// The chatbot may start with "Hi, which topics would you like me to chat with you about?". And the chatbot may want to only use a max of 4-7 new words in each message to not strain the short-term memory. It's functionality will be different than the product of WriteTXT function; it will connect the paths between the topics, but the pagerank functionality will work somehow indifferable; IMPORTANT: it will sum up pageranks if a neighboring article has 2 base articles as neighbors. (should maybe also be used for the WriteTXT function).
// Create a neural network for the chatbot with a light prediction analysis algorithm (prediction algorithm trained by user's behaviour) to make it smarter.. Maybe the neural network is not needed since Dijkstra and Pagerank makes up for it; the shortest part and the top pageranking path. Maybe add a neural network to find the path from shortest path and top pageranking path to make it better to chat with humans rather than machines(?).
// Here I may also add a xml dump file of articles from HowStuffWorks if it provides such to add a functinality to the chatbot to answer "how" questions.
// TODO: if the chatbot functionality works well; hook it up to 1 or more irc channels.
// TODO: make a GUI for Android for the chatbot functionality
