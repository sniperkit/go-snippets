package main

// 22. В каком фрагменте кода приведен правильный способ перехвата паники функции enterDangeon?

// Ответы:
//

func enterDangeo(level int) int {
	return level
}

func farmGold(level int) int {
	gold := enterDangeo(level)
	err := recover()
	if err != nil {
		return 0
	}

	return gold
}

//func farmGold2(level int) int {
//	var gold int
//	err := recover(func() {
//		gold = enterDangeo(level)
//	})
//}


func main() {

}

