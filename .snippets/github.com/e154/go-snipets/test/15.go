package main

// 15. Какие из присваиваний корректны

// Ответы:
// i = Int(j)

type Int int64

var (
	i Int
	j int64
)

func main() {

	//i = j
	i = Int(j)
	//i = j.(Int)
}

