package main

import(
	"github.com/gocql/gocql"
)

func main() {
	cluster := gocql.NewCluster("192.168.1.1", "192.168.1.2", "192.168.1.3")

