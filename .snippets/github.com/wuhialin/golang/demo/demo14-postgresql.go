/*
* @Author: wuhailin
* @Date:   2017-12-04 15:24:52
* @Last Modified by:   wuhailin
* @Last Modified time: 2017-12-05 08:39:11
 */
package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"math/rand"
	"runtime"
	"strings"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.Gosched()
	db, err := sql.Open("postgres", `postgres://postgres:123456@localhost/template1?sslmode=disable`)
	defer db.Close()
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	var i uint
	var randString, sqlString string
	var md5Byte [16]byte
	var j uint64 = 0
	dbExecChan := make(chan bool, 10)
	for {
		dbExecChan <- true
		go func() {
			var md5s []string
			var params []interface{}
			m, n := 1, 2
			for i = 0; i < 5000; i++ {
				randString = fmt.Sprintf("%d%d", time.Now().UnixNano(), rand.Int63())
				md5Byte = md5.Sum([]byte(randString))
				md5s = append(md5s, fmt.Sprintf("($%d, $%d)", m, n))
				m += 2
				n += 2
				params = append(params, fmt.Sprintf("%x", md5Byte), time.Now().Unix())
			}
			sqlString = `INSERT INTO test_md5 (md5, create_time) VALUES ` + strings.Join(md5s, ", ")
			_, err = db.Exec(sqlString, params...)
			if err != nil {
				fmt.Println(err)
			}
			<-dbExecChan
		}()
		j++
		println(j)
	}
}
