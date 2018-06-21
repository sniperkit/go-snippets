//package videos

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
	. "videos"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const APIDB = "apiv2.db"

var mux map[string]func(http.ResponseWriter, *http.Request)

func main() {
	dbInit()
	db, err := gorm.Open("sqlite3", APIDB)
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
}

func api(w http.ResponseWriter, r *http.Request) {
	var keys []string
	var ok bool
	var search string
	var page, id, ty int
	var rows []VideoORM
	var ret []byte
	errInfo, _ := json.Marshal(map[string]string{"ret": "params erros"})
	// fmt.Println(r.URL.Path)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != "GET" && r.URL.Path != "/api" {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "Nothing")
		return
	}
	db, err := gorm.Open("sqlite3", APIDB)
	checkErr(err)
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	keys, ok = r.URL.Query()["id"]
	if ok && len(keys) >= 1 {
		id, err = strconv.Atoi(keys[0])
		if err != nil {
			w.Write(errInfo)
			return
		}
	}

	keys, ok = r.URL.Query()["page"]
	if ok && len(keys) >= 1 {
		page, err = strconv.Atoi(keys[0])
		if err != nil {
			w.Write(errInfo)
			return
		}
	}

	keys, ok = r.URL.Query()["q"]
	if ok && len(keys) >= 1 {
		search = keys[0]
		search = "%" + search + "%"
	}

	keys, ok = r.URL.Query()["tid"]
	if ok && len(keys) >= 1 {
		ty, err = strconv.Atoi(keys[0])
		if err != nil {
			w.Write(errInfo)
			return
		}
	}

	if page > 0 && ty > 0 {
		page = (page - 1) * 10
		db.Limit(20).Where("tid = ?", ty).Offset(page).Find(&rows)
		ret, err = json.Marshal(V2VXML(rows))
		checkErr(err)
	}

	if page > 0 {
		page = (page - 1) * 10
		db.Limit(8).Offset(page).Find(&rows)
		ret, err = json.Marshal(V2VXML(rows))
		checkErr(err)
	}

	if len(search) > 0 {
		db.Limit(10).Where("name LIKE ?", search, search).Find(&rows)
		ret, err = json.Marshal(V2VXML(rows))
		checkErr(err)
	}

	if id > 0 {
		db.Find(&rows, id)
		ret, err = json.Marshal(V2VXML(rows))
		checkErr(err)
	}

	_, err = w.Write(ret)
	checkErr(err)
}

func dbInit() {
	db, err := gorm.Open("sqlite3", APIDB)
	checkErr(err)
	defer db.Close()
	db.AutoMigrate(&VideoORM{})
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
