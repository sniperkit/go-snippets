package main

import (
	"os"
	"text/template"
)

type User struct {
	FirstName string
	LastName  string
	Age       int
}

var tplUser string = `FirstName: {{.FirstName}}
LastName: {{.LastName}}
Age: {{ .Age }}`

func main() {

	tpl := template.New("UserTemplate")
	tpl, _ = tpl.Parse(tplUser)

	maik := User{
		FirstName: "Maik",
		LastName:  "Ellerbrock",
		Age:       35}

	tpl.Execute(os.Stdout, maik)
}
