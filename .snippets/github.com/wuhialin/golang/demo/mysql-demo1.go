/*
* @Author: wuhailin
* @Date:   2017-10-05 17:34:17
* @Last Modified by:   wuhailin
* @Last Modified time: 2017-10-06 10:30:45
 */

package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func main() {
	start := time.Nanosecond
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1)/ifs?charset=utf8")
	if nil != err {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT id FROM b_story")
	if nil != err {
		panic(err)
	}
	defer rows.Close()

	var count uint64
	count = 0
	for count < 4000000000 {
		count++
		//println(count)
	}
	time.Sleep(3 * time.Second)
	// for rows.Next() {
	// 	var id uint32
	// 	if err := rows.Scan(&id); nil != err {
	// 		panic(err)
	// 	}
	// 	count += 1
	// }
	end := time.Nanosecond
	println((end - start) / 1000000000)
	//println(count)
}
