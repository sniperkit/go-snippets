package main
import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"strconv"
	"regexp"
	"github.com/pressure679/go-wikiparse"
	"github.com/mvryan/fasttag"
	"io"
	"bufio"
	// "sync"
	"runtime"
	// "time"
)
type Sentence struct { Start, End []byte }
type Sentences []Sentence
type ArticleElements struct {
	Title string
	// Links string
	// Words map[string]Sentences
	Offset int64
}
type OrderArticles map[string][]ArticleElements
func GetFilesFromArticlesDir(directory string) (files []string, err error) {
	osFileInfo, err := ioutil.ReadDir(directory); if err != nil { return nil, err }
	for _, fileInfo := range osFileInfo { if !fileInfo.IsDir() { files = append(files, fileInfo.Name()) } }; return files, err
}
func WriteIndex(myMap map[string]OrderArticles, dir string) error {
	// var indexFile *os.File
	for file, _ := range myMap {
		for char, articles := range myMap[file] {
			// fmt.Println(dir+char+"/"+file+"-"+char+".index")
			indexFile, err := os.OpenFile(dir+char+"/"+file+"-"+char+".index", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0770)
			if err != nil { return err }
			for _, article := range articles {
				// fmt.Fprintln(indexFile, article.Title + ":" + strconv.Itoa(int(article.Offset)) + ":links:" + article.Links)
				fmt.Fprintln(indexFile, article.Title + "&:&" + strconv.Itoa(int(article.Offset)))
			}
			err = indexFile.Close()
			if err != nil { return err }
		}
		// fmt.Println("written", file, "in WriteIndex(...)")
	}
	return nil
}
func GetWords(page string) (wordMap map[string]Sentences) {
	var sentenceIndex, previousSentenceIndex, curNounIndex int = -1, -1, -1
	const (
		sentenceStop = ":,.!?"
	)
	wordMap = make(map[string]Sentences)
	words := fasttag.WordsToSlice(page)
	posTags := fasttag.BrillTagger(words)
	for posNum, _ := range posTags {
		if string(posTags[posNum]) != "N" { continue }
		curNounIndex = strings.Index(page[curNounIndex:], words[posNum])
		if curNounIndex != -1 {
			for i := 0; i < 3; i++ {
				previousSentenceIndex, sentenceIndex = strings.LastIndexAny(page[previousSentenceIndex:curNounIndex], sentenceStop), strings.IndexAny(page[curNounIndex+len(words[posNum])+sentenceIndex+1:], sentenceStop)
			}
			// article.Nouns[*graph.HashTable.Words].Sentences.appendSentence(Sentence{Start: []byte(strconv.Itoa(previousSentenceIndex)), End: []byte(strconv.Itoa(sentenceIndex))})
			wordMap[words[curNounIndex]] = append(wordMap[words[curNounIndex]], Sentence{Start: []byte(strconv.Itoa(previousSentenceIndex)), End: []byte(strconv.Itoa(sentenceIndex))})
		}
	}
	// fmt.Println(wordMap)
	return wordMap
}
func ReadMem() (float64, error) {
	var memTF []string
	memReader, err := os.Open("/proc/meminfo")
	if err != nil { return 0, err }
	bufioReader := bufio.NewReader(memReader)
	if err != nil { return 0, err }
	numbers := regexp.MustCompile(`\d+`)
	for i := 0; i < 2; i++ {
		line, _, err := bufioReader.ReadLine()
		if err != nil { return 0, err }
		mem := numbers.Find(line)
		memTF = append(memTF, string(mem))
	}
	total, err := strconv.Atoi(memTF[0])
	if err != nil { return 0, err }
	free, err := strconv.Atoi(memTF[1])
	if err != nil { return 0, err }
	return float64(free) / float64(total), nil
}
func main() {
	var chars string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var memstats runtime.MemStats
	var myMap map[string]OrderArticles = make(map[string]OrderArticles)
	dir := "/home/naamik/Documents/wikipedia/xml_articles/"
	indexDir := "/home/naamik/Documents/wikipedia/index/"
	dumpFiles, err := GetFilesFromArticlesDir(dir)
	if err != nil { panic(err) }
	var article ArticleElements
	var char string
	var counter uint = 100
	var buffer [2]int64
	for _, file := range dumpFiles {
		ioReader, err := os.Open(dir + file)
		if err != nil { panic(err) }
		wikiParser, err :=  wikiparse.NewParser(ioReader)
		if err != nil { panic(err) }
		myMap[file] = make(map[string][]ArticleElements)
		buffer[0] = int64(0)
		for {
			// wg.Wait()
			if _, ok := myMap[file]; !ok { myMap[file] = make(map[string][]ArticleElements) }
			page, err := wikiParser.Next()
			if err != nil {
				if err == io.EOF { break }
				panic(err) }
			if len(page.Title) == 0 { continue }
			for i := 0; i < len(chars); i++ { if strings.EqualFold(string(page.Title[0]), string(chars[i])) { char = string(chars[i]) } else { char = "etc" } // "etc" added, I noticed the wikipedia articles with special start characters were not added (highly likely in foreign languages).
			/* if strings.EqualFold(string(page.Title[0]), " ") {
				char = "_"
			} else if strings.EqualFold(string(page.Title[0]), "/") {
				char = "&slash;"
			} else { char = string(page.Title[0]) } */
			if !strings.EqualFold(page.Revisions[0].Text, "") {
				// if bufferCounter == 1 { buffer[0] = buffer[1]; bufferCounter = 0 }
				// bufferCounter++
				article.Offset = buffer[0]
				buffer[1] = wikiParser.GetOffset() + 1
				article.Title = page.Title
			}
			myMap[file][char] = append(myMap[file][char], article)
			buffer[0] = buffer[1]
			if counter == 0 {
				counter = 100
				runtime.ReadMemStats(&memstats)
				if memstats.Alloc / 1000000 > 2000 { // If InUseBytes is greater than 1500 MB
					fmt.Println(memstats.Alloc / 1000000)
					fmt.Println("GC'ing")
					go func() { WriteIndex(myMap, indexDir); /* wg.Done() */ }()
					for key, _ := range myMap[file] {
						delete(myMap[file], key)
					}
					myMap[file] = make(map[string][]ArticleElements)
				}
			}
			counter--
		}
		fmt.Println("processed", file)
		WriteIndex(myMap, indexDir)
		for key, _ := range myMap[file] {
			delete(myMap[file], key)
		}
	}
}
