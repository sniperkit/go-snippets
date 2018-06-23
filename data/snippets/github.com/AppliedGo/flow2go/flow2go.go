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
title = "Flow To Go"
description = "Flow Based Programming in pure Go"
author = "Christoph Berger"
email = "chris@appliedgo.net"
date = "2017-03-11"
publishdate = "2017-03-11"
draft = "false"
domains = ["Concurrent Programming"]
tags = ["FBP", "flow-based programming", "dataflow"]
categories = ["Tutorial"]
+++

If you want to do Flow-Based Programming in Go, there are a couple of frameworks and libraries available. Or you simply do it with pure Go. After all, how difficult can it be?

<!--more-->

The [previous article](https://appliedgo.net/flow/) peeked into Flow-Based Programming (FBP), a paradigm that puts the flow of the data above the code that makes the data flow. An FBP application can be described as a network of loosely coupled processing nodes, only connected through data pipelines. The article's code made use of a quite sophisticated FBP library that made the magic of a convenient syntax happen through reflection (hidden within the library, but still).

The article triggered a couple of comments on [Reddit](https://www.reddit.com/r/golang/comments/5wg6jm/get_into_the_flow_flowbased_programming_applied_go/) that were suggesting pure Go approaches, without using third-party libraries.

This made me curious.

*How well would the code from the previous article translate into "just stdlib" code?*

I gave it a try, and here is the result.


## Constructing an FBP net in pure stdlib Go

The approach is as simple as it can get:

**Data flow:** Every node reads from one or more input channels, and writes to one or more output channels.

**Network construction:** Weaving the net happens in `main()` through creating channels and connecting them to the input and output ports of the processing nodes.

**Starting the net:** The net starts by calling a `Process()` method on each node and feeding data into the network's input channel.

**Stopping the net:** The net stops when the net's input channel is closed. Then every node whose input channels get closed closes his output channels and shuts down, and this way the shutdown propagates through the network until the last node (the "sink" node with no output channel) stops.


## Changes to the code

The code below is based on a 1:1 copy of the code from the previous article. Then the following changes were applied.

### Input channels

The `goflow` framework takes care of each node's input channels, and the nodes need special "`On...()`" functions that received a single channel item at a time.

I changed the nodes to have their own input channels, and I replaced the `On...()` methods with `Process()` methods that take no arguments and start a goroutine to read from the input channel(s) and write to the output channel(s). This is substantially more code compared to the `On...()` methods that mostly were one-liners; however, in real life where each node would contain much more code, the overhead for input and output handling would be negligible.


### No more fan-in

The original code used one channel between the two counter nodes and the printer node. Go channels trivially support a fan-in scenario with multiple writers and one reader on the same channel.

I had to change this so that the printer node now has two input channels, and the two counter nodes do not share the same output channel anymore but send their results into separate channels.

Why? The reason is the network's shutdown mechanism. As described above, each node shuts down when its input channels are closed. Piece of cake, you might think, but things get difficult when a channel has multiple writers, as in the counter/printer part of our network.

As you know, closing a channel closes it completely, and other writers panic when trying to write to this channel. (Personally, I would prefer a fan-in semantic where all writers except the last one would only close their own end of the channel they share rather than the whole channel at once, but this is not how channels work in Go.)

So we need to split every multi-writer channel into separate channels. Then we can write a `merge` function that merges all the channels into one, and also takes care of closing the output channels when the last of the input channels closes.

Or, rather than writing one, we can take a ready-made `merge()` function [from the Go blog](https://blog.golang.org/pipelines#TOC_4.) With some very minor changes, the `merge` function is now a method of the `printer` node. Problem solved!


### Signaling shutdown completion to the outside

Without the `goflow` framework, we also need to add a mechanism to tell the outside that the network has shut down. This is the duty of the final node in the network. Similar to how `goflow` does it, our `printer` node closes a channel of empty structs when concluding work.

 An empty, unbuffered channel blocks its readers. When it is closed, however, it starts delivering the channels zero value. Any read operation on this channel then unblocks, and this is how we can make `main()` wait for the network to shut down.

(Side note: This behavior may seem counterintuitive and difficult to deal with, but remember that the ["comma, ok" idiom for the receive operator](https://golang.org/ref/spec#Receive_operator) tells you if the channel has been closed.)


## The code

Ok, enough talking about the whats and whys, now let's dive into the code!
*/

// Imports
package main

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

// Every type or func below this point was taken over from the previous article's code. Wherever I had to make a change, a comment explains what and why.
//
// In the GitHub repository for this article you can find the file `flow.go` that contains the original code from the previous article, for easy comparison. (As always, find the `go get` instructions at the end of the article.)
type count struct {
	tag   string
	count int
}

type splitter struct {
	In         <-chan string
	Out1, Out2 chan<- string
}

// In the old code, this method was called `OnIn()` and was fed with a string via the `goflow` framework. This new method now reads the input channel directly within a goroutine. When the channel is closed and drained, the goroutine closes its output channels and exits.
func (t *splitter) Process() {
	fmt.Println("Splitter starts.")
	go func() {
		for {
			s, ok := <-t.In
			if !ok {
				fmt.Println("Splitter has finished.")
				close(t.Out1)
				close(t.Out2)
				return
			}
			t.Out1 <- s
			t.Out2 <- s
		}
	}()
}

type wordCounter struct {
	Sentence <-chan string
	Count    chan<- *count
}

// Previously, this function was named OnSentence, and it was a one-liner that received a string via the `goflow` framework. As with `splitter`'s `OnIn()`, let's replace it by a function that reads the input channel directly.
func (wc *wordCounter) Process() {
	fmt.Println("WordCounter starts.")
	go func() {
		for {
			sentence, ok := <-wc.Sentence
			if !ok {
				fmt.Println("WordCounter has finished.")
				close(wc.Count)
				return
			}
			wc.Count <- &count{"Words", len(strings.Split(sentence, " "))}
		}
	}()
}

type letterCounter struct {
	Sentence <-chan string
	Count    chan<- *count
	re       *regexp.Regexp
}

// As with `wordCounter`,  `letterCounter`'s `OnSentence` function also got replaced by a function that reads the input channel directly.
func (lc *letterCounter) Process() {
	fmt.Println("LetterCounter starts.")
	go func() {
		lc.Init()
		for {
			sentence, ok := <-lc.Sentence
			if !ok {
				fmt.Println("LetterCounter has finished.")
				close(lc.Count)
				return
			}
			lc.Count <- &count{"Letters", len(lc.re.FindAllString(sentence, -1))}
		}
	}()
}

func (lc *letterCounter) Init() {
	lc.re = regexp.MustCompile("[a-zA-Z]")
}

// printer now has two input channels instead of one, so that each sender can simply close its channel when the data flow ends.
type printer struct {
	Line1 <-chan *count
	Line2 <-chan *count
	// Close this channel when all input channels are closed; the network has stopped then.
	Done chan<- struct{}
}

// As we now have two input channels, we need to merge them into one.
// For this we use a slightly modified version of the `merge` function
// from the [Go blog](https://blog.golang.org/pipelines#TOC_4.).
func (p *printer) merge() <-chan *count {
	var wg sync.WaitGroup
	out := make(chan *count)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan *count) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(2)
	go output(p.Line1)
	go output(p.Line2)

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// `printer`'s `OnLine()` method was also replaced by a `Process()` method that reads directly from the input channel.
func (p *printer) Process() {
	fmt.Println("Printer starts.")
	in := p.merge()
	go func() {
		for {
			c, ok := <-in
			if !ok {
				fmt.Println("Printer has finished.")
				close(p.Done)
				return
			}
			fmt.Println(c.tag+":", c.count)
		}
	}()
}

// Now let's build the flow network with pure Go only.
func main() {
	// Create the processor nodes.
	s := &splitter{}
	wc := &wordCounter{}
	lc := &letterCounter{}
	p := &printer{}

	// Create the channels for the network.
	// We do not want to synchronize the nodes, so we use buffered
	// channels. The channel capacity was chosen arbitrarily.
	in := make(chan string, 10)
	sToWc := make(chan string, 10)
	sToLc := make(chan string, 10)
	wcToP := make(chan *count, 10)
	lcToP := make(chan *count, 10)

	// The `done` channel is used by the last node (the "sink") to signal that
	// the network has stopped.
	done := make(chan struct{})

	// Connect the nodes to each other.
	s.In = in
	s.Out1 = sToWc
	s.Out2 = sToLc

	wc.Sentence = sToWc
	wc.Count = wcToP

	lc.Sentence = sToLc
	lc.Count = lcToP

	p.Line1 = wcToP
	p.Line2 = lcToP
	p.Done = done

	// Start the nodes.
	fmt.Println("Start the nodes.")
	s.Process()
	wc.Process()
	lc.Process()
	p.Process()

	// Now feed the network with data.

	fmt.Println("Send the data into the network.")
	in <- "I never put off till tomorrow what I can do the day after."
	in <- "Fashion is a form of ugliness so intolerable that we have to alter it every six months."
	in <- "Life is too important to be taken seriously."
	// Closing the input channel shuts the network down.
	close(in)
	// Wait until the network has shut down.
	<-done
	fmt.Println("Network has shut down.")
}

/*

And that's it! For more complex networks, you can always define an interface like this...

```go
type processor interface {
	Process()
}
```

...then define a network...

```go
type counterNet map[string]processor
```

...and in `main()`, create the network from the nodes...

```go
net := counterNet{
	"splitter": &splitter{
		In:   in,
		Out1: sToWc,
		Out2: sToLc,
	},
	"wordCounter": &wordCounter{
		Sentence: sToWc,
		Count:    wcToP,
	},
	// ...
}
```

...and then start all nodes within a loop (thanks to the interface defined above):

```go
for node := range net {
	net[node].Process()
}
```

When you `go get` the code (see below), an extra file with a runnable interface version of the code is included (`interfaceVersion/flow2goWithInterface.go`).

## How to get and run the code

Step 1: `go get` the code. Two notes here:

* Use the `-d` flag to prevent auto-installing the binary into `$GOPATH/bin`.
* Ensure to include the ellipsis (...) at the end to also fetch the alternate versions in the subdirectories `goflowVersion` (from the previous article) and `interfaceVersion` (see above).

```
go get -d github.com/appliedgo/flow2go/...
```

Step 2: `cd` to the source code directory.

    cd $GOPATH/src/github.com/appliedgo/flow2go

Step 3. Run the binary.

    go run flow2go.go

You should see an output similar to this:

	Start the nodes.
	Splitter starts.
	WordCounter starts.
	LetterCounter starts.
	Printer starts.
	Send the data into the network.
	Splitter has finished.
	WordCounter has finished.
	Words: 13
	Words: 17
	Words: 8
	Letters: 45
	LetterCounter has finished.
	Letters: 70
	Letters: 36
	Printer has finished.
	Network has shut down.

## Conclusion

With only some basic Go mechanisms - goroutines, channels, and a WaitGroup (in the `merge` method), it is possible to re-implement the FBP network from the previous article without any third-party library. The code size increased a bit but in a manageable way that should scale quite well with the number of nodes.


**Happy coding!**


*/
