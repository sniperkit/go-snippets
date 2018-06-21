package videos

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

const host = "http://www.80s.tw/movie/"

type Video struct {
	gorm.Model
	Url        string `gorm:"index"`
	ImageSrc   string
	Name       string `gorm:"index"`
	Alias      string
	Desc       string
	VideoSrc   string
	ThunderSrc string
}

var mux map[string]func(http.ResponseWriter, *http.Request)

func main() {
	dbInit()
	db, err := gorm.Open("sqlite3", "80s.db")
	checkErr(err)
	defer db.Close()
	server := http.Server{
		Addr:        ":8888",
		Handler:     &apiHandler{},
		ReadTimeout: 5 * time.Second,
	}

	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/api"] = api
	fmt.Println("http server at: 8888")
	err = server.ListenAndServe()
	checkErr(err)
}

type apiHandler struct{}

func (*apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}
	api(w, r)
	// w.WriteHeader(http.StatusNotFound)
	// io.WriteString(w, "something wrong!")
}

func api(w http.ResponseWriter, r *http.Request) {
	var keys []string
	var ok, find bool
	var rel, search string
	var page, id int
	var rows []Video
	var ret []byte
	// fmt.Println(r.URL.Path)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != "GET" && r.URL.Path != "/api" {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "Nothing")
		return
	}
	db, err := gorm.Open("sqlite3", "80s.db")
	checkErr(err)
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	keys, ok = r.URL.Query()["id"]
	if ok && len(keys) >= 1 {
		id, err = strconv.Atoi(keys[0])
		checkErr(err)
		db.Where(&Video{Url: rel}).Find(&rows, id)
		ret, err = json.Marshal(rows)
		checkErr(err)
		find = true
	}

	keys, ok = r.URL.Query()["rel"]
	if ok && len(keys) >= 1 {
		rel = keys[0]
		db.Where(&Video{Url: rel}).Find(&rows)
		ret, err = json.Marshal(rows)
		checkErr(err)
		find = true
	}

	keys, ok = r.URL.Query()["page"]
	if ok && len(keys) >= 1 {
		page, err = strconv.Atoi(keys[0])
		checkErr(err)
		page = page - 1
		db.Limit(10).Offset(page * 10).Find(&rows)
		ret, err = json.Marshal(rows)
		checkErr(err)
		find = true
	}

	keys, ok = r.URL.Query()["search"]
	if ok && len(keys) >= 1 {
		search = keys[0]
		search = "%" + search + "%"
		db.Limit(10).Where("name LIKE ? OR alias LIKE ?", search, search).Find(&rows)
		ret, err = json.Marshal(rows)
		checkErr(err)
		find = true
	}

	if !find {
		ret, _ = json.Marshal(map[string]string{"ret": "nothing"})
	}
	_, err = w.Write(ret)
	checkErr(err)
}

func dbInit() {
	db, err := gorm.Open("sqlite3", "80s.db")
	checkErr(err)
	defer db.Close()
	db.AutoMigrate(&Video{})
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
