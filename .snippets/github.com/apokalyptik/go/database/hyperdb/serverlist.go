package hyperdb

import (
	"math/rand"
	"sort"
)

// serverList is a map of priority to a slice of pointers to servers
type serverList map[int][]*Server

func (s serverList) get() *Server {
	// Build a list of priority numbers
	var priorities = []int{}
	for k := range s {
		if k == 0 {
			continue
		}
		priorities = append(priorities, k)
	}
	// Sort the list (ends up ascending)
	sort.Ints(priorities)
	// Loop over the list keys
	for k := range priorities {
		// Turn the ascending key number into a descending key number, and use that
		// to get the priority we want to work on first (because we want higher
		// priorities to be higher)
		var v = priorities[len(priorities)-(1+k)]
		// Build a list of server pointers for that priority
		var list = []*Server{}
		for _, sv := range s[v] {
			list = append(list, sv)
		}
		// Infinite loop, because we want to loop over the list in random order for
		// load distribution purposes
		for {
			// If our list is empty (started, or we ran out of possibilities) then we're
			// done.
			if len(list) < 1 {
				break
			}
			// Get a random number which falls within the range of valid keys for our
			// server pointer slice
			try := rand.Intn(len(list))
			if db := list[try].getDB(); db != nil {
				// If that one worked, then use it
				return list[try]
			}
			// Otherwise remove it from the possibilities and try again
			list = append(list[:try], list[try+1:]...)
		}
	}
	return nil
}
