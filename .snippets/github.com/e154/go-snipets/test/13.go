package main

// 13. При следующем объявлении типов какие строки (statements) корректны

// Ответы:
// s = string(r), s = string(b)

var (
	s string
	b []byte
	r []rune
	i int
	f float32
)

func main() {

	//s = b
	//s = r
	s = string(r)
	s = string(b)
	//i = f
}

