package main

import "fmt"

func main() {
	user := "Maik"
	group := info(user)
	fmt.Println(user, " is in group ", group)
}

func info(user string) (group string) {
	if user == "Maik" {
		group = "root"
	} else {
		group = "user"
	}

	return
}
