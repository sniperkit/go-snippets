package main

//import "github.com/mattn/go-sqlite3"
//import "github.com/mxk/go-sqlite/sqlite3"
//..or see https://github.com/mxk/go-sqlite
import (
 "database/sql"
 _ "github.com/go-sql-driver/mysql"
 "fmt"
 "log"
)
func main() {
//    conn,err := Open('ex1')
	//fmt.Println(sql.Drivers())
	db, err := sql.Open(sql.Drivers()[0], "root:admin@/test")
	//fmt.Println(db)
	if(err!=nil){
	   fmt.Println(err)
	}
	//res:=db.Ping()
	//fmt.Println(res)
	rows,err:=db.Query("select * from user")
	for(rows.Next()){
		var id int
		var userid string 
		if err := rows.Scan(&id,&userid); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v is %v\n", id, userid)
	}
	
	
	
	/*columns,err2:=rows.Columns()
	if(err2!=nil){
		fmt.Print("error selecting columns:\n")
		fmt.Println(err2)
	}
	
	//column[i] throws an error below
	for i,column := range columns{
		if(column[i]!=nil){
			fmt.Println("not nil")//column[i])
		}
	}*/
}