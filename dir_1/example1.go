package main

import (
	"fmt"
)

type Artist struct {
	Name, Genre string
	Songs       int
}

func newReleaseNonMutable(a Artist) int {
	a.Songs++
	return a.Songs
}

func newReleaseMutable(a *Artist) int {
	a.Songs++
	return a.Songs
}

func main() {
	me := Artist{Name: "Nexus", Genre: "Black Metal", Songs: 10}

	fmt.Println("-------- \tNon Mutabe\t --------")

	fmt.Printf("Name: %s\nGenre: %s\nSongs: %d\n", me.Name, me.Genre, me.Songs)
	fmt.Printf("%s released their %d songs\n", me.Name, newReleaseNonMutable(me))
	fmt.Printf("%s has in total %d songs\n", me.Name, me.Songs)

	fmt.Println("----------------------------------------")

	fmt.Println("-------- \tMutable\t\t --------")

	me2 := &Artist{Name: "Nexus", Genre: "Death Metal", Songs: 20}

	fmt.Printf("Name: %s\nGenre: %s\nSongs: %d\n", me2.Name, me2.Genre, me2.Songs)
	fmt.Printf("%s released their %d songs\n", me2.Name, newReleaseMutable(me2))
	fmt.Printf("%s has in total %d songs\n", me2.Name, me2.Songs)

	fmt.Println("----------------------------------------")

}
