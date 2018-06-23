package trie

import (
	"errors"
	"fmt"
)

type tree struct {
	root node
}

func NewTree() *tree {
	return &tree{root: node{}}
}

/// Finds associated value for key
/// returns error if key was not found in trie
func (t *tree) Find(key string) (interface{}, error) {
	return t.root.find(key)
}

/// Insert key-value-pair into trie
func (t *tree) Insert(key string, value interface{}) {
	t.root.insert(key, value)
}

/// If the tree contains the given value, the associated value is removed from the trie
/// returns true if key was deleted successfully
func (t *tree) Delete(key string) bool {
	return t.root.delete(key)
}

//Checks if the given key is associated with a value in the trie
func (t *tree) Contains(key string) bool {
	v, e := t.root.find(key)
	return e == nil && v != nil
}

/// Calls given function for every key-value-pair
/// argument:
/// 	- f func(string,interface{}) : function, which is called with key and value as arguments
func (t *tree) ForEach(f func(string, interface{})) {
	t.root.iterate("", f)
}

/// Returns all existing keys in trie as array of strings
func (t *tree) Keys() []string {
	return t.root.iterate("", nil)
}

/// Returns all key-value-pairs as map (golang type)
func (t *tree) ToMap() map[string]interface{} {
	m := map[string]interface{}{}
	t.root.iterate("", func(key string, value interface{}) {
		m[key] = value
	})
	return m
}

/// Returns all existing values in trie
func (t *tree) Values() []interface{} {
	values := []interface{}{}
	t.root.iterate("", func(key string, value interface{}) {
		values = append(values, value)
	})
	return values
}

/// Returns the amount of key-value-pairs in the trie as integer
func (t *tree) Size() int {
	return len(t.Keys())
}

// Converts trie to string e.g {test=2,go=test,key3={2,5}}
func (t *tree) String() string {
	result := "{"
	i := 0
	mLen := t.Size()
	t.ForEach(func(key string, value interface{}) {
		result += key + "=" + fmt.Sprintf("%v", value)
		if i+1 < mLen {
			result += ","
		}
		i++
	})
	return result + "}"
}

type node struct {
	Value interface{} //abstract value
	chars []byte      // character array
	tries []node      // children
}

func (t *node) find(key string) (interface{}, error) {
	if len(key) == 0 {
		return t.Value, nil
	} else {
		for i := 0; i < len(t.chars); i++ {
			if t.chars[i] == key[0] {
				return t.tries[i].find(string(key[1:]))
			}
		}
		return 0, errors.New("Key not found")
	}
}

//sorts a nodes chars and tries alphabetically with selection-sort
func (t node) sort() {
	for j := 0; j < len(t.chars); j++ {
		for i := 0; i < len(t.chars)-1; i++ {
			if t.chars[i] > t.chars[i+1] {
				t.chars[i], t.chars[i+1] = t.chars[i+1], t.chars[i]
				t.tries[i], t.tries[i+1] = t.tries[i+1], t.tries[i]
			}

		}
	}
}

func (t *node) insert(key string, value interface{}) {
	if len(key) == 0 {
		t.Value = value
	} else {
		for i := 0; i < len(t.chars); i++ {
			if t.chars[i] == key[0] {
				t.tries[i].insert(string(key[1:]), value)
				return
			}
		}
		newTrie := node{}
		newTrie.insert(string(key[1:]), value)
		t.chars = append(t.chars, key[0])
		t.tries = append(t.tries, newTrie)
		t.sort()
	}
}

func (t *node) delete(key string) bool {
	if len(key) == 0 {
		t.Value = nil
		return true
	} else {
		for i := 0; i < len(t.chars); i++ {
			if t.chars[i] == key[0] {
				result := t.tries[i].delete(string(key[1:]))
				if t.tries[i].isEmpty() {
					removeIthChar(t.chars, t.chars[i])
					removeIthTrie(t.tries, t.tries[i])
				}
				return result
			}
		}
	}
	return false
}

/// iterates recursively through tree and calls the given function for every found key-value-pair in trie
/// arguments:
/// 	- pre string: 			The previous keys as string to create the key as string
///	- f func(string,interface{}): 	The function to be called for every pair
func (t *node) iterate(pre string, f func(string, interface{})) []string {
	arr := []string{}
	if t.Value != nil {
		arr = append(arr, pre)
		if f != nil {
			f(pre, t.Value)
		}
	}
	if len(t.chars) > 0 {
		for i, e := range t.chars {
			arr = append(arr, t.tries[i].iterate(pre+string(e), f)...)
		}
	}
	return arr
}

// Checks if node is not filled yet
func (t *node) isEmpty() bool {
	return t.Value == nil && len(t.chars) == 0
}

//removes first occurrence of value in array
func removeIthChar(list []byte, val byte) {
	index := -1
	for i, e := range list {
		if &val == &e {
			index = i
			break
		}
	}
	if index > 0 {
		list = append(list[:index], list[index+1:]...)
	}
}

//removes first occurrence of value in array
func removeIthTrie(list []node, val node) {
	index := -1
	for i, e := range list {
		if &val == &e {
			index = i
			break
		}
	}
	if index > 0 {
		list = append(list[:index], list[index+1:]...)
	}
}
