//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//	 Unless required by applicable law or agreed to in writing, software
//	 distributed under the License is distributed on an "AS IS" BASIS,
//	 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	 See the License for the specific language governing permissions and
//	 limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package patriciaDB

import (
	"io"
	"sort"
)

type childList interface {
	length() int
	head() *Trie
	add(child *Trie) childList
	replace(b byte, child *Trie)
	remove(child *Trie)
	next(b byte) *Trie
	nextWithLongestPrefixMatch(b byte) (trie *Trie, exact bool)
	walk(prefix *Prefix, visitor VisitorFunc) error
	walkAndUpdate(prefix *Prefix, visitor UpdateFunc, handle Item) error
	print(w io.Writer, indent int)
	//total() int
}

type tries []*Trie

func (t tries) Len() int {
	return len(t)
}

func (t tries) Less(i, j int) bool {
	strings := sort.StringSlice{string(t[i].prefix), string(t[j].prefix)}
	return strings.Less(0, 1)
}

func (t tries) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

type sparseChildList struct {
	children tries
}

func newSparseChildList(maxChildrenPerSparseNode int) childList {
	return &sparseChildList{
		children: make(tries, 0, maxChildrenPerSparseNode),
	}
}

func (list *sparseChildList) length() int {
	return len(list.children)
}

func (list *sparseChildList) head() *Trie {
	return list.children[0]
}

func (list *sparseChildList) add(child *Trie) childList {
	// Search for an empty spot and insert the child if possible.
	//if len(list.children) != cap(list.children) {
	list.children = append(list.children, child)
	return list
	//}

	// Otherwise we have to transform to the dense list type.
	//return newDenseChildList(list, child)
}

func (list *sparseChildList) replace(b byte, child *Trie) {
	// Seek the child and replace it.
	for i, node := range list.children {
		if node.prefix[0] == b {
			list.children[i] = child
			return
		}
	}
}

func (list *sparseChildList) remove(child *Trie) {
	for i, node := range list.children {
		if node.prefix[0] == child.prefix[0] {
			list.children, list.children[len(list.children)-1] =
				append(list.children[:i], list.children[i+1:]...),
				nil
			return
		}
	}

	// This is not supposed to be reached.
	panic("removing non-existent child")
}

func (list *sparseChildList) nextWithLongestPrefixMatch(b byte) (trie *Trie, exact bool) {
	//logger.Println("Looking for byte ", b)
	var longestPrefixChild *Trie = nil
	for _, child := range list.children {
		//logger.Println("Scanning child ", child.prefix, " child.prefix[0] = ", child.prefix[0])
		if child != nil && len(child.prefix) > 0 && child.prefix[0] == b {
			//logger.Println("returning child ", child.prefix, "exact byte match")
			return child, true
		}
		if child != nil && len(child.prefix) > 0 && uint(child.prefix[0]) < uint(b) {
			//logger.Println("potential child ", child.prefix, " a potential match")
			if longestPrefixChild == nil || (uint(longestPrefixChild.prefix[0]) < uint(child.prefix[0])) {
				longestPrefixChild = child
			}
		}
	}
	return longestPrefixChild, exact
}

func (list *sparseChildList) next(b byte) *Trie {
	for _, child := range list.children {
		if child.prefix[0] == b {
			return child
		}
	}
	return nil
}

func (list *sparseChildList) walkAndUpdate(prefix *Prefix, visitor UpdateFunc, handle Item) error {

	sort.Sort(list.children)
	for i := 0; i < len(list.children); i++ {
		//for _, child := range list.children {
		child := list.children[i]
		*prefix = append(*prefix, child.prefix...)
		curr_len := len(list.children)
		if child.item != nil {
			err := visitor(*prefix, child.item, handle)
			if err != nil {
				if err == SkipSubtree {
					*prefix = (*prefix)[:len(*prefix)-len(child.prefix)]
					continue
				}
				*prefix = (*prefix)[:len(*prefix)-len(child.prefix)]
				return err
			}
		}

		err := child.children.walkAndUpdate(prefix, visitor, handle)
		*prefix = (*prefix)[:len(*prefix)-len(child.prefix)]
		if err != nil {
			return err
		}
		if len(list.children) < curr_len { //the current node was deleted
			i = i - 1
		}
	}

	return nil
}

func (list *sparseChildList) walk(prefix *Prefix, visitor VisitorFunc) error {

	sort.Sort(list.children)

	for _, child := range list.children {
		*prefix = append(*prefix, child.prefix...)
		if child.item != nil {
			err := visitor(*prefix, child.item)
			if err != nil {
				if err == SkipSubtree {
					*prefix = (*prefix)[:len(*prefix)-len(child.prefix)]
					continue
				}
				*prefix = (*prefix)[:len(*prefix)-len(child.prefix)]
				return err
			}
		}

		err := child.children.walk(prefix, visitor)
		*prefix = (*prefix)[:len(*prefix)-len(child.prefix)]
		if err != nil {
			return err
		}
	}

	return nil
}

/*
func (list *sparseChildList) total() int {
	tot := 0
	for _, child := range list.children {
		if child != nil {
			tot = tot + child.total()
		}
	}
	return tot
}
*/
func (list *sparseChildList) print(w io.Writer, indent int) {
	for _, child := range list.children {
		if child != nil {
			child.print(w, indent)
		}
	}
}

type denseChildList struct {
	min      int
	max      int
	children []*Trie
}

func newDenseChildList(list *sparseChildList, child *Trie) childList {
	var (
		min int = 255
		max int = 0
	)
	for _, child := range list.children {
		b := int(child.prefix[0])
		if b < min {
			min = b
		}
		if b > max {
			max = b
		}
	}

	b := int(child.prefix[0])
	if b < min {
		min = b
	}
	if b > max {
		max = b
	}

	children := make([]*Trie, max-min+1)
	for _, child := range list.children {
		children[int(child.prefix[0])-min] = child
	}
	children[int(child.prefix[0])-min] = child

	return &denseChildList{min, max, children}
}

func (list *denseChildList) length() int {
	return list.max - list.min + 1
}

func (list *denseChildList) head() *Trie {
	return list.children[0]
}

func (list *denseChildList) add(child *Trie) childList {
	b := int(child.prefix[0])

	switch {
	case list.min <= b && b <= list.max:
		if list.children[b-list.min] != nil {
			panic("dense child list collision detected")
		}
		list.children[b-list.min] = child

	case b < list.min:
		children := make([]*Trie, list.max-b+1)
		children[0] = child
		copy(children[list.min-b:], list.children)
		list.children = children
		list.min = b

	default: // b > list.max
		children := make([]*Trie, b-list.min+1)
		children[b-list.min] = child
		copy(children, list.children)
		list.children = children
		list.max = b
	}

	return list
}

func (list *denseChildList) replace(b byte, child *Trie) {
	list.children[int(b)-list.min] = nil
	list.children[int(child.prefix[0])-list.min] = child
}

func (list *denseChildList) remove(child *Trie) {
	i := int(child.prefix[0]) - list.min
	if list.children[i] == nil {
		// This is not supposed to be reached.
		panic("removing non-existent child")
	}
	list.children[i] = nil
}

func (list *denseChildList) next(b byte) *Trie {
	i := int(b)
	if i < list.min || list.max < i {
		return nil
	}
	return list.children[i-list.min]
}
func (list *denseChildList) nextWithLongestPrefixMatch(b byte) (trie *Trie, exact bool) {
	return nil, exact
}
func (list *denseChildList) walkAndUpdate(prefix *Prefix, visitor UpdateFunc, handle Item) error {
	for _, child := range list.children {
		if child == nil {
			continue
		}
		*prefix = append(*prefix, child.prefix...)
		if child.item != nil {
			if err := visitor(*prefix, child.item, handle); err != nil {
				if err == SkipSubtree {
					*prefix = (*prefix)[:len(*prefix)-len(child.prefix)]
					continue
				}
				*prefix = (*prefix)[:len(*prefix)-len(child.prefix)]
				return err
			}
		}

		err := child.children.walkAndUpdate(prefix, visitor, handle)
		*prefix = (*prefix)[:len(*prefix)-len(child.prefix)]
		if err != nil {
			return err
		}
	}

	return nil
}

func (list *denseChildList) walk(prefix *Prefix, visitor VisitorFunc) error {
	for _, child := range list.children {
		if child == nil {
			continue
		}
		*prefix = append(*prefix, child.prefix...)
		if child.item != nil {
			if err := visitor(*prefix, child.item); err != nil {
				if err == SkipSubtree {
					*prefix = (*prefix)[:len(*prefix)-len(child.prefix)]
					continue
				}
				*prefix = (*prefix)[:len(*prefix)-len(child.prefix)]
				return err
			}
		}

		err := child.children.walk(prefix, visitor)
		*prefix = (*prefix)[:len(*prefix)-len(child.prefix)]
		if err != nil {
			return err
		}
	}

	return nil
}

func (list *denseChildList) print(w io.Writer, indent int) {
	for _, child := range list.children {
		if child != nil {
			child.print(w, indent)
		}
	}
}

/*
func (list *denseChildList) total() int {
	tot := 0
	for _, child := range list.children {
		if child != nil {
			tot = tot + child.total()
		}
	}
	return tot
}*/
