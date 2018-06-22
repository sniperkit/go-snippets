package main

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"log"
	"net/http"
	"os"
	"strings"
//	"time"
	"bufio"
	"strconv"
	"io"
	"sync"
	"path/filepath"
	"time"
)

type LinguaResp struct {
	ErrorMsg string      `json:"error_msg"`
	Count    uint      	`json:"count_words"`
	ShowMore bool      	`json:"show_more"`
	Userdict []Userdict  `json:"userdict3"`
}

type Userdict struct {
	Name  string `json:"name"`
	Count uint    `json:"count"`
	Words []Word `json:"words"`
}

type Word struct {
	Id          uint            `json:"word_id"`
	Value       string          `json:"word_value"`
	Transcript  string          `json:"transcription"`
	//	Created     time.Time       `json:"created_at"`
	//	LastUpdated time.Time       `json:"last_updated_at"`
	Translates   []UserTranslate `json:"user_translates"`
	SoundUrl     string          `json:"sound_url"`
	PictureUrl   string          `json:"picture_url"`
}

type UserTranslate struct {
	Value string `json:"translate_value"`
}

const (
	linguaDictUrl  = "http://lingualeo.com/userdict/json"
	linguaLoginUrl = "http://api.lingualeo.com/api/login"
	linguaTranslateUrl = "http://api.lingualeo.com/gettranslates"
	pageCount      = 1 // restriction policy, you have 116
	httpTimeout    = 15
)

var client *gorequest.SuperAgent
var showMore = true
var importFile *os.File
var supplyChan chan string
var wg sync.WaitGroup

type Config struct {
	Email    string
	Password string
}

func readConfig(configFile string) Config {
	file, err := os.Open(configFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	config := Config{}
	configScanner := bufio.NewScanner(file)
	if configScanner.Scan() {
		config.Email = configScanner.Text()
	}
	if configScanner.Scan() {
		config.Password = configScanner.Text()
	}
	//	fmt.Printf("%v \n", config)
	return config
}

func main() {
	config := readConfig("settings.txt")
	authLeo(config.Email, config.Password)
	fmt.Printf("I'll grab %v pages", pageCount)
	supplyChan = make(chan string, 2)
	importFileName := "import.anki." + fmt.Sprintf("%v", time.Now().Local()) + ".txt"
	importFile, _ = os.OpenFile(importFileName, os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666); defer importFile.Close()
	go downloadResources()
	
	for i := 1; showMore; i++ {
		leoAskPage(i)
	}
	
	wg.Wait()
	fmt.Println("Finished")
}

/**
 * json requested by pages, @see showMore flag
 */
func leoAskPage(page int) {
	pageNumber := strconv.Itoa(page)
	fmt.Println("Url " + linguaDictUrl + "sortBy=date&wordType=0&filter=all&page=" + pageNumber)
	_, body, errs := client.Post(linguaDictUrl).
									Query("page=" + pageNumber).
									Query("sortBy=date").Query("wordType=0").Query("filter=all").
									End()
	if errs != nil {
		log.Fatalf("Error %v \n", errs)
		os.Exit(1)
	}
//	fmt.Printf("\nPage %v Body %v\n", pageNumber, body)
	var linguaResp LinguaResp
	json.NewDecoder(strings.NewReader(body)).Decode(&linguaResp)
	showMore = linguaResp.ShowMore
//	fmt.Printf("\n ShowMore %v \n", showMore)
	userdicts := linguaResp.Userdict
	fmt.Printf("\n === %v User Dictionaries \n", len(userdicts))
	for i := 0; i < len(userdicts); i++ {
		ankiImport(userdicts[i].Words, importFile, supplyChan)
//		userdicts[i].Print()
	}
	//	fmt.Printf("Decoded %v \n", linguaResp)
}

func (d *Userdict) Print() {
	fmt.Printf("\n === Dictionary '%v' [ %v words] \n", d.Name, d.Count)
	words := d.Words
	for i := 0; i < len(words); i++ {
		words[i].Print()
	}
}

func (w *Word) Print() {
	fmt.Printf("= %v [ %v ] \n", w.Value, w.Transcript)
	//	fmt.Printf("  picture %v, sound %v \n", w.PictureUrl, w.SoundUrl)
	wordTranslates := w.Translates
	for i := 0; i < len(wordTranslates); i++ {
		fmt.Printf("  - %v \n", wordTranslates[i])
	}
}

func authLeo(email, password string) {
	client = gorequest.New() //.Timeout(httpTimeout*time.Second)
	resp, body, errs := client.Get(linguaLoginUrl).Query("email=" + email).Query("password=" + password).End()
	if errs != nil {
		log.Fatalf("%v \n", errs)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed login %v", resp.Status)
	}
	var loginResp LoginResp
	err := json.NewDecoder(strings.NewReader(body)).Decode(&loginResp)
	if err != nil {
		log.Fatalf("Failed decode %v", body)
	}
	//	log.Printf("login parse %v \n", loginResp)
}

func linguaLeoAPI() string {
	result := fmt.Sprintf(linguaDictUrl)
	log.Println("Api " + result)
	return result
}

func ankiImport(words []Word, file *os.File, sinkChan chan string) error {
	w := bufio.NewWriter(file)
	for i := 0; i < len(words); i++ {
		word := words[i]
		str := word.plainImport()
		wg.Add(2)
		sinkChan <- word.SoundUrl
		sinkChan <- word.PictureUrl
		fmt.Printf("ankiImport %v", str)
		_, err := w.WriteString(str)
		if err != nil {
			panic(err)
		}
//		fmt.Printf("Written %v", written)
	}
	w.Flush()
	return nil
}

func (w *Word) plainImport() string {
	localSound := getName(w.SoundUrl)
	localPiclure := getName(w.PictureUrl)
	return fmt.Sprintf("%v\t%v\t%v\t%v\t\n", w.Value, w.Translates, localPiclure, localSound)
}

func downloadResources() error {
	for {
		r := <- supplyChan
		go downloadFromUrl(r)		
	}
	
}

/**
 * path: FS path or Url
 */
func getName(path string) string {
	tokens := strings.Split(path, "/")
	fileName := tokens[len(tokens)-1]
	return fileName 
}

func downloadFromUrl(url string) {
	defer wg.Done()
	//todo remove stub
//	return
	
	if url == "" {
		return
	}
	if !strings.HasPrefix(url, "http:") {
		 url = "http:" + url
	}
	fileName := getName(url)
//	fmt.Println("Downloading", url, "to", fileName)

	absPath, _ := filepath.Abs("./resources/" + fileName)
	if _, err := os.Stat(absPath); err == nil {
		return
	}
	output, err := os.Create(absPath)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
//	fmt.Println(n, "bytes downloaded.")
}

type LoginResp struct {
	ErrorMsg   string `json:"error_msg"`
	User       User   `json:"user"`
}

type User struct {
	Username     string `json:"nickname"`
	Id           int    `json:"user_id"`
	AutologinKey string `json:"autologin_key"`
}

func addWord(word, tword string) {
	url := "http://api.lingualeo.com/addword"
	params := fmt.Sprintf("?word=%v&tword=%v", word, tword)
	fmt.Printf("%v %v", url, params)
}

func getTranslates(word string) {
	client.Get(linguaTranslateUrl).Query("word=" + word);
}

