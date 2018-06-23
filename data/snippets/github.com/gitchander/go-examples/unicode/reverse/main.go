package main

import "fmt"

func StrRev(s string) string {
	rs := []rune(s)
	i, j := 0, len(rs)-1
	for i < j {
		rs[i], rs[j] = rs[j], rs[i]
		i, j = i+1, j-1
	}
	return string(rs)
}

func main() {
	s := ".укнараб учрев ежкат есв я а ,укнаразопс олатсв ецнлоС"
	fmt.Println(StrRev(s))
}
