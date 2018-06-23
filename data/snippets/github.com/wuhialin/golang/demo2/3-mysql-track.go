package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	source := "pms_select:gVjw6Ysr@tcp(192.168.0.144)/pms?charset=utf8"
	db, err := sql.Open("mysql", source)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(1000)
	wait := new(sync.WaitGroup)
	var payId uint
	f, err := os.OpenFile("track.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
	for i := 0; i < 10; i++ {
		wait.Add(1)
		if err != nil {
			panic(err)
		}
		go func(i int) {
			start := time.Now()
			var payIds []string
			rows, err := db.Query("SELECT pay_id FROM pms_finance_order WHERE pay_id > 2425150 OR update_time > '2017-12-25 11:27:52' ORDER BY NULL LIMIT 100")
			if err != nil {
				panic(err)
			}
			for rows.Next() {
				rows.Scan(&payId)
				payIds = append(payIds, fmt.Sprint(payId))
			}
			log.Println(i, time.Since(start))
			_, err = f.WriteString(fmt.Sprintln(i, time.Since(start), strings.Join(payIds, ",")))
			if err != nil {
				panic(err)
			}
			wait.Done()
		}(i)
	}
	wait.Wait()
}
