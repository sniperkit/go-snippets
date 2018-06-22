package main

import (
	"database/sql"
	_"github.com/go-sql-driver/mysql"

	"fmt"
	"log"
)

// TODO 批量插入
// TODO 这种连接方式如何实现同时进行不同协程操作

func main() {
	query2()
}

func query() {
	db, err := sql.Open("mysql", "root:1406Wr641231,.@tcp(192.168.201.133:3307)/mydb")
	check(err)

	rows, err := db.Query("SELECT * FROM mydb.announcement")
	check(err)

	for rows.Next() {
		columns, _ := rows.Columns()

		scanArgs := make([]interface{}, len(columns))
		values := make([]interface{}, len(columns))

		for i := range values {
			scanArgs[i] = &values[i]
		}

		// 将每行数据保存为一个map数据结构
		// 因为map是无序的，所以每次获取到的结果也是无序的
		err = rows.Scan(scanArgs...)
		record := make(map[string]string)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		fmt.Println(record)
	}
	rows.Close()
}

func query2() {
	db, err := sql.Open("mysql","root:1406Wr641231,.@tcp(" +
		"192.168.201.133:3307)/mydb?charset=utf8")
	check(err)

	rows, err := db.Query("SELECT id,state,imgUrl,createDate FROM announcement")
	check(err)

	defer rows.Close()

	for rows.Next() {
		var id int
		var state int
		var imgUrl string
		var createDate string

		// 这里需要注意 Scan 中定义变量（的指针位置）的顺序必须和Query中的顺序一一对应
		if err := rows.Scan(&id,&state,&imgUrl,&createDate); err != nil{
			log.Fatal(err)
		}
		fmt.Println(imgUrl, id, createDate, state)
		// 若以其他形式将查询到的数据输出，则需要注意各个数据的类型
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

}


//插入数据
func insert()  {
	db, err := sql.Open("mysql", "root:@/shopvisit")
	check(err)

	stmt, err := db.Prepare(`INSERT announcement (imgUrl, detailUrl, createDate, state) VALUES (?, ?, ?, ?)`)
	check(err)

	res, err := stmt.Exec("/visitshop/img/ann/cofox1.png",nil,"2017-09-06",0)
	check(err)

	id, err := res.LastInsertId()
	check(err)

	fmt.Println(id)
	stmt.Close()

}

//修改数据
func update() {
	db, err := sql.Open("mysql", "root:@/shopvisit")
	check(err)

	stmt, err := db.Prepare("UPDATE announcement set imgUrl=?, detailUrl=?, createDate=?, state=? WHERE id=?")
	check(err)

	res, err := stmt.Exec("/visitshop/img/ann/cofox2.png", nil, "2017-09-05", 1, 7)
	check(err)

	num, err := res.RowsAffected()
	check(err)

	fmt.Println(num)
	stmt.Close()
}

//删除数据
func remove() {
	db, err := sql.Open("mysql", "root:@/shopvisit")
	check(err)

	stmt, err := db.Prepare("DELETE FROM announcement WHERE id=?")
	check(err)

	res, err := stmt.Exec(7)
	check(err)

	num, err := res.RowsAffected()
	check(err)

	fmt.Println(num)
	stmt.Close()

}


func check(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
