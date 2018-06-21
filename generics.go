package main //go:generate genny -in=capsule/capsule.go -out=capsule/uint32capsule.go gen "Item=uint32"
import (
	"fmt"
	"reflect"
)

/*
<!--
Copyright (c) 2016 Christoph Berger. Some rights reserved.
Use of this text is governed by a Creative Commons Attribution Non-Commercial
Share-Alike License that can be found in the LICENSE.txt file.

The source code contained in this file may import third-party source code
whose licenses are provided in the respective license files.
-->

<!--
NOTE: The comments in this file are NOT godoc compliant. This is not an oversight.

Comments and code in this file are used for describing and explaining a particular topic to the reader. While this file is a syntactically valid Go source file, its main purpose is to get converted into a blog article. The comments were created for learning and not for code documentation.
-->

+++
title = "Who needs generics? Use ... instead!"
description = "The Go language has no generics. This article is a survey of techniques that can be used instead."
author = "Christoph Berger"
email = "chris@appliedgo.net"
date = "2016-07-14"
publishdate = "2016-07-14"
domains = ["Patterns and Paradigms"]
tags = ["generics", "interface", "reflection", "code generation"]
categories = ["Essay"]
+++

What are generics? Why are they considered useful? Why does Go have no generics? What shall Gophers use instead?

This article examines the nature of generics, and surveys various techniques that can be used to work around the absence of this programming paradigm.

<!--more-->

- - -

**UPDATE: The "..." in the title is meant to be a plain English ellipsis!** Big thanks to the readers who have pointed out that "..." in this context can (and will) be understood as Go's ellipsis. In this case the title would of course be complete rubbish.

**UPDATE 2:** The snappy title is often misunderstood as "Go needs no generics". This is far from what the article wants to say. I do see the usefulness of generics in certain problem domains, and I am the last one to balk at the idea of generics in Go. If anyone feels misled by the title, I apologize.
- - -

## First, an important note

_The question about generics in Go is years old, and has been discussed up and down and forth and back across the Go forums, newsgroups, and email lists. **It is NOT the goal of this article to re-ignite this discussion.** I think that all that can be said, has been said._

*A good summary of the state of the discussion is [here](https://docs.google.com/document/d/1vrAy9gMpMoS3uaVphB32uVXX4pi-HnNjkMEgyAHX4N4/edit?usp=sharing).*

***This article solely focuses on alternate ways of achieving some of the goals that other languages try to solve with generics.***

Let's start by briefly looking at the motivation for generics as well as the possible downsides.


## The problem

Many data structures and algorithms are applicable to a range of data types. A sorted tree, for example, could be defined to hold elements of type `int`, or `map[string]string`, or some `struct` type. A sorting algorithm is able to sort any type whose elements are comparable to each other. So if you need, say, a sorted tree for strings, you could sit down and write one. Easy enough!

But what if you need a sorted tree for many different types? The tree datatype comes with a couple of methods like `Insert`, `Find`, `Delete`, etc. If there are *N* tree methods, and *M* element types to implement the tree for, you would end up with *N x M* methods to implement! Sounds like days of tedious, repetitive work.

Can't we let the compiler do that?


## Enter generics

To address this problem, many programming languages have a concept called *'generics'*. The point is to write the code only once, in an abstract form, with opaque type placeholders instead of real types, as in this Java example:

```Java
// Java pseudo code

public class SortedTree<E implements Comparable<E>> {
	void insert(E comparableSortTreeElement) {
		//...
	}
}
```

`E` is a type placeholder; it can appear in class or interface definitions, as well as in function parameter lists. `SortedTree` is used by substituting a real type for the placeholder:

```Java
// Java pseudo code (again)

SortedTree<String> sortedTreeOfStrings = new SortedTree<String>();
sortedTreeOfStrings.insert("abcde");
```

These lines instantiate a sorted tree with string elements. The `insert` method then only accepts strings, and the sort algorithm uses string comparison methods. Had we used `Integer` instead, like in `SortedTree<Integer>`, the tree would be a sorted tree of integers instead.

To summarize, generics allow to define classes and functions with type parameters. The type parameters can be restricted to a certain subset (e.g., parameter T must be comparable), but otherwise they are unspecific. The result is a code template that can be applied to different types as the need arises.


## The downsides

While generics may come in handy, they also have some strings attached.

1. **Performance:** Turning a generic code template into actual code takes time, either at compile time or at runtime. Or as [Russ Cox stated](http://research.swtch.com/generic) in 2009:

    > The generic dilemma is this: *do you want slow programmers, slow compilers and bloated binaries, or slow execution times?*

	(The 'slow programmers' part refers to having no generics at all, nor any suitable substitute.)

2. **Complexity**: Generics are not complex *per se*, but they can become complex when integrated with other language features such as inheritance, or when developers start creating nested generics like,

	```C#
	// C#
	List<Dictionary<string<IEnumerable<HttpRequest>>>>
	```

	(Thanks to [Jonathan Oliver](http://blog.jonathanoliver.com/golang-has-generics/) for this nice example), or even seemingly recursive inheritance like,

	```Java
	// Java
	public abstract class Enum<E extends Enum<E>>
	```

	(taken from [this Java Generics FAQ](http://www.angelikalanger.com/GenericsFAQ/FAQSections/TypeParameters.html#FAQ106) by Angelika Langer).


## Go has no generics

As we all know, Go has deliberately been designed with simplicity in mind, and generics are considered to add complexity to a language (see the previous section). So along with inheritance, polymorphism, and some other features of the 'state of the art' languages at that time, generics were left off from Go's feature list.

### Actually, Go does have *some* generics--sort of...

There are indeed a few 'generic' language constructs in Go.

First, there are three generic data types you can make use of (and probably already did so without noticing):

* arrays
* slices
* maps

All of these can be instantiated on arbitrary element types. For the `map` type, this is even true for both the key and the value. This makes maps quite versatile. For example, `map[string]struct{}` can be used as a Set type where every element is unique.

Second, some built-in functions operate on a range of data types, which makes them almost act like generic functions. For example, `len()` works with strings, arrays, slices, and maps.

Although this is definitely not 'the real thing', chances are that those 'internal generics' already cover your needs.


## But what if a project just seems to cry for generics?

At times you might come across a problem where generics would seem to be very handy, if not downright indispensable. What can you do then?

Luckily, there are a couple of options.


### 1. Review the requirements

Step back and revisit the requirements. Review the technical or functional specification (you should have one). Do the specs really demand the use of generics? Consider that while other languages may support a design that is based on type systems, Go's philosophy is different:

> If C++ and Java are about type hierarchies and the taxonomy of types, Go is about composition. [(Rob Pike)](https://commandcenter.blogspot.de/2012/06/less-is-exponentially-more.html)

So think of the paradigms that come with Go--most notably composition and embedding--and verify if any of these would help approaching the problem in a way more natural to Go.



### 2. Consider copy & paste

This may sound like a foolish advice (and if applied improperly, it surely is), but do not hastily reject it.


Here is the point:

Every time you think that you need to create a generic object, do a quick Litmus test: Ask yourself, "How many times would I ever have to instantiate this generic object in my application or library?" In other words, is it worth to construct a generic object when there may only be one or two actual implementations of it? In this case, creating a generic object would just be an over-abstraction.

(At this point, I cannot help but thinking of Joel Spolsky's witty article on [Architecture Astronauts](http://www.joelonsoftware.com/articles/fog0000000018.html) who build abstractions on top of other abstractions until they run out of oxygen. (Small caveat: the article is from 2001. Expect a couple of references to outdated software concepts of which you never may have heard if you are young enough.)) *<-- Yes, I nested two parens here. So?*

An excellent real-life example can be found right in the standard library. The two packages `strings` and `bytes` have almost identical API's, yet no generics were used in the making of these packages (or so I strongly guess).

Remember that this article tries to remain neutral and takes no stand on generics. I don't want to indicate that copy&paste is better than generics. I only want to encourage you to seek out other options even if they look strange at a first glance.


### 3. Use interfaces

Interfaces define behavior without requiring any implementation details. This is ideal for defining 'generic' behavior:

* Find a set of basic operations that your generic algorithm or data container can use to process the data.
* Define an `interface` containing these operations.
* To 'instantiate' your generic entity on a given data type, implement the interface methods for that type.

The `sort` package is an example for this technique. It contains a sort interface (appropriately called `sort.Interface`) that declares three methods: `Len()`, `Less()`, and `Swap()`.

```go
type Interface interface {
	Len() int
	Less(i, j int) bool
	Swap(i, j int)
}
```

By just implementing these three methods for a data container (say, a slice of structs), `sort.Sort()` can be
applied to any kind of data, as long as the data has a reasonable definition of 'less than'.

A nice aspect here is that the code of `sort.Sort()` does not know anything about the data it sorts, and actually it does not have to. It simply relies on the three interface methods `Len`, `Less`, and `Swap`.

There are already a couple of examples in the [sort package doc](https://golang.org/pkg/sort/), and the following technique makes heavy use of an empty interface, so I'll skip the code part and move right on to type assertions.


### 4. Use type assertions

Generic container types usually do not need to care much about the actual type of their contents. (Except maybe for basic operations like sorting, but we already know a solution for that.) Hence the value can be stored in the container as a 'type without properties'. Such a type is already built into Go: the empty interface, declared as `interface{}`. This interface has no particular behavior, hence objects with any behavior satisfy this interface.

It is quite easy to build a container type based on `interface{}`. We only need a way to recover the actual data type when reading elements from the container. For that purpose, Go has type assertions. Here is an example that implements a generic container object.

To keep the code short and concise, the container object is kept super-simple. It features only two methods, Put and Get.

*/

// `Container` is a generic container, accepting anything.
type Container []interface{}

// Put adds an element to the container.
func (c *Container) Put(elem interface{}) {
	*c = append(*c, elem)
}

// Get gets an element from the container.
func (c *Container) Get() interface{} {
	elem := (*c)[0]
	*c = (*c)[1:]
	return elem
}

// The calling code does the type assertion when retrieving an element.
func assertExample() {
	intContainer := &Container{}
	intContainer.Put(7)
	intContainer.Put(42)
	elem, ok := intContainer.Get().(int) // assert that the actual type is int
	if !ok {
		fmt.Println("Unable to read an int from intContainer")
	}
	fmt.Printf("assertExample: %d (%T)\n", elem, elem)
}

/*

This looks so straightforward that we might easily forget the downsides of this technique. We give up compile-time type checking here, exposing the application to increased risk of type-related failure at runtime. Also, conversions to and from interfaces come with a cost. And finally, all 'magic' happens outside our Container type, at the caller's level. Usually you would rather want a technique that hides the type conversion mechanisms from the caller.


### 5. Use reflection

Remember the *generic dilemma* from the 'downsides' section? One of the three options to pick from was "slow execution times". The technique discussed below is a manifestation of this option. Yes, I am talking about reflection.

Reflection allows a program to work with objects whose types are not known at compile time. Putting performance concerns aside, this sounds like a good fit for solving the generics problem, so let's give it a try.

The code below implements a simple generic container. You can see a lot of type checking going on via the `reflect` package. Worth noting: A couple of standard operations do not work on variables of type `reflect.Value` even if the actual type would allow this operation. You need to substitute these operations with methods from `reflect`. I added comments to the parts of the code where this occurs.

*/

// Cabinet has one field, `s`, that holds a slice of a given type. (As the name `Container` is already taken, I had to choose another name.)
type Cabinet struct {
	s reflect.Value
}

// NewCabinet creates a new Cabinet struct where `s` holds a slice of type `[]t`. Formally, `s` remains a `reflect.Value` as defined in the struct. The code needs to deal with that fact in some places below.
func NewCabinet(t reflect.Type) *Cabinet {
	return &Cabinet{
		s: reflect.MakeSlice(reflect.SliceOf(t), 0, 10), // cap is arbitrary, we need to pass one here
	}
}

// Set sets the element at index i to the passed-in value.
func (c *Cabinet) Put(val interface{}) {
	// The passed-in `val` must match the type of the elements of slice `s`. Otherwise, Put panics, as this is a programmer error.
	if reflect.ValueOf(val).Type() != c.s.Type().Elem() {
		panic(fmt.Sprintf("Put: cannot put a %T into a slice of %s", val, c.s.Type().Elem()))
	}
	// `AppendSlice` is a replacement for the builtin `append` function, which fails on a `reflect.Value` even if the actual value's type is a slice.
	c.s = reflect.Append(c.s, reflect.ValueOf(val))
}

// Get gets the element at index `i`. There is no way (or so it seems) to have a function return a `reflect.Value` type that turns into the actual type of the returned data. Hence the Get function has a second parameter that must be a reference of the receiving variable. See `reflectExample()`.
func (c *Cabinet) Get(retref interface{}) {
	// `Index(i)` replaces the index operator `[i]` as `s` is only a reflect.Value (even though it effectively contains a slice).
	retref = c.s.Index(0)
	c.s = c.s.Slice(1, c.s.Len())
}

//
func reflectExample() {
	f := 3.14152
	g := 0.0
	c := NewCabinet(reflect.TypeOf(f))
	// Try c.Put(0, "blabla") to see the type check panicking
	c.Put(f)
	// The syntax `g = c.Get(0)` is not possible, see the comment on `Get()`.
	c.Get(&g)
	fmt.Printf("reflectExample: %f (%T)\n", g, g)
}

/*
Frankly, I really had a hard time wrapping my head around the semantics of the `reflect` package until I got this code running. 😓

And I realized that when moving from real types to the level of `reflect.Value` and `reflect.Type`, a couple of things are not possible anymore. For example, you need replacement functions for `append()`, the `[i]` operator, etc, and there is no way to return a `reflect.Value` and have it magically turn into the actual type that it contains. (That's why `Get()` needs to receive a reference to fill with the return value.)

I feel that I have put most of the time of writing this article into the Reflection code. A Go proverb says, "Clear is better than clever", but this code is not only anything but clear, it is also far from being clever.

My $0.02 on generics through reflection?

Avoid.


### 6. Use a code generator

Far better than reflection IMHO is code generation. The concept is rather simple:

* Write a template file that uses generic 'mock' types. These types are just placeholders for the real types.
* Run the file through a converter. The converter identifies the mock types and replaces them by real types.

To illustrate this, here is a quick walkthrough using [genny](https://github.com/cheekybits/genny). (There are also other converters with similar features out there, like [generic](https://github.com/taylorchu/generic), [gengen](https://github.com/joeshaw/gengen), [gen](https://github.com/clipperhouse/gen), and more.)

So let's write another generic container. Whenever we need a generic type, we can declare one like so:

```go
type APlaceholder generic.Type
```

`APlaceholder` then can be used in the template code like any other type. The template code compiles without problems, and it is even possible to write tests against this code. This is how our container template looks like:

```go
package capsule

import "github.com/cheekybits/genny/generic"

type Item generic.Type

type ItemCapsule struct {
	s []Item
}

func NewItemCapsule() *ItemCapsule {
	return &ItemCapsule{s: []Item{}}
}

func (c *ItemCapsule) Put(val Item) {
	c.s = append(c.s, val)
}

func (c *ItemCapsule) Get() Item {
	r := c.s[0]
	c.s = c.s[1:]
	return r
}
```
(When you `go get` the code of this article, this template code is in `capsule/capsule.go`.)

In the main file, `generics.go`, we want to have a Capsule containing uint32 elements. So we place a `go:generate` directive at the top of the file:

	//go:generate genny -in=template.go -out=uint32capsule.go gen "Item=uint32"

Before compiling the main binary, calling `go generate` executes genny, which produces an output file where `Item` is replaced by `uint32`. Note that this also applies to names that contain the string 'Item' - like 'ItemCapsule' in our example. This way, it is possible to create multiple replacements from the same template without name conflicts. (For example, a Uint32Capsule and a StringCapsule.)

After running `go generate generics.go`, we can find the file `uint32capsule.go` in the project dir with the following code:
*/

// Generated code
type Uint32Capsule struct {
	s []uint32
}

func NewUint32Capsule() *Uint32Capsule {
	return &Uint32Capsule{s: []uint32{}}
}

func (c *Uint32Capsule) Put(val uint32) {
	c.s = append(c.s, val)
}

func (c *Uint32Capsule) Get() uint32 {
	r := c.s[0]
	c.s = c.s[1:]
	return r
}

// Now we can write the calling function as if the Uint32Capsule was handcoded. Everything is just plain, clear, idiomatic Go.
func generateExample() {
	var u uint32 = 42
	c := NewUint32Capsule()
	c.Put(u)
	v := c.Get()
	fmt.Printf("generateExample: %d (%T)\n", v, v)
}

/* For completness of our test code, here is the `main` function.
 */

// `main` simply runs all examples.
func main() {
	assertExample()
	reflectExample()
	generateExample()
}

/*

Get the complete test code from GitHub:

		go get -d github.com/appliedgo/generics
		cd $GOPATH/src/github.com/appliedgo/generics
		go build
		./generics

You don't need to get `genny` nor run `go generate` if you just want to run the examples in `generics.go`. The code generated by `genny` is included in generics.go.


## Summary

Let me do a quick roundup of the examined techniques.

### Requirements review

* Pro: Makes you think about your design.
* Con: Likely to apply to only a small set of problem types.

### Copy&Paste

* Pros:
	* Quick
	* Needs no external libraries or tools
* Cons:
	* Code bloat
	* Breaks the [DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself) principle

### Interfaces (with methods)

* Pros:
	* Clean and idiomatic
	* Needs no external libraries or tools
* Cons:
	* Additional coding required - each interface method needs to be implemented for each target type.
	* Runtime overhead (however maybe not that dramatic)

### Type assertions (with empty interfaces)

* Pros:
	* Code stays quite clear
	* Needs no external libraries or tools
* Cons:
	* No compile-time type checking
	* Runtime overhead from converting to and from interfaces
	* Caller is required to do the type assertion

### Reflection

* Pros:
	* Versatile
	* Needs no external libraries or tools
* Cons:
	* Anything but clear
	* No compile-time type checking
	* Considerable runtime overhead

### Code generation

* Pros:
	* Very clear code possible (depending on the tool), both in the templates and in the generated code
	* Compile-time type checking
	* Some tools even allow writing tests against the generic code template
	* No runtime overhead
* Cons:
	* Possible binary bloat, if a template gets instantiated many times for different target types (but how often does this happen in your apps?)
	* The build process requires an extra step and a third party tool (but as far as I have seen, these tools usually are `go get`table)


That's it for this time. I hope you enjoyed reading, even though this article ended up being quite long and without any images (sorry for that). I try to keep it shorter in the future. Until next time!

*/
