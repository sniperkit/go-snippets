package main
import (
	// "github.com/mvryan/fasttag"
	// "github.com/Obaied/RAKE.Go"
	"github.com/pressure679/go-wikiparse"
	// "github.com/dustin/go-wikiparse"
	"github.com/neurosnap/sentences"
	// "github.com/neurosnap/sentences/data"
	// "github.com/dataence/porter2"
	"os"
	"io"
	"io/ioutil"
	"bufio"
	"strings"
	"strconv"
	"runtime"
	"fmt"
	"regexp"
	"math"
	"sort"
	"sync"
)
type Profile struct {
	Epoch uint8
	MemStats runtime.MemStats
}
type Word string
type Sentence struct { Start, End []byte }
type WordElements struct {
	Sentences []Sentence
}
type WordMap map[*Word]WordElements
type Section string
type SectionElements map[*Section]WordMap
type Article string
type ArticleElements struct {
	// Title string
	// Offset int64

	Text string
	Words map[*HashWord][]Sentence
	
	Links []*Article
	Parents [7][]*Article

	ZValue float64
}
func (ae *ArticleElements) AppendLinkIfNotExists(article *Article) (offset int) {
	for linkNum, _ := range ae.Links {
		if ae.Links[linkNum] == article { return linkNum }
	}
	ae.Links = append(ae.Links, article)
	return len(ae.Links) -1
}
func (ae ArticleElements) AppendWeightedParent(weight int, article *Article) {
	ae.Parents[weight] = append(ae.Parents[weight], article)
}
type HashWord struct {
	Word Word
	Parents []*Article
	ZScore float64
}
type HashTable struct {
	Articles []Article
	Words []HashWord
}
func (HashTable *HashTable) SearchArticle(article Article) int {
	for num, _ := range HashTable.Articles {
		if HashTable.Articles[num] == article { return num }
	}
	return -1
}
func (HashTable *HashTable) SearchWord(word Word) int {
	for num, _ := range HashTable.Words {
		if HashTable.Words[num].Word == word { return num }
	}
	return -1
}
func (HashTable *HashTable) AppendArticleIfNotExists(article Article) int {
	offset := HashTable.SearchArticle(article)
	if offset == -1 { HashTable.Articles = append(HashTable.Articles, article); return len(HashTable.Articles) } else { return offset }
	return -1
}
func (HashTable HashTable) AppendWordIfNotExists(word Word) int {
	offset := HashTable.SearchWord(word)
	if offset == -1 { HashTable.Words = append(HashTable.Words, HashWord{Word: word}); return len(HashTable.Words) - 1 }
	return -1
}
type LinearRegressionElements struct {
	// ParentedMin, ParentedMax map[*Article]uint
	// TODO: ParentedMin, ParentedMax uint map[*Article]uint // This is going to have an individual loop method after all articles have been read and vertice connected. - their values well for linear regression-pagerank leak security of z-scores.
	SumArticles, SumLinks uint
	ArticleMinLinked, ArticleMaxLinked, ArticleSumLinked map[*Article]uint
	
	SumWords, SumSentences uint
	WordMinSentences, WordMaxSentences, WordSumSentences map[*Word]uint
}
type Graph struct {
	Profile Profile
	wg sync.WaitGroup
	osFileWG sync.WaitGroup
	Hash HashTable
	LrE LinearRegressionElements

	RegEx []*regexp.Regexp
	NeuroSnapSentencer *sentences.DefaultSentenceTokenizer
	XMLUnescapers [][]string
	
	WikiFiles []string
	IndexFiles []string
	BaseDir string

	Map map[*Article]ArticleElements
	
}
type wordSorter struct {
	words []*HashWord
	by func(w1, w2 *HashWord) bool
}
type ByWords func(w1, w2 *HashWord) bool
func (by ByWords) Sort(words []*HashWord) {
	ws := &wordSorter{
		words: words,
		by: by,
	}
	sort.Sort(ws)
}
func (ws *wordSorter) Len() int {
	return len(ws.words)
}
func (ws *wordSorter) Swap(i, j int) {
	ws.words[i], ws.words[j] = ws.words[j], ws.words[i]
}
func (ws *wordSorter) Less(i, j int) bool {
	return ws.by(ws.words[i], ws.words[j])
}
type articleSorter struct {
	articles []*ArticleElements
	by func(a1, a2 *ArticleElements) bool
}
type ByArticles func(a1, a2 *ArticleElements) bool
func (by ByArticles) Sort(articles []*ArticleElements) {
	as := &articleSorter{
		articles: articles,
		by: by,
	}
	sort.Sort(as)
}
func (s *articleSorter) Len() int {
	return len(s.articles)
}
func (s *articleSorter) Swap(i, j int) {
	s.articles[i], s.articles[j] = s.articles[j], s.articles[i]
}
func (s *articleSorter) Less(i, j int) bool {
	return s.by(s.articles[i], s.articles[j])
}
const (
	Chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)
var counter uint = 0
var Dirs []string
var Files [][]string
func main() {
	var osFile *os.File
	var MyGraph Graph
	baseArticle := "Computer Science"
	
	runtime.ReadMemStats(&MyGraph.Profile.MemStats)

	// 1st) nowiki text, 2) comment (no wiki text), 3) wikipedia links, 4) section headers.
	MyGraph.RegEx = append(MyGraph.RegEx, regexp.MustCompile(`(?ms)<nowiki>.*</nowiki>`))
	MyGraph.RegEx = append(MyGraph.RegEx, regexp.MustCompile(`(?ms)<!--.*-->`))
	MyGraph.RegEx = append(MyGraph.RegEx, regexp.MustCompile(`\[\[([^\|\]]+)`)) // \[\[
	MyGraph.RegEx = append(MyGraph.RegEx, regexp.MustCompile("[=]{2,}(.+)[=]{2,}"))
	// MyGraph.RegEx = append(MyGraph.RegEx, regexp.MustCompile(`<ref>(.+)</ref>`)
	MyGraph.XMLUnescapers = [][]string{{"&lt", "<"}, {"&gt", ">"}, {"&quot", "\""}, {"&amp", "&"}}

	var err error
	Dirs, err = GetSubDirs("/home/naamik/Documents/wikipedia/index/")
	Files = make([][]string, len(Dirs))
	if err != nil { panic(err) }
	for dirNum, _ := range Dirs {
		Files[dirNum], err = GetFilesFromDir("/home/naamik/Documents/wikipedia/index/" + Dirs[dirNum])
	}
	// fmt.Println(Dirs)

	/* if err != nil { panic(err) }
	trainingData, err := sentences.LoadTraining(dataAsset)
	if err != nil { panic(err) } */
	// MyGraph.Sentence = sentences.NewSentenceTokenizer(trainingData)
	
	MyGraph.Map = make(map[*Article]ArticleElements)
	MyGraph.Hash = HashTable{}
	MyGraph.Hash.Articles = []Article{}
	offset := MyGraph.Hash.AppendArticleIfNotExists(Article(baseArticle))
	MyGraph.Map[&MyGraph.Hash.Articles[offset-1]] = ArticleElements{}
	MyGraph.Profile.Epoch = 0
	
	MyGraph.ReadArticle(0, baseArticle, 0, osFile)
	// ReadArticle(0, baseArticle, 0, &MyGraph, osFile)
	MyGraph.RankArticles()
	// as := MyGraph.OrderByPagerank
	// fmt.Println(as)
	/* osFile, err := os.Create("/home/naamik/article/text-2.org")
	if err != nil { panic(err) }
	fmt.Fprintln(osFile, MyGraph)
	err = osFile.Close()
	if err != nil { panic(err) } */
}
func (ae ArticleElements) SetText(page string) {
	ae.Text = page
}
var started, fileOpen, fileOpen2, started2 bool = false, false, false, false
func GetIndex(article string, osFile *os.File, Graph *Graph) (file string, offset int, err error) {
	// if fileOpen && !fileOpen2 { Graph.osFileWG.Wait() }
	// fileOpen = true
	// fileOpen2 = true
	Graph.osFileWG.Add(1)
	// dirs, err := GetSubDirs("/home/naamik/Documents/wikipedia/index/")
	// if err != nil { return "", -1, err }
	for dirNum, _ := range Dirs {
		if len(article) > 0 {
			// fmt.Println(Dirs[dirNum])
			if !strings.EqualFold(string(Dirs[dirNum][0]), string(article[0])) { continue }
			// files, err := GetFilesFromDir("/home/naamik/Documents/wikipedia/index/" + Dirs[dirNum])
			// if err != nil { return "", -1, err }
			for fileNum, _ := range Files[dirNum] {
				// for fileNum2, _ := range Files[fileNum] {
				osFile, err = os.Open("/home/naamik/Documents/wikipedia/index/" + Dirs[dirNum] + "/" + Files[dirNum][fileNum])
				if err != nil { osFile.Close(); return "", -1, err }
				bufioReader := bufio.NewReader(osFile)
				for {
					line, _, err := bufioReader.ReadLine()
					if err != nil {
						if err == io.EOF { break }
						return "", -1, err }
					parts := strings.Split(string(line), "&:&")
					if !strings.EqualFold(parts[0], article) { continue }
					file = string(Files[dirNum][fileNum][:len(Files[dirNum][fileNum])-8])
					offset, err = strconv.Atoi(string(parts[1]))
					if err != nil { osFile.Close(); return "", -1, err }
				}
				err = osFile.Close()
				if err != nil { return "", -1, err }
				// }
			}
		}
	}
	Graph.osFileWG.Done()
	// fileOpen = false
	return file, offset, nil
}
func GetFilesFromDir(directory string) (files []string, err error) {
	osFileInfo, err := ioutil.ReadDir(directory); if err != nil { return nil, err }
	for _, fileInfo := range osFileInfo { if !fileInfo.IsDir() { files = append(files, fileInfo.Name()) } }; return files, err
}
func GetSubDirs(directory string) (dirs []string, err error) {
	// fmt.Println("Getting subdirectories of index")
	osFileInfo, err := ioutil.ReadDir(directory)
	if err != nil { return nil, err }
	for _, fileInfo := range osFileInfo {
		if fileInfo.IsDir() {
			dirs = append(dirs, fileInfo.Name())
		}
	}
	return dirs, err
}
// func (Graph *Graph) ReadArticle(depth uint8, BaseArticle string, rootOffset int, osFile *os.File) {
func (Graph *Graph) ReadArticle(depth uint8, BaseArticle string, rootOffset int, osFile *os.File) {
	defer recover()
	var text string
	file, index, err := GetIndex(string(BaseArticle), osFile, Graph)
	if err != nil { panic(err) }
	if !strings.EqualFold(file, "") {
		osFile, err = os.Open("/home/naamik/Documents/wikipedia/xml_articles/" + file)
		if err != nil { panic(err) }
		_, err = osFile.Seek(int64(index), 0)
		wikiParser, err := wikiparse.NewParser(osFile)
		if err != nil { panic(err) }
		page, err := wikiParser.Next()
		if err != nil { panic(err) }
		err = osFile.Close()
		if err != nil { panic(err) }
		if !strings.EqualFold(page.Revisions[0].Text, "") {
			text = Graph.RegEx[0].ReplaceAllString(Graph.RegEx[1].ReplaceAllString(page.Revisions[0].Text, ""), "")
		}
		for num, _ := range Graph.XMLUnescapers {
			text = strings.Replace(text, Graph.XMLUnescapers[num][0], Graph.XMLUnescapers[num][1], -1)
		}
	}
	links := Graph.RegEx[2].FindAllString(text, -1)
	var absLink string
	fmt.Println(BaseArticle)
	Graph.Map[&Graph.Hash.Articles[rootOffset]].SetText(text)
	for linkNum, _ := range links {
		if strings.Contains(links[linkNum][2:], "File:") { continue }
		if strings.Contains(links[linkNum][2:], "Image:") { continue }
		if strings.Contains(links[linkNum][2:], "User:") { continue }
		if strings.Contains(links[linkNum][2:], "User talk:") { continue }
		if strings.Contains(links[linkNum][2:], "User_talk:") { continue }
		var offset int
		if strings.Contains(string(links[linkNum][2:]), "#") {
			parts := strings.Split(links[linkNum][2:], "#")
			absLink = parts[0]
		} else {
			absLink = links[linkNum][2:]
		}
		offset = Graph.Hash.SearchArticle(Article(absLink))
		if offset == -1 {
			Graph.Hash.Articles = append(Graph.Hash.Articles, Article(absLink))
			offset = len(Graph.Hash.Articles) -1
			Graph.Profile.Epoch++
			if Graph.Profile.Epoch > 100 {
				runtime.ReadMemStats(&Graph.Profile.MemStats)
				Graph.Profile.Epoch = 0
				if Graph.Profile.MemStats.Alloc / 1000000 > 2600 { break }
			}
			if depth < 6 {
				ReadArticle(depth + 1, absLink, offset, Graph, osFile)
			}
		} else {
			Graph.Map[&Graph.Hash.Articles[offset]].AppendWeightedParent(int(depth), &Graph.Hash.Articles[rootOffset])
		}
	}
}
func ReadArticle(depth uint8, BaseArticle string, rootOffset int, Graph *Graph, osFile *os.File) {
	// if started && !started2 { Graph.wg.Wait() }
	// started2 = true
	// started = true
	// Graph.wg.Add(1)
	fmt.Println(BaseArticle, depth)
	defer recover()
	var text string
	Graph.osFileWG.Wait()
	file, index, err := GetIndex(string(BaseArticle), osFile, Graph)
	if err != nil { panic(err) }
	if !strings.EqualFold(file, "") {
		osFile, err = os.Open("/home/naamik/Documents/wikipedia/xml_articles/" + file)
		if err != nil { panic(err) }
		_, err = osFile.Seek(int64(index), 0)
		wikiParser, err := wikiparse.NewParser(osFile)
		if err != nil { panic(err) }
		page, err := wikiParser.Next()
		if err != nil { panic(err) }
		err = osFile.Close()
		if err != nil { panic(err) }
		if !strings.EqualFold(page.Revisions[0].Text, "") {
			text = Graph.RegEx[0].ReplaceAllString(Graph.RegEx[1].ReplaceAllString(page.Revisions[0].Text, ""), "")
		}
		for num, _ := range Graph.XMLUnescapers {
			text = strings.Replace(text, Graph.XMLUnescapers[num][0], Graph.XMLUnescapers[num][1], -1)
		}
	}
	links := Graph.RegEx[2].FindAllString(text, -1)
	var absLink string
	Graph.Map[&Graph.Hash.Articles[rootOffset]].SetText(text)
	for linkNum, _ := range links {
		if strings.Contains(links[linkNum][2:], "File:") { continue }
		if strings.Contains(links[linkNum][2:], "Image:") { continue }
		if strings.Contains(links[linkNum][2:], "User:") { continue }
		if strings.Contains(links[linkNum][2:], "User talk:") { continue }
		if strings.Contains(links[linkNum][2:], "User_talk:") { continue }
		var offset int
		if strings.Contains(string(links[linkNum][2:]), "#") {
			parts := strings.Split(links[linkNum][2:], "#")
			absLink = parts[0]
		} else {
			absLink = links[linkNum][2:]
		}
		offset = Graph.Hash.SearchArticle(Article(absLink))
		if offset == -1 {
			// if strings.EqualFold(absLink, "Computer Science") { fmt.Println(BaseArticle, "-", absLink) }
			Graph.Hash.Articles = append(Graph.Hash.Articles, Article(absLink))
			offset = len(Graph.Hash.Articles) -1
			Graph.Profile.Epoch++
			if Graph.Profile.Epoch > 100 {
				runtime.ReadMemStats(&Graph.Profile.MemStats)
				Graph.Profile.Epoch = 0
				if Graph.Profile.MemStats.Alloc / 1000000 > 2600 { break }
			}
			if depth < 6 {
				go ReadArticle(depth + 1, absLink, offset, Graph, osFile)
			}
		} else {
			Graph.Map[&Graph.Hash.Articles[offset]].AppendWeightedParent(int(depth), &Graph.Hash.Articles[rootOffset])
		}
	}
	// started = false
	// Graph.wg.Done()
}
func (Graph Graph) RankArticles() {
	// This loop will highly likely give high mean values. Consider modifying the linear regression algorithm (below).
	var articleMeanLinked map[*Article]float64
	for article, _ := range Graph.LrE.ArticleSumLinked {
		// articleMeanLinked[article] = articles.LrE.ArticleSumLinked[article] / articles.LrE.SumArticles
		articleMeanLinked[article] = float64(Graph.LrE.ArticleSumLinked[article]) / float64(len(Graph.Map))
	}
	var articleStdDev map[*Article]float64
	for article, _ := range articleMeanLinked {
		articleStdDev[article] = math.Pow(float64(Graph.LrE.ArticleSumLinked[article]) - articleMeanLinked[article], 2)
	}
	for article, _ := range articleStdDev {
		articleStdDev[article] = math.Sqrt(float64(Graph.LrE.ArticleSumLinked[article]) * articleMeanLinked[article])
	}
	for article, _ := range articleStdDev {
		Graph.Map[article].SetArticleZ(articleStdDev[article] / float64(Graph.LrE.ArticleMaxLinked[article]-Graph.LrE.ArticleMinLinked[article]))
	}
}
func (ae ArticleElements) SetArticleZ(value float64) {
	ae.ZValue = value
}
func (Graph Graph) RankWords() {
	var wordStdDev, wordMeanSentences map[*Word]float64
	for word, _ := range Graph.LrE.WordSumSentences {
		wordMeanSentences[word] = float64(Graph.LrE.WordSumSentences[word]) / float64(Graph.LrE.SumWords) / float64(len(Graph.Map))
	}
	for word, _ := range wordMeanSentences {
		wordStdDev[word] = math.Sqrt(float64(Graph.LrE.WordSumSentences[word]) * wordMeanSentences[word])
	}
	for wordNum, _ := range Graph.Hash.Words {
		Graph.Hash.Words[wordNum].ZScore = wordStdDev[&Graph.Hash.Words[wordNum].Word] / float64(Graph.LrE.WordMaxSentences[&Graph.Hash.Words[wordNum].Word]-Graph.LrE.WordMinSentences[&Graph.Hash.Words[wordNum].Word])
	}
}
func (Graph Graph) OrderByPagerank() (as articleSorter) {
	for _, ae := range Graph.Map {
		as.articles = append(as.articles, &ae)
	}
	z := func(article1, article2 *ArticleElements) bool {
		return article1.ZValue < article2.ZValue
	}
	ByArticles(z).Sort(as.articles)
	return as
}
func (Graph Graph) OrderByWord() (ws wordSorter) {
	for article, _ := range Graph.Map {
		for word, _ := range Graph.Map[article].Words {
			// for sentenceNum, _ := range Graph.Map[article].Words[word].Sentences {
			// ws.words = append(ws.words, Graph.Map[article].Words[word])
			ws.words = append(ws.words, word)

			// }
		}
	}
	z := func(word1, word2 *HashWord) bool {
		return word1.ZScore < word2.ZScore
	}
	ByWords(z).Sort(ws.words)
	return ws
}
