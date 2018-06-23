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
title = "MapReduce - munching through Big Data"
description = "The essence of the MapReduce algorithm, explained in Go"
author = "Christoph Berger"
email = "chris@appliedgo.net"
date = "2016-12-02"
publishdate = "2016-12-02"
draft = "false"
domains = ["Big Data"]
tags = ["mapreduce", "parallelization", "data processing", "performance"]
categories = ["Tutorial"]
+++

How Google tackled the problem of processing enormous amounts of data, and how you can do the same with Go.

<!--more-->

*It's been a while since the last post, and I have to apologize for the long wait. The last weeks have been quite busy, but I finally managed to complete another article. I hope you'll enjoy it.*

## Map and Reduce

This is going to be a boring article about two boring functions, `map()` and `reduce()`. Here is the story:

You have a list with elements of type, say, `string`.

```go
var list []string
```

You define a function that takes a `string` and produces an `int`. Let's say you want to know the length of a string.

```go
func length(s string) int {
	return len(s)
}
```

Now you define a function called `map()` that takes this function and applies it to each of the elements in the list and returns a list of all results.

```go
func mäp(list []string, fn func(string)int) []int { // "map" is a reserved word, "mäp" isn't
	res := make([]int, len(list))
	for i, elem := range list {
		res[i] = fn(elem)
	}
	return res
}
```

Finally, you define another function `reduce()` that takes the result list and boils it down to a single result..

```go
func reduce(list []int, fn func(int, int)int) (res int) {
	for _, elem := range list {
		res = fn(res, elem)
	}
	return res
}

func sum(a,b int) int {
	return a+b
}
```

Now you can wire it all up.

```go
func main() {
	list := []string{"a", "bcd", "ef", "g", "hij"}
	res := reduce(mäp(list, len), sum)
	fmt.Println(res)
}
```
[(Playground link)](https://play.golang.org/p/P7-1ro4a_d)

Here is the whole thing visualized. (Click on Play to start the animation.)

HYPE[Map and Reduce](MapAndReduce.html)

That's it. End of the story. Pretty boring, eh?


## But wait! ...

... what's this?

**Looks like we just abstracted away the concept of `for` loops!**

Now if that's not something to brag about on the next Gopher meetup...

However, does this buy us anything else? Indeed it does:

* First, no more one-off index errors.

* Second, and more importantly, if the mapped function `fn` does not depend on previous results, it can be trivially called in a concurrent manner.

How to do this? Easy: Split the list into *n* pieces and pass them to *n* independently running mappers. Next, have the mappers run on separate CPU cores, or even on separate CPU's.

Imagine the speed boost you'll get. Map and reduce, as it seems, form a fundamental concept for efficient distributed loops.

> Lemme repeat that. By abstracting away the very concept of looping, you can implement looping any way you want, including implementing it in a way that scales nicely with extra hardware.
>
> Joel Spolsky, [Can Your Programming Language Do This?](http://www.joelonsoftware.com/items/2006/08/01.html) (2006)


## From map() and reduce() to MapReduce

Google researchers took the map/reduce concept and scaled it up to search engine level (I leave the exact definition of "search engine level" as an exercise for the reader). [MapReduce was born](http://research.google.com/archive/mapreduce.html).

The result was a highly scalable, fault-tolerant data processing framework with the two functions `map()` and `reduce()` at its core.

Here is how it works.

Let's say we have a couple of text files and we want to calculate the average count of nouns & verbs per file.

Our imaginary test machine has eight CPU cores. So we can set up eight processing entities/work units/actors (or whatever you want to call them):

* One input reader
* Three mappers
* One shuffler, or partitioner
* Two reducers
* One output writer


### The input reader

The input reader fetches the documents, turns each one into a list of words, and distributes the lists among the mappers.


### The mapper

Each of the mappers reads the input list word by word and counts the nouns and verbs in that list.

The result is a key-value list of word types (noun, verb) and counts. For example, our three mappers could return these counts:

	mapper 1:
		nouns: 7
		verbs: 4

	mapper 2:
		nouns: 5
		verbs: 8

	mapper 3:
		nouns: 6
		verbs: 3

When a mapper has finished, it passes the result on to the shuffler.


### The shuffler

The shuffler receives the output lists from the mappers. It rearranges the data by key; that's why it is also referred to as "partitioning function". In our example, the shuffler generates two lists, one for nouns and one for verbs:

	list 1:
		nouns: 7
		nouns: 5
		nouns: 6

	list 2:
		verbs: 4
		verbs: 8
		verbs: 3

The shuffler then passes each list to one of the two reducers.


### The reducer

Each reducer receives a list with a couple of counts. It simply runs through the list, adds up all the counts, and divides the result by the number of counts. Both reducers then send their output to the output writer.

Back to our example. The first reducer would calculate an average of

	(7 + 5 + 6) / 3 = 6

and the other one would return

	(4 + 8 + 3) / 3 = 5


### The output writer

All the output writer has to do is collecting the results from the reducers and write them to disk or pass them on to some consumer process.


### Summary

To make all this less abstract, here is the same as an animation. (Click on Play.)

HYPE[MapReduce](MapReduce.html)

This concept easily scales beyond a single multi-CPU machine. The involved entities - input reader, mapper, shuffler, reducer, and output writer - can even run on different machines if required.

But MapReduce is more than just some distributed version of `map()` and `reduce()`. There are a couple of additional bonuses that we get from a decent MapReduce implementation.

* A good deal of the functionality is the same for any kind of map/reduce task. These parts can be implemented as a MapReduce framework where the user just needs to provide the `map` and `reduce` functions.
* The MapReduce framework can provide fault recovery. If a node fails, the framework can re-execute the affected tasks on another node.
* With fault tolerance mechanisms in place, MapReduce can run on large clusters of commodity hardware.


## The code

The code below is a very simple version of the noun/verb average calculation. To keep the code short and clear, the mapper does not actually identify nouns and verbs. Instead, the input text is just a list of strings that read either "noun" or "verb". Also, the reducer does not receive key/value pairs but rather just the values. We already know that one reducer receives the nouns and the other receives the verbs.

*/
package main

import (
	"fmt"
	"sync"
)

// mapper receives a channel of strings and counts the occurrence of each unique word read from this channel. It sends the resulting map to the output channel.
func mapper(in <-chan string, out chan<- map[string]int) {
	count := map[string]int{}
	for word := range in {
		count[word] = count[word] + 1
	}
	out <- count
	close(out)
}

// reducer receives a channel of ints and adds up all ints until the channel is closed. Then it divides through the number of received ints to calculate the average.
func reducer(in <-chan int, out chan<- float32) {
	sum, count := 0, 0
	for n := range in {
		sum += n
		count++
	}
	out <- float32(sum) / float32(count)
	close(out)
}

// inputDistributor receives three output channels and sends each of them some input.
func inputReader(out [3]chan<- string) {
	// "Read" some input.
	input := [][]string{
		{"noun", "verb", "verb", "noun", "noun"},
		{"verb", "verb", "verb", "noun", "noun", "verb"},
		{"noun", "noun", "verb", "noun"},
	}

	for i := range out {
		go func(ch chan<- string, word []string) {
			for _, w := range word {
				ch <- w
			}
			close(ch)
		}(out[i], input[i])
	}
}

// shuffler gets a list of input channels containing key/value pairs like
// "noun: 5, verb: 4". For each "noun" key, it sends the corresponding value
// to out[0], and for each "verb" key to out[1].
// The input channles are multiplexed into one, based on the `merge` function
// from the [Pipelines article](https://blog.golang.org/pipelines) of the
// Go Blog.
func shuffler(in []<-chan map[string]int, out [2]chan<- int) {
	var wg sync.WaitGroup
	wg.Add(len(in))
	for _, ch := range in {
		go func(c <-chan map[string]int) {
			for m := range c {
				nc, ok := m["noun"]
				if ok {
					out[0] <- nc
				}
				vc, ok := m["verb"]
				if ok {
					out[1] <- vc
				}
			}
			wg.Done()
		}(ch)
	}
	go func() {
		wg.Wait()
		close(out[0])
		close(out[1])
	}()
}

// outputWriter starts a goroutine for each input channel and writes out
// the averages that it receives from each channel.
func outputWriter(in []<-chan float32) {
	var wg sync.WaitGroup
	wg.Add(len(in))
	// `out[0]` contains the nouns, `out[1]` the verbs.
	name := []string{"noun", "verb"}
	for i := 0; i < len(in); i++ {
		go func(n int, c <-chan float32) {
			for avg := range c {
				fmt.Printf("Average number of %ss per input text: %f\n", name[n], avg)
			}
			wg.Done()
		}(i, in[i])
	}
	wg.Wait()
}

func main() {
	// Set up all channels used for passing data between the workers.
	//
	// I could have used loops instead, to create arrays or
	// slices of channels. Apparently, copy/paste has won.
	size := 10
	text1 := make(chan string, size)
	text2 := make(chan string, size)
	text3 := make(chan string, size)
	map1 := make(chan map[string]int, size)
	map2 := make(chan map[string]int, size)
	map3 := make(chan map[string]int, size)
	reduce1 := make(chan int, size)
	reduce2 := make(chan int, size)
	avg1 := make(chan float32, size)
	avg2 := make(chan float32, size)

	// Start all workers in separate goroutines, chained together via channels.
	go inputReader([3]chan<- string{text1, text2, text3})
	go mapper(text1, map1)
	go mapper(text2, map2)
	go mapper(text3, map3)
	go shuffler([]<-chan map[string]int{map1, map2, map3}, [2]chan<- int{reduce1, reduce2})
	go reducer(reduce1, avg1)
	go reducer(reduce2, avg2)

	// The outputWriter runs in the main thread.
	outputWriter([]<-chan float32{avg1, avg2})
}

/*

The code is `go get`able from [GitHub](github.com/appliedgo/mapreduce). Ensure to use -d so that the binary does not make it into `$GOPATH/bin`.

	go get -d github.com/appliedgo/mapreduce
	cd $GOPATH/github.com/appliedgo/mapreduce
	go run mapreduce.go

This code also runs in the [Go Playground](https://play.golang.org/p/cipGuzMNT3).

Homework assignment 😉: Add logging to the code to visualize the control flow. I intentionally did not add logging statements to keep the code easy to read.


*/
