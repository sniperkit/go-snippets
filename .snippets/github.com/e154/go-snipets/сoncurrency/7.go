package main

import (
	"fmt"
	"time"
)

//Буферизированный канал
//При инициализации канала можно использовать второй параметр
//Обычно каналы работают синхронно - каждая из сторон ждёт, когда другая сможет
//получить или передать сообщение. Но буферизованный канал работает асинхронно — получение или отправка
//сообщения не заставляют стороны останавливаться. Но канал теряет пропускную способность,
//когда он занят, в данном случае, если мы отправим в канал 1 сообщение, то мы не сможем отправить
//туда ещё одно до тех пор, пока первое не будет получено.

func pinger(c chan <- string) {
	for i := 0;;i++ {
		c <- "ping"
		time.Sleep(time.Second * 2)
	}
}

func ponger(c chan <- string) {
	for i := 0;;i++ {
		c <- "ponger"
		time.Sleep(time.Second * 3)
	}
}

func main() {

	//буферизированный канал, с ёмкостью в 20 сообщений
	var c1 chan string = make(chan string, 20)
	var c2 chan string = make(chan string, 20)

	go pinger(c1)
	go ponger(c2)

	go func() {
		for {
			select {
			case msg1 := <- c1:
				fmt.Println(msg1)

			case msg2 := <- c2:
				fmt.Println(msg2)

			case <- time.After(time.Second):
				fmt.Println("timeout")

//			default:
//				fmt.Println("nothing ready")

			}

			time.Sleep(time.Second * 1)
		}
	}()

	var input string
	fmt.Scanln(&input)
}

