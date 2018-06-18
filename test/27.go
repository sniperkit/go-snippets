package main

import (
	"fmt"
	"time"
)

// 27. Что будет выведено в результате выполнения следующего кода?

// Ответы:
// Числа от 0 до 9 без повторений

func main() {
	for i:=0;i<10;i++ {
		go fmt.Println(i)
	}

	time.Sleep(time.Second)
}

