package main

import (
	"sort"
	"fmt"
)

type Item struct {
	Name		string
	Items		[]string
}

type Items []*Item

func p(str string, items *Items){
	fmt.Println(str)
	for _, i := range *items {
		fmt.Println(i.Name)
	}
}

func (i Items) Len() int {
	return len(i)
}

func (i Items) Swap(a, b int) {
	i[a], i[b] = i[b], i[a]
}

func (i Items) Less(a, b int) bool {
	return len(i[a].Items) < len(i[b].Items)
}

func (i Items) Push(item *Item) {
	l := len(item.Items)
	if l == 0 {return}

	k := sort.Search(len(i), func(k int) bool {return len(i[k].Items) >= l})

	var q Items

	q = append(i[:k])
	q = append(q, item)
	q = append(q, i[k:len(i)] ...)

	p("q:", &q)
}

func main(){
	items := Items{
		&Item{"obj5", []string{"item1","item2","item3","item4","item5",}},
		&Item{"obj4", []string{"item1","item2","item3","item4",}},
		&Item{"obj2", []string{"item1", "item2",}},
		&Item{"obj1", []string{"item1",}},
	}

//	p("before:", &items)

	sort.Sort(items)
//	p("after:", &items)

	items.Push(&Item{"obj3", []string{"item1", "item2", "item3", }})

//	p("after:", &items)
}
