// application framework
// Copyright (C) 2014  geosoft1@gmail.com
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
	"os"
)

//global variable to access template from modules
var t *template.Template

//standard framework page format
//	page title
//	status field, can keep loged user or other informations
//	body is generate from modules
type Page struct {
	Title, Status string
	Body          template.HTML
}

func init() {
	t = template.New("templ")
	t.ParseGlob("templates/*.html")
}

func main() {
	os.Remove("./foo.db")
	db, _ := sql.Open("sqlite3", "./foo.db")

	//some tables
	//	user table keep application users
	//	person table is only for demo modules
	db.Exec(`
	create table user (id integer not null primary key autoincrement, username text, password text);
	insert into user (username,password) values ("george","");
	insert into user (username,password) values ("john","");
	create table person (id integer not null primary key autoincrement, name text, age int, address text);
	insert into person (name,age,address) values ("George",38,"Sesame Street,Romania");
	insert into person (name,age,address) values ("Gill Bates",55,"Linux Street 10");
	insert into person (name,age,address) values ("Tinus Lorvalds",42,"Windows Bay 1");
	//delete from foo;
	`)
	db.Close()

	//set a log file for apllication
	//	use log.Println("message") in module to log actions
	f, err := os.OpenFile("logfile.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		println("cannot open logfile")
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("app started")

	// start web server
	http.ListenAndServe(":8080", nil)
}
