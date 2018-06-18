package main

func main() {
	a := App{}
	a.Initialize("172.17.0.4", "movies_db")
	a.Run(":5000")
}
