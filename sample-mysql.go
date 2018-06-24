// The Wikipedia articles can be downloaded at https://dumps.wikimedia.org/enwiki/latest/
// Package main provides ...

// TODO: Big picture from here is to make sorting algorithm for pagerank function, update SqlInsert, then the "installation/CreateDB" part of the program is finished.
// Then prettify the output of WriteTXT with a CSS just for the sake of it.
// Then when the program is to utilize the dijkstra's algorithm and have found a path, in the main function, use the pagerank part of the DB to load all articles with a max distance/depth of 7 from the MySQL DB (test it to utilize RAM and/or CPU best, e.g how many articles to load at a time), then discard the neighboring articles from the RAM if they do not have the base article as a neighbor.
// Then get the top 5 or 10 best pageranking neighboring articles from each article in the path (len(path) / 7 * 3, then 5, then 7 if RAM is clocked), then use WriteTXT function to write the summaries of each article (path and top 5 or 10 pageranking algorithms).

package main
import (
	"fmt"
	"io"
	"io/ioutil"
	"compress/bzip2"
	"os"
	// "bufio"
	"sort"
	// "strings"
	"strconv"
	"regexp"
	// "flag"
	// "errors"
	"github.com/dustin/go-wikiparse"
	"github.com/boltdb/bolt"
	// "github.com/arnauddri/algorithms/data-structures/graph"
	"github.com/alixaxel/pagerank"
	"rosettacode/dijkstra" // See RosettaCode, I cannot take credit for this package.
	// "sync"
	// "testing"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

type SuccessorPr map[string][]byte

// TODO: Implement this into func SqlPrepare and make a belonging SqlInsert function.
type PageItems struct {
	// Sections, title indicated by key, item is content/text
	Sections map[string]string

	Data map[string](map[string]sql.Stmt)

	// The NodeID, used for grahping/treeing
	NodeID []byte

	// Weight/distance from root
	Weight uint8

	// Links from this article, used to collect them for the MySQL Db, after that the program will use them to utilize the Db for Dijkstra's algorithm and the Pagerank algorithm.
	Links string

	Pageranks map[uint8]SuccessorPr
}

func main() {
	// createDb := flag.String("create-db", "", "Whether to create db - It has to be to start with, otherwise error will occur.")
	// updateDb := flag.String("update-db", "", "Whether or not to update db")
	// base := flag.Bool("base", false, "A base article, the article used to communicate with other users.")
	// link := flag.Bool("link", false, "A target article, an article and etc. to link by either getting the Bloom-circle from another user or to create it your own pc, after the Bloom-circle is created it will link the your base articles with the new.")
	// flag.Parse()
	
	// var articles map[string]PageItems
	
	dumpFiles, err := GetFilesFromArticlesDir()
	defer recover()
	if err != nil { panic(err) }
	// fmt.Println(dumpFiles)

	db, err := sql.Open(
		"mysql",
		"<username>:<password>@tcp(localhost:3306)/wiki")
	if err != nil { panic(err) }

	// I do not think my OS supports the Prepare method, so execute directly, gobyexample says it does not matter anyway.
	// Also, the boltdb version would probably benefit from a counter too (commiting to disk after 1000 read articles).
	/* StmtCreateSqlTable, err := db.Prepare("CREATE TABLE `?` (NodeID INT, Links CHAR(`?`))")
	if err != nil { panic(err) }
	StmtInsertNodeIDAndLinks, err := db.Prepare("INSERT INTO `?` (NodeID, Links) VALUES (?, ?))")
	if err != nil { panic(err) }
	StmtAddSection, err := db.Prepare("ALTER TABLE `?` ADD `?` blob")
	if err != nil { panic(err) }
	StmtInsertSection, err := db.Prepare("INSERT INTO `?` (`?`) VALUES (`?`)")
	if err != nil { panic(err) } */

	// tx, err := db.Begin()
	// if err != nil { panic(err) }

	// articles := make(map[string]PageItems)
	
	for _, file := range dumpFiles {
		// fmt.Println(file)
		ioReader, err := DecompressBZip(file)
		if err != nil { panic(err) }
		
		wikiParser, err :=  wikiparse.NewParser(ioReader)
		if err != nil { panic(err) }

		var cnt uint16 = 0
		for {
			page, err := wikiParser.Next()
			if err != nil { break }
			cnt++

			// articles[page.Title], err = ReadWikiXML(*page)
			pageItems, err := ReadWikiXML(*page)
			if err != nil { panic(err) }

			// _, err = StmtCreateSqlTable.Exec(page.Title, strconv.Itoa(len(articles[page.Title].Links)))
			// _, err = StmtInsertNodeIDAndLinks.Exec(page.Title, page.ID, articles[page.Title].Links)

			// _, err = StmtCreateSqlTable.Exec(page.Title, len(articles[page.Title].Links))
			// if err != nil { panic(err) }
			// _, err = StmtInsertNodeIDAndLinks.Exec(page.Title, page.ID, articles[page.Title].Links)
			// if err != nil { panic(err) }
			/* _, err = StmtCreateSqlTable.Exec(page.Title, len(pageItems.Links))
			if err != nil { panic(err) }
			_, err = StmtInsertNodeIDAndLinks.Exec(page.Title, page.ID, pageItems.Links)
			if err != nil { panic(err) } */
			// _, err = db.Exec("CREATE TABLE `" + page.Title + "` (NodeID INT, Links BLOB")
			_, err = db.Exec("CREATE TABLE `" + page.Title + "` (NodeID INT, Links BLOB)")
			if err != nil { panic(err) }
			_, err = db.Exec("INSERT INTO `" + page.Title + "` (NodeID, Links) Values (?, ?)", page.ID, pageItems.Links)
			if err != nil { panic(err) }
			for sectionTitle, sectionBody := range pageItems.Sections {
			// for sectionTitle, sectionBody := range articles[page.Title].Sections {
				_, err = db.Exec("ALTER TABLE `" + page.Title + "` ADD `" + sectionTitle + "` BLOB")
				if err != nil { panic(err) }
				_, err = db.Exec("INSERT INTO `" + page.Title + "` (`" + sectionTitle + "`) VALUES (?)", sectionBody)
				if err != nil { panic(err) }

				/* _, err = StmtAddSection.Exec(page.Title, sectionTitle)
				if err != nil { panic(err) }
				_, err = StmtInsertSection.Exec(page.Title, sectionTitle, sectionBody)
				if err != nil { panic(err) } */
				/* _, err = tx.Exec("ALTER TABLE `" + page.Title + "` ADD `" + sectionTitle + "` BLOB")
				_, err = tx.Exec("INSERT INTO `" + page.Title + " (`" + sectionTitle + "`) VALUES `" + sectionBody + "`") */

				// if err := SqlInsert(*db, articles); err != nil { panic(err) }
				// articles = nil
			}
			// if cnt == 1000 {
				// tx.Commit()
				// cnt = 0
				// articles = nil
				// articles = make(map[string]PageItems)
			// }
		}
		
		// if cnt > 0 {
			// tx.Commit()
			// cnt = 0
			// articles = nil
			// articles = make(map[string]PageItems)
		// }
		
		fmt.Println("appended " + file)
	}
	db.Close()

	/* if *base {
		baseFile, err := os.Open("bases.txt")
		if err != nil { panic(err) }
		fileReader := bufio.NewReader(baseFile)
		for err == nil {
			baseArticle, _, err := fileReader.ReadLine()
			if err != nil { panic(err) }
			// wikiGraph, err := BoltGetGraph(string(baseArticle))
			if err != nil { panic(err) }
			articles, err := BoltGetArticles(string(baseArticle))
			if err != nil { panic(err) }
			articles, err = BoltGetChildren(string(baseArticle), articles)
			if err != nil { panic(err) }
			articlesPr, err := PagerankGraph(string(baseArticle), articles[string(baseArticle)].Pageranks)
			if err != nil { panic(err) }
			for depth, item := range articlesPr {
				articles[string(baseArticle)].Pageranks[depth] = item
			}
			if err := BoltInsertPagerank(articles); err != nil { panic(err) }
			err = WriteTxt(articles)
			if err != nil { panic(err) }
		} */
		// Then use the "link" flag to execute the dijkstra function from each created base from earlier. The top ranking shared pages with a max distance of 7 from bases should be returned as well.
		// Use the WriteTxt function for 
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
func DecompressBZip (file string) (ioReader io.Reader, err error) {
	osfile, err := os.Open("/run/media/naamik/Data/articles/" + file)
	if err != nil { return nil, err }
	ioReader = bzip2.NewReader(osfile)
	return
}

// Reads Wikipedia articles from a Wikimedia XML dump bzip file, return the Article with titles as map keys and PageItems (Links, Sections and Text) as items - Also add Section "See Also"
func ReadWikiXML(page wikiparse.Page) (pageItems PageItems, err error) {
	for i := 0; i < len(page.Revisions); i++ {
		// if text is not nil then add to articles text and sections to articles 
		if page.Revisions[i].Text != "" {
			pageItems.Sections = make(map[string]string)
			pageItems.Sections, err = GetSections(page.Revisions[i].Text)
			if err != nil { return pageItems, err }
			for num, link := range wikiparse.FindLinks(page.Revisions[i].Text) {
				if num == 1 {
					pageItems.Links = link
				} else {
					pageItems.Links = pageItems.Links + "-" + link
				}
			}
			if err != nil { return pageItems, err }
			pageItems.NodeID = []byte(strconv.Itoa(int(page.Revisions[i].ID)))
		}
	}
	return
}
// Gets sections from a wikipedia article, page is article content, title is article title
func GetSections(page string) (sections map[string]string, err error) {
	sections = make(map[string]string)
	// Make a regexp search object for section titles
	re, err := regexp.Compile("[=]{2,5}.{1,50}[=]{2,5}")
	if err != nil { return }
	// fmt.Println(title + "\n" + page)
	index := re.FindAllStringIndex(page, -1)
	if len(index) == 0 { return nil, nil }
	// Look at how regex exactly does this
	for i := 0; i < len(index); i++ {
		if i == 0 {
			sections["Summary"] = page[:index[i][0]] // Error: slice out of bounds
		} else if i < len(index)-1 {
			sections[page[index[i][0]:index[i][1]]] = page[index[i][1]:index[i+1][0]] // Assume this will create an error like the "if i == 0" condition error (maybe not)
			// sections[page[index[i][0]:index[i]]] = page[index[i]:index[i+1]]
		} else {
			sections[page[index[i][0]:index[i][1]]] = page[index[i][1]:len(page)] // Assume this will create an error like the "if i == 0" condition error /maybe not)
			// sections[page[index[i]:index[i]]] = page[index[i]:len(page)]
		}
	}
	return
}

// I realized this was wrong, the intend was to make prepared statements. I thought the db.Prepare method returned prepared statements with the variables' content, but it seems the sql driver I use does not do so. Instead of this I will use db.Prepare and create a transaction to insert the variables into the given transaction's statements.
/* func SqlPrepare(db *sql.DB, articles map[string]PageItems) (data map[string](map[string]sql.Stmt), err error) {
	data = make(map[string](map[string]Sql.Stmt))
	for title, items := range articles {
		// sections := make([]string, len(items.Sections))
		data[title]["SetTable"], err = db.Prepare("CREATE TABLE `" + title + "` (NodeID int, Links char(" + strconv.Itoa(len(items.Links)) + ")")
		data[title]["InsertLinks"], err = db.Prepare("INSERT INTO `" + title + "` (NodeID, Links) VALUES (?, ?)", items.NodeID, items.Links)
		if err != nil { return nil, err }
		// TODO: Update to fit with CREATE TABLE
		for sectionTitle, sectionBody := range items.Sections {
			data[title]["AlterTableAdd" + sectionTitle], err = db.Prepare("ALTER TABLE `" + title + "` ADD `" + sectionTitle + "` blob")
			if err != nil { return nil, err }
			data[title].SqlStmtsAndArgs["Insert" + sectionTitle], err = db.Prepare("INSERT INTO `" + title + "` (" + sectionTitle + ") VALUES (?)", sectionBody)
			if err != nil { return nil, err }
		}
	}
	return sqlStmts, nil
} */

// This is not needed, although the function is very much reusable
/* func SqlInsert(tx *sql.Tx, articles map[string]PageItems) (err error) {
	for title, items := range articles {
		// sections := make([]string, len(items.Sections))
		sqlStmts[title]["tree"], err = tx.Exec("CREATE TABLE `" + title + "` (NodeID int, Links char(" + strconv.Itoa(len(items.Links)) + ")")
		_, err = db.Prepare("INSERT INTO `" + title + "` (NodeID, Links) VALUES (?, ?)", items.NodeID, items.Links)
		if err != nil { return err }
		// TODO: Update to fit with CREATE TABLE
		for sectionTitle, sectionBody := range items.Sections {
			sqlStmts[title][sectionTitle], err = tx.Exec("ALTER TABLE `" + title + "` ADD `" + sectionTitle + "` blob")
			if err != nil { return err }
			_, err = tx.Exec("INSERT INTO `" + title + "` (" + sectionTitle + ") VALUES (?)", sectionBody)
			if err != nil { return err }
		}
	}
	return sqlStmts, nil
} */

/* func SqlSelectArticles(db sql.DB, table []string) (articles map[string]PageItems, err error) {
	var tmp PageItems
	for _, title := range table {
		rows, err := db.Query("SELECT * FROM " + title);
		if err != nil { return nil, err }
		columns, err := rows.Columns()
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for rows.Next() {
			for i, _ := range columns {
				valuePtrs[i] = &values[i]
			}
			rows.Scan(valuePtrs...)
			
			for i, _ := range columns {
				tmp = articles[title]
				switch columns[i] {
				case "NodeID":
					if values[i] != nil {
						tmp.NodeID = []byte(values[i].(string))
						if err != nil { return nil, err }
					}
				// TODO: Update to fit with PageItems.Pageranks - Also update PagerankGraph to your original idea (better ranking or the tree already there, but also relative to vertices out of bounds of root, but still related to root).
				case "Pagerank_1":
					if values[i] != nil {
						fmt.Println("case \"Pagerank\": columns[i]:\t", columns[i], values[i])
						tmp.Pagerank[0], err = strconv.ParseFloat(values[i].(string), 32)
						if err != nil { return nil, err }
					}
				case "Pagerank_2":
					if values[i] != nil {
						fmt.Println("case \"Pagerank\": columns[i]:\t", columns[i], values[i])
						tmp.Pagerank[1], err = strconv.ParseFloat(values[i].(string), 32)
						if err != nil { return nil, err }
					}
				case "Pagerank_3":
					if values[i] != nil {
						fmt.Println("case \"Pagerank\": columns[i]:\t", columns[i], values[i])
						tmp.Pagerank[2], err = strconv.ParseFloat(values[i].(string), 32)
						if err != nil { return nil, err }
					}
				case "Pagerank_4":
					if values[i] != nil {
						fmt.Println("case \"Pagerank\": columns[i]:\t", columns[i], values[i])
						tmp.Pagerank[3], err = strconv.ParseFloat(values[i].(string), 32)
						if err != nil { return nil, err }
					}
				case "Pagerank_5":
					if values[i] != nil {
						fmt.Println("case \"Pagerank\": columns[i]:\t", columns[i], values[i])
						tmp.Pagerank[4], err = strconv.ParseFloat(values[i].(string), 32)
						if err != nil { return nil, err }
					}
				case "Pagerank_6":
					if values[i] != nil {
						fmt.Println("case \"Pagerank\": columns[i]:\t", columns[i], values[i])
						tmp.Pagerank[5], err = strconv.ParseFloat(values[i].(string), 32)
						if err != nil { return nil, err }
					}
				case "Pagerank_7":
					if values[i] != nil {
						fmt.Println("case \"Pagerank\": columns[i]:\t", columns[i], values[i])
						tmp.Pagerank[6], err = strconv.ParseFloat(values[i].(string), 32)
						if err != nil { return nil, err }
					}
				default:
					if values[i] != nil {
						fmt.Println("default: columns[i]:\t", columns[i], values[i])
						tmp.Sections[columns[i]] = values[i].(string)
					}
				}
				articles[title] = tmp
			}
		}
	}
	return articles, nil
} */

// TODO: Comment SqlUpdate and SqlDelete out or update them.
/* func SqlUpdate(table, column string, value map[string]interface{}) (query string, err error) {
	query, _, err = squirrel.Update(table).
    Set(column, value).
		ToSql()
	if err != nil { return "", err }
	return query, nil
}
func SqlDelete(table, column, value string) (query string, err error) {
	query, _, err = squirrel.Delete(value).
		From(table).
		ToSql()
	if err != nil { return "", err }
	return query, nil
} */
func PagerankGraph(title string, children map[uint8]SuccessorPr) (map[uint8]SuccessorPr, error) {
	articlesDb, err := bolt.Open("/home/naamik/go/wikiproj/articles.boltdb", 0666, nil)
	if err != nil { return nil, err }
	articlesTx, err := articlesDb.Begin(false)
	if err != nil { return nil, err }
	pagerankGraph := pagerank.NewGraph()
	articlesTitle := make(map[string]string)
	articlesDepth := make(map[string]uint8)
	// TODO: add vertices out of bounds of root (get vertices 7 - depth from "absNeighbor")
	for depth := uint8(1); depth < uint8(7); depth += 2 {
		for neighbor, _ := range children[depth - 1] {
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
func (a SortedPageranks) Len() int { return len(a) }
func (a SortedPageranks) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortedPageranks) Less(i, j int) bool { return a[i] < a[j] }
func WriteTxt(articles map[string]PageItems) (err error) {
	var pageranks map[string][]Pagerank
	// fWriter := bufio.NewWriter(ioWriter)
	for articleName, _ := range(articles) {
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

