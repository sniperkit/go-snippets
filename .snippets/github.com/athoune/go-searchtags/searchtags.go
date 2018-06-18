package main

var searchQueue chan queryAnswer

func main() {
	searchQueue = make(chan queryAnswer, 200)
	for i := 0; i < 4; i++ {
		go StartSearch()
	}
	loadData()
	startHttp()
}
