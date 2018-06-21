/*
<!--
Copyright (c) 2016 Christoph Berger. Some rights reserved.
Use of this text is governed by a Creative Commons Attribution Non-Commercial
Share-Alike License that can be found in the LICENSE.txt file.

The source code contained in this file may import third-party source code
whose licenses are provided in the respective license files.
-->

+++
title = "Dependency Injection in a nutshell"
description = "How to implement the Dependency Rule in Go"
author = "Christoph Berger"
email = "chris@appliedgo.net"
date = "2016-06-23"
publishdate = "2016-06-23"
domains = ["Architecture"]
tags = ["Dependency Injection", "Separation Of Concerns", "Interface"]
categories = ["Tutorial"]
+++

Layered software architectures adhere to the *Dependency Rule:* Source code in a lower-level layer can make use of code in higher-level layers, but never vice versa. Control flow, however, goes in both directions. How is this possible, given that higher-level code must not know anything about the code in lower levels?

<!--more-->

## Software Architecture

Software architecture strives to provide structure to software systems, in order to make them robust, maintainable, extendable, testable, easier to develop, and easier to document.

Many different [architecture patterns](https://en.wikipedia.org/wiki/List_of_software_architecture_styles_and_patterns) have evolved over time, at different abstraction levels and with different levels of complexity. Among these architecture patterns, layered architectures seem to represent a fairly versatile concept that is applicable to a large range of scenarios.

Just to see how such an architecture may look like, let's have a brief look at the *Clean Architecture*, a layered architecture model that summarizes the idea of layering very well.


## The Clean Architecture

At the center of any layered software architecture is the separation of concerns. In simple words: The less each software module knows about the other modules, the better.

To achieve this, the modules are organized into layers. Each layer represents a certain level of abstraction. The Clean Architecture model describes them as (at least) four concentric circles, with the innermost circle representing the highest abstraction level.

![The Clean Architecture](cleanarchitecture.png)

(Discussing each layer in detail is outside the scope of this article. I briefly introduced the Clean Architecture here so that the dependency problem that is discussed below becomes clear. You can read more about The Clean Architecture [in this article](http://blog.8thlight.com/uncle-bob/2012/08/13/the-clean-architecture.html). Definitely recommended!)

The central rule of The Clean Architecture is the **Dependency Rule**, which says,

> Source code dependencies can only point inwards.

In other words, the source code of each circle can only access code in an inner circle but never any code in an outer circle.

But what is this good for?


## A small example without the Dependency Rule

As a completely made up and utterly pointless scenario, imagine a poet who writes, well, poems. Poems have to be stored somewhere, so the Chief Software Architect of ACME Poem Processing, Inc. comes up with this architecture:

* A top layer (or "inner ring") containing poem documents, and
* A bottom layer (or "outer ring") containing poem storage entities.

![A two-layer poem architecture](poemarchitecture.png)

(Granted, this is a rather simplified version of a layered architecture but for our scenario, it is just enough.)

A document object obviously needs to access the services of a storage object to store and retrieve its contents (blue arrow). Thus it would seem natural to add a storage service directly to the document.

In our example, our poet surely wants to write the poems into a small notebook, and thus the Lead Programmer creates this document layer:

```go
type Poem struct {
	content []byte
	storage acmeStorageServices.PoemNotebook
}

func NewPoem() *Poem {
	return &Poem {
		storage: acmeStorageServices.NewPoemNotebook(),
	}
}

func (p *Poem) Load(title string) {
	p.content = p.storage.Load(title)
}
func (p *Poem) Save(title string) {
	storage.Save(title, p.content)
}
```

Easy enough! But wait--what if our poet decides to write a poem on a napkin? Or on 4x6 index cards? **The document layer would have to be modified and recompiled!** We have created an unwanted dependency on a particular storage type.

How can we remove that dependency?


## Abstraction to the rescue

As a first step, we can replace the storage service by an abstraction of that service. Using Go's `interface` type, this becomes really easy.

```go
type PoemStorage interface {
	Load(string) []byte
	Save(string, []byte)
}
```

The interface describes only a behavior, and our Poem object can call the interface functions without worring about the object that implements this interface.

Now we can define the Poem struct without any dependency on the storage layer:

```go
type Poem struct {
	content []byte
	storage PoemStorage
}
```

Remember, `PoemStorage` is just an interface but we can assign any type to `storage` that satisfies this interface.

![Poem storage interface](poemabstraction.png)


## Adding dependency injection

Right now the Poem only talks to an empty abstraction. As the next step, we need a way to connect a real storage object to the Poem.

In other words, we need to *inject a dependency* on a PoemStorage object into the Poem layer.

We can do this, for example, through a constructor:

```go
func NewPoem(ps *PoemStorage) {
	return &Poem{
		storage: ps
	}
}
```
When called, the constructor receives an actual PoemStorage object, yet the returned Poem still just talks to the abstract PoemStorage interface.

Finally, in `main()` or in some dedicated setup function, we can wire up all higher-level objects with their lower-level dependencies.

```go
func main() {
	storage := NewNapkin()
	poem := NewPoem(storage)  // wired up.
}
```

Boom! We have just injected a dependency on a Napkin object into our new Poem object. To point it out again, at no point did the Poem object learn about the Napkin object, yet we just made it use one.


HYPE[Clean Poem Architecture with dependency injection](poem.html)

*This is the gist of dependency injection.* There is surely more to it than we were able to go through in this article. The interface/constructor pattern is not the only approach to implementing dependency injection. Still, it is a quite appealing one because it is clear and concise and builds upon just a few basic language constructs.


## Verba docent exempla trahunt

*Words teach, examples lead.* With this in mind let me finish this article with a working example.

(Note: The complete lack of error handling or any other kind of sanity checks is intentional for brevity's sake, yet it is anything but exemplary. If you think this sets a bad example for inexperienced readers, then you are probably right and I apologize. Dear inexperienced readers: Use proper error handling. Wherever you can. I am serious about this.)

*/

// ## Imports and globals
package main

import "fmt"

// ### The "inner ring"

// A `Poem` contains some poetry and an abstract storage reference.
type Poem struct {
	content []byte
	storage PoemStorage
}

// `PoemStorage` is just an interface that defines the behavior of a poem storage.
// This is all that `Poem` knows (and needs to know) about storing and retrieving poems.
// Nothing from the "outer ring" appears here.
type PoemStorage interface {
	Type() string        // Return a string describing the storage type.
	Load(string) []byte  // Load a poem by name.
	Save(string, []byte) // Save a poem by name.
}

// `NewPoem` constructs a `Poem` object. We use this constructor to inject an object
// that satisfies the `PoemStorage` interface.
func NewPoem(ps PoemStorage) *Poem {
	return &Poem{
		content: []byte("I am a poem from a " + ps.Type() + "."),
		storage: ps,
	}
}

// `Save` simply calls `Save` on the interface type. The `Poem` object neither knows
// nor cares about which actual storage object receives this method call.
func (p *Poem) Save(name string) {
	p.storage.Save(name, p.content)
}

// `Load` also invokes the injected storage object without knowing it.
func (p *Poem) Load(name string) {
	p.content = p.storage.Load(name)
}

// `String` makes Poem a Stringer, allowing us to drop it anywhere a string would be
// expected.
func (p *Poem) String() string {
	return string(p.content)
}

// ### The "outer ring"

// #### The notebook

// A `Notebook` is the classic storage device of a poet.
type Notebook struct {
	poems map[string][]byte
}

func NewNotebook() *Notebook {
	return &Notebook{
		poems: map[string][]byte{},
	}
}

// After adding `Save` and `Load`, `Notebook` implicitly satisfies `PoemStorage`.
func (n *Notebook) Save(name string, contents []byte) {
	n.poems[name] = contents
}

func (n *Notebook) Load(name string) []byte {
	return n.poems[name]
}

// `Type` returns an informal description of the storage type.
func (n *Notebook) Type() string {
	return "Notebook"
}

// A `Napkin` is the emergency storage device of a poet.
// It can store only one poem.
type Napkin struct {
	poem []byte
}

func NewNapkin() *Napkin {
	return &Napkin{
		poem: []byte{},
	}
}

func (n *Napkin) Save(name string, contents []byte) {
	n.poem = contents
}

func (n *Napkin) Load(name string) []byte {
	return n.poem
}

func (n *Napkin) Type() string {
	return "Napkin"
}

// ### Wiring everything up

// Create and connect objects, then save and load a few poems from different storage objects.
func main() {
	notebook := NewNotebook()
	napkin := NewNapkin()

	// First, write a poem into a notebook.
	// `NewPoem()` injects the dependency.
	poem := NewPoem(notebook)
	poem.Save("My first poem")

	// Create a new poem object to prove that the notebook storage works.
	poem = NewPoem(notebook)
	poem.Load("My first poem")
	fmt.Println(poem)

	// Now we do the same with a napkin as storage.
	poem = NewPoem(napkin)
	// Note the poem still just uses `Save` and `Load`. "Notebook? Napkin? I don't care."
	poem.Save("My second poem")
	poem = NewPoem(napkin)
	poem.Load("My second poem")
	fmt.Println(poem)
}

/* As usual, you can `go get` the code from GitHub. Don't forget to use -d if you do not wish to have the exectuable in your $GOPATH/bin directory.

    go get -d github.com/appliedgo/di
	cd $GOPATH/src/github.com/appliedgo/di
	./di

## Conclusion

Outside the world of poetry, dependency injection is a useful tool for decoupling logical entities, especially in multi-layered architectures as we have seen above.

Besides its benefits for layered architectures, dependency injection can also help with testing. Instead of reading a poem from a real notebook, a test can read from a notebook mockup that either is easier to set up, or delivers consistent test data, or both.


## Further reading

I definitely recommend reading the aforementioned [article about the Clean Architecture](http://blog.8thlight.com/uncle-bob/2012/08/13/the-clean-architecture.html) by [Robert C. Martin, a.k.a. "Uncle Bob"](https://de.wikipedia.org/wiki/Robert_Cecil_Martin).

The excellent article [Applying The Clean Architecture to Go applications](http://manuel.kiessling.net/2012/09/28/applying-the-clean-architecture-to-go-applications/) is a deep dive into implementing DI in Go that builds upon all four layers of the Clean Architecture. This is a great opportunity to see how entities, use cases, interfaces, and frameworks (speaking in Clean Architecture lingo) are utilized to build a (toy) shop system.

Dependency Injection can be seen as one specific form of *loose coupling*, a term referring to interconnecting components without making them too dependent on each other. Another option for loose coupling in Go (besides interfaces) is to use *higher-order functions*. I found a quick and easy intro to this topic in the blog article [Loose Coupling in Go lang](https://blog.8thlight.com/javier-saldana/2015/02/06/loose-coupling-in-go-lang.html).
*/
