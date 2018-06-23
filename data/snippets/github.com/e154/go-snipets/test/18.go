package main

// 18. В какой строке ошибка

// Ответы:
//

type Incrementer interface {
	Inc()
}

type Num int

func (n *Num) Inc() {
	*n += 1
}

func Levelup(c Incrementer) {
	c.Inc()
}

func main() {
	var n Num = 0
	var np *Num = new(Num)

	//Levelup(&n)
	Levelup(n) // <-- тут ошибка
	Levelup(np)

	//fmt.Println(n)
	//fmt.Println(*np)
}

