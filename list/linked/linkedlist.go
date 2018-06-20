// Package linked implements a doubly linked list of the same typed items
package linked

const Prev = 0
const Next = 1

func Insert(link func(*) *[2]*, list **, elm *) {
	if nil == *list {
		*list = elm
		(*link(elm))[Prev] = elm
		(*link(elm))[Next] = elm

	} else if (*link(elm))[Prev] == nil {
		(*link(elm))[Prev] = *list
		(*link(elm))[Next] = (*link(*list))[Next]
		(*link(*list))[Next] = elm
		(*link((*link(elm))[Next]))[Prev] = elm
	} else {
		panic("One link cannot be in two lists")
	}
}

// add adds element to a list another element is already member of
func Add(link func(*) *[2]*, already *, elm *) {
	if (*link(already))[Prev] == nil || (*link(already))[Next] == nil {
		panic("Already is not already in the list")
	}
	(*link(elm))[Prev] = already
	(*link(elm))[Next] = (*link(already))[Next]
	(*link(already))[Next] = elm
	(*link((*link(elm))[Next]))[Prev] = elm
}

func Remove(link func(*) *[2]*, list **, elm *) {
	if *list == elm {
		if (*link(elm))[Prev] == elm {
			*list = nil
			goto finally
		} else {
			*list = (*link(elm))[Next]
		}
	}
	(*link((*link(elm))[Prev]))[Next] = (*link(elm))[Next]
	(*link((*link(elm))[Next]))[Prev] = (*link(elm))[Prev]
finally:
	(*link(elm))[Prev] = nil
	(*link(elm))[Next] = nil
}

func Empty(link func(*) *[2]*, list **) bool {
	return nil == *list
}

// do count of the list items
func Len(link func(*) *[2]*, list **) (count int) {

	if nil == *list {
		return 0
	}

	var e *
	e = (*link(*list))[Next]

	for e != *list {
		e = (*link(e))[Next]
		count++
	}
	count++
	return count
}

// apply function to all link elements
func Foreach(link func(*) *[2]*, list **, direction byte, f func(*)) {

	if nil == *list {
		return
	}

	var end *
	end = *list

	var e *
	e = (*link(*list))[direction]
	f(end)

	for (e != end) && ((*link(e))[direction] != nil) {
		var newe = (*link(e))[direction]
		f(e)
		e = newe
	}
	return
}
