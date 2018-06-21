/*
<!--
Copyright (c) 2017 Christoph Berger. Some rights reserved.

Use of the text in this file is governed by a Creative Commons Attribution Non-Commercial
Share-Alike License that can be found in the LICENSE.txt file.

Use of the code in this file is governed by a BSD 3-clause license that can be found
in the LICENSE.txt file.

The source code contained in this file may import third-party source code
whose licenses are provided in the respective license files.
-->

<!--
NOTE: The comments in this file are NOT godoc compliant. This is not an oversight.

Comments and code in this file are used for describing and explaining a particular topic to the reader. While this file is a syntactically valid Go source file, its main purpose is to get converted into a blog article. The comments were created for learning and not for code documentation.
-->

+++
title = "Go slices are not dynamic arrays"
description = "Go slices are based on a smart concept that does not like being ignored"
author = "Christoph Berger"
email = "chris@appliedgo.net"
date = "2017-08-03"
draft = "false"
domains = ["Patterns and Paradigms"]
tags = ["slice", "append", "split", "memory management", "gotcha"]
categories = ["Background"]
+++

Go's slices are cleverly designed. They provide the look-and-feel of truly dynamic arrays while being optimized for performance. However, not being aware of the slice mechanisms can bring you into trouble.

<!--more-->

*Background: Just recently I observed a few discussions--again--about seemingly inconsistent behavior of slice operations. I take this as an opportunity to talk a bit about slice internals and the mechanics around slice operations, especially `append()` and `bytes.Split()`.*


## Go's slices

The concept of slices in Go is really a clever one. A slice represents a flexible-length array-like data type while providing full control over memory allocations.

This concept is not seen in other languages, and so people new to Go often consider the behavior of slice operations as quite confusing. (Believe me, it happened to me as well.) Looking at the inner workings of slices removes much (if not all) of the confusion, so let's first have a look at the basics: What are slices, and how do they work?


## A slice is just a view on an array

In Go, arrays have a fixed size. The size is even part of the definition of an array, so the two arrays `[10]int` and `[20]int` are not just two `int` arrays of different size but are in fact different types.

Slices add a dynamic layer on top of arrays. Creating a slice from an array  neither allocates new memory nor copies anything. A slice is nothing but a "window" to some part of the array. Technically, a slice can be seen as a struct with a pointer to the array element where the slice starts, and two ints describing length and capacity.

This means that typical slice manipulations are cheap. Creating a slice, expanding it (as far as the available capacity permits), moving it back and forth on the underlying array--all that requires nothing more than changing the pointer value and/or one or both of the two int values. The data location does not change.

HYPE[slice basics](slices01.html)

*Fig.1: Slices are just "windows" to an array (click the buttons to see the operations)*

This also means that two slices created from the same array can overlap, and after assigning a slice to a new slice variable, both variables now share the same memory cells. Changing one item in one of the slices also change the same item in the other slice. If you want to create a true copy of a slice, create a new slice and use the built-in function `copy()`.

All of this is based on simple and consistent mechanisms. The problems arise when not being aware of these mechanisms.


## Some slice functions work in place

Since slices are just efficient "dynamic windows" on static arrays, it does make sense that most slice manipulations also happen in place.

As an example, `bytes.Split()` takes a slice and a separator, splits the slice by the separator, and returns a slice of byte slices.

But: All the byte slices returned by `Split()` still point to the same underlying array as the original slice. This may come unexpected to many who know similar split functions from other languages that rely on allocate-and-copy semantics (at the expense of efficiency at runtime).

HYPE[split](slices02.html)

*Fig. 2: `bytes.Split()` is an in-place operation*

Code that ignores the fact that the result of `Split()` still points to the original data may cause data corruption in a way that neither the compiler nor the runtime can detect as being wrong.

Another unexpected behavior can happen when combining `bytes.Split()` and `append()` - but first, let's have a look at `append()` alone.


## append() adds convenience--and some "magic"

`append()` adds new elements to the end of a slice, thus expanding the slice. `append()` has two convenience features:

* First, it can append to a `nil` slice, making it spring into existence in the moment of appending.
*  Second, if the remaining capacity is not sufficient for appending new values, `append()` automatically takes care of allocating a new array and copying the old content over.

Especially the second one can cause confusion, because after an `append()`, sometimes the original array has been changed, and sometimes a new array has been created, and the original one stays the same. If the original array was referenced by different parts of the code, one reference then may point to stale data.

HYPE[split](slices03.html)

*Fig. 3: The two outcomes of `append()`*

This behavior could be easily characterized as "random", although the behavior is in fact quite deterministic. An observer who always knows the values of slice length, capacity, and the number of items to append can trivially determine whether `append()` needs to allocate a new array.

In combination with `bytes.Split()`, `append()` can also create unexpected results. The slices that `bytes.Split()` returns have their `cap()` set to the end of the underlying array. Now when `append()`ing to the first of the returned slices, the slice grows within the same underlying array, overwriting subsequent slices.

HYPE[split and append](slices04.html)

*Fig. 4: After splitting (see fig. 2), append to the first returned slice*

If `bytes.Split()` returned all slices with their capacity set to their length, `append()` would not be able to overwrite subsequent slices, as it would immediately allocate a new array, to be able to extend beyond the slice's current capacity.

## A few demos

The code below demonstrates the discussed `Split()` and `append()` scenarios. It also shows how to do achieve an "always copy" semantics when appending.
*/

//
package main

import (
	"bytes"
	"fmt"
)

// Split the byte slice `a` at each comma, then update one of the split slices.
func splitDemo() {
	fmt.Println("Split demo")
	// bytes.Split splits in place.
	a := []byte("a,b,c")
	b := bytes.Split(a, []byte(","))
	fmt.Printf("a before changing b[0][0]: %q\n", a)

	// `b``'s byte slices use `a``'s underlying array. Changing `b[0][0]` also changes `a`.
	b[0][0] = byte('*')
	fmt.Printf("a after changing b[0][0]:  %q\n", a)

	// Appending to slice `b[0]` can write into slices `b[1]` and even `b[2], as `b[0]`'s capacity extends until the end of the underlying array that all slices share.
	fmt.Printf("b[1] before appending to b[0]: %q\n", b[1])
	b[0] = append(b[0], 'd', 'e', 'f')
	fmt.Printf("b[1] after appending to b[0]:  %q\n", b[1])
	fmt.Printf("a after appending to b[0]: %q\n", a)
}

// Append numbers to a slice; first, within capacity, then beyond capacity.
func appendDemo() {
	fmt.Println("\nAppend demo")
	s1 := make([]int, 2, 4)
	s1[0] = 1
	s1[1] = 2
	fmt.Printf("Initial address and value: %p: %[1]v\n", s1)
	s1 = append(s1, 3, 4)
	// Note the same address as before.
	fmt.Printf("After first append:        %p: %[1]v\n", s1)
	s1 = append(s1, 5)
	// Note the changed address. Append allocated a new, larger array.
	fmt.Printf("After second append:       %p: %[1]v\n", s1)
}

// How to get "always copy" semantics: simply `copy()` the slice before appending. Ensure the target slice is large enough for the subsequent `append()`, or else `append()` might again allocate a new array.
func alwaysCopy() {
	fmt.Println("\nAppend and always copy")
	s1 := []int{1, 2, 3, 4}
	fmt.Printf("s1: %p: %[1]v\n", s1)
	// Create a new slice with sufficient len (for copying) and cap (for appending - to avoid allocating and copying twice).
	s2 := make([]int, 4, 8)
	// Destination is always the first parameter, analogous to Fprintf, http.HandleFunc, etc.
	copy(s2, s1)
	// Note the different addresses of s1 and s2 in the output.
	fmt.Printf("s2: %p: %[1]v\n", s2)
	s2 = append(s2, 5, 6, 7, 8)
	// s2 has enough capacity so that append() does not allocate again.
	fmt.Printf("s2: %p: %[1]v\n", s2)
}

func main() {
	splitDemo()
	appendDemo()
	alwaysCopy()
}

/*

Output:

```
Split demo
a before changing b[0][0]: "a,b,c"
a after changing b[0][0]:  "*,b,c"
b[1] before appending to b[0]: "b"
b[1] after appending to b[0]:  "e"
a after appending to b[0]: "*defc"

Append demo
Initial address and value: 0xc42000a340: [1 2]
After first append:        0xc42000a340: [1 2 3 4]
After second append:       0xc420012380: [1 2 3 4 5]

Append and always copy
s1: 0xc42000a3c0: [1 2 3 4]
s2: 0xc4200123c0: [1 2 3 4]
s2: 0xc4200123c0: [1 2 3 4 5 6 7 8]
```


## How to get and run the code

Step 1: `go get` the code. Note the `-d` flag that prevents auto-installing
the binary into `$GOPATH/bin`.

    go get -d github.com/appliedgo/slices

Step 2: `cd` to the source code directory.

    cd $GOPATH/src/github.com/appliedgo/slices

Step 3. Run the code.

    go run slices.go


## Takeaways

### Remember that append() may or may not allocate a new slice.

In many cases, this is absolutely ok, as a single slice does not care if it gets relocated. Only when two or more slices interact, the behavior of `append()` can lead to unexpected results.

**To avoid ambiguous results, use the correct techniques to ensure the desired outcome:**

* If you absolutely want to avoid allocation and copying, use a large underlying array, re-slice your slice, and strictly avoid `append()`.

* If you absolutely need copy semantics, create a destination slice of sufficient size, and use the built-in `copy()` function.


### Never assume exclusive ownership of a slice that you did not create.

Any function that returns a slice may return a *shared* slice. `bytes.Split` splits a slice in-place, `append()` returns a slice header that still might point to the slice that it received before.

Hence if you receive a slice from a function, keep in mind that other code may still modify that slice. Again, `copy()` is your friend.


### Read the docs.

Functions that create or return copies of slices usually mention this in their documentation:

*"...returns a copy of..."*

*"...returns a new byte slice..."*

Whereas the documentation of in-place operations often talks about *"slicing"* or *"subslices"*, which indicates that no allocation takes place and the returned data may still be accessed by other code.


**Happy coding!**

- - -

Update 2017-08-05:

* Fixed fig. 3 to correctly show that len(t) == 3
* Added new case: Split then append

- - -

*/
