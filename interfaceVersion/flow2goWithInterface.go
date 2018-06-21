package main

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

// The interface that unites all node structs.
type processor interface {
	Process()
}

// The network is just a map from node name to the node struct (represented
// as a `processor`).
type counterNet map[string]processor

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
// from the [Go blog](https://blog.golang.org/pipelines).
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

	// PROBLEM: net["x"] is only a Processor (interface type). No way to access the
	// actual struct fields except through reflection.

	// Create the processor nodes. We need to initialize all structs here.
	// Later, any `net["abc"]` is just a Processor (an interface type) and
	// we have no more access to the structs' fields.
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
		"letterCounter": &letterCounter{
			Sentence: sToLc,
			Count:    lcToP,
		},
		"printer": &printer{
			Line1: wcToP,
			Line2: lcToP,
			Done:  done,
		},
	}

	// Start the nodes.
	fmt.Println("Start the nodes.")
	for node := range net {
		net[node].Process()
	}

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
