package main

/*
<table>
<tr> <th>head 1</th> <th>head 2</th> </tr>
<tr> <td>cell 1</td> <td>cell 2</td> </tr>
<tr> <td>cell 3</td> <td>cell 4</td> </tr>
</table>
*/

import (
	"fmt"
	"net/http"
	"text/template"
)

// TODO create and parse templates at compile time

const tableA = "<table>\n"
const tableO = "</table>\n"
const lineA = "<tr>"
const lineO = "</tr>\n"
const templH = "<th>{{.}}</th>"
const templC = "<td>{{.}}</td>"

// saving key strokes

func tableHead(w http.ResponseWriter, s string) {
	tableOut(w, s, templH)
}

func tableCell(w http.ResponseWriter, s string) {
	tableOut(w, s, templC)
}

func tableOut(w http.ResponseWriter, s string, templ string) {
	t, err := template.New("th").Parse(templ)
	checkY(err)
	err = t.Execute(w, s)
	checkY(err)
}

func tableDuo(w http.ResponseWriter, s1 string, s2 string) {
	fmt.Fprint(w, lineA)
	tableCell(w, s1)
	tableCell(w, s2)
	fmt.Fprint(w, lineO)
}
