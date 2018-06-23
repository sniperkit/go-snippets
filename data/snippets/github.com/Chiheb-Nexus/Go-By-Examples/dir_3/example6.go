package main

import (
	"fmt"
)

var (
	coins = 50
	users = []string{
		"Matthew", "Sarah", "Augustus", "Heidi", "Emilie",
		"Peter", "Giana", "Adriano", "Aaron", "Elizabeth",
	}
)

func CoinsForUsers(users []string) (map[string]int, int) {

	coins_per_user := make(map[string]int, len(users))

	returnTotal := func(name string) int {
		var total int
		for i := 0; i < len(name); i++ {
			switch string(name[i]) {
			case "a", "A":
				total++
			case "e", "E":
				total++
			case "i", "I":
				total = total + 2
			case "o", "O":
				total = total + 3
			case "u", "U":
				total = total + 4
			}
		}
		return total
	}

	for _, name := range users {
		v := returnTotal(name)
		if v > 10 {
			v = 10
		}

		coins_per_user[name] = v
		coins = coins - v
	}

	return coins_per_user, coins
}

func main() {
	coins_per_user, coin := CoinsForUsers(users)

	for name, u_coin := range coins_per_user {
		fmt.Printf("%s:%d\n", name, u_coin)
	}
	fmt.Printf("Final coins: %d\n", coin)
}
