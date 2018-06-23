package main

// 26. Что делает данная конструкция?

// Ответы:
// объявляет канал только для чтения

func test26(ch <- chan int)  {

}

func main() {
	var c <- chan int
	c = make(<- chan int)


	go test26(c)
}

