commits := map[string]int{
	"rsc": 3771,
	"huhu": 5688,
	"gri": 987,
	"adhy": 4529,
}

m = map[string]int{}


// For instance, a map of boolean values can be used as a set-like data structure (recall that the zero value for the boolean type is false). 
// This example traverses a linked list of Nodes and prints their values. It uses a map of Node pointers to detect cycles in the list.


type Node struct {
	Next *Node
	Value interface{}
}

var first *Node

visited := make(map[*Node]bool)
for n := first; n != nil; n = n.Next {
	if visited[n] {
		fmt.Println("cycle detected")
		break
	}
	visited[n] = true
	fmt.Println(n.Value)
}

type Person struct {
	Name string
	Likes []string
}

var people []*Person

likes := make(map[string][]*Person)
for _,p := range people {
	for _, l := range p.Likes {
		likes[l] = append(likes[l], p)
	}
}

// To print a list of people who like cheese:

    for _, p := range likes["cheese"] {
        fmt.Println(p.Name, "likes cheese.")
    }
// To print the number of people who like bacon:

    fmt.Println(len(likes["bacon"]), "people like bacon.")

hits := make(map[string]map[string]int)

// This is map of string to (map of string to int). Each key of the outer map is the path to a web page with its own inner map. 
// Each inner map key is a two-letter country code. This expression retrieves the number of times an Australian has loaded the documentation page:
n := hits["/doc/"]["au"]

func add(m map[string]map[string]int, path, country string){
	mm, ok := m[path]
	if !ok {
		mm = make(map[string]int)
		m[path] = mm
	}
	mm[country]++
}
add(hist, "/doc/", "au")


