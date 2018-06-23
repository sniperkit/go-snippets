package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math"
	"os"
	"strconv"
)

var total uint64
var gDb *sql.DB
var maxConn int = 5

func main() {
	db, err := sql.Open(`mysql`, `root:123456@tcp(127.0.0.1)/test?charset=utf8`)
	if nil != err {
		log.Fatalln(err)
	}
	db.SetMaxOpenConns(maxConn)
	defer db.Close()
	defer gDb.Close()
	gDb = db
	db.QueryRow(`SELECT max(id) FROM test`).Scan(&total)
	log.Println(total)
	process()
}

var limit uint16 = 100
var offset uint64 = 0
var countChan chan int
var f *os.File

func process() {
	tmpF, err := os.OpenFile(`D:/test.log`, os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if nil != err {
		log.Fatalln(err)
	}
	f = tmpF
	defer tmpF.Close()
	defer f.Close()
	size := uint32(math.Ceil(float64(total / uint64(limit))))
	count := 0
	for size > 0 {
		if size > uint32(maxConn) {
			count = maxConn
			size = size - uint32(maxConn)
		} else {
			count = int(size)
			size = 0
		}
		countChan = make(chan int, count)
		for i := 0; i < count; i++ {
			go pageData(countChan)
		}
		for i := 0; i < count; i++ {
			<-countChan
		}
	}
}

func pageData(c chan<- int) {
	var id uint64
	var name string
	offset = offset + uint64(limit)
	rows, err := gDb.Query("SELECT id, name FROM test LIMIT ? OFFSET ?", limit, offset)
	if nil != err {
		log.Fatalln(err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&id, &name); nil != err {
			log.Fatalln(err)
		}
		log.Println(id)
		f.WriteString(strconv.Itoa(int(id)) + "," + name + "\n")
	}
	c <- 1
}
