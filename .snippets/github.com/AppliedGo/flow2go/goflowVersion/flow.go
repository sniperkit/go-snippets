package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/trustmaster/goflow"
)

/*
First, we define the nodes. Each node is a struct with an embedded `flow.Component` and input and output channels (at least one of each kind, except for a sink node that only has input channels).

Nodes can act on input by functions that are named after the input channels. For example, if an input channel is named "Strings", the function that triggers on new input is called "OnStrings" by convention.

We define these nodes:

* A splitter that takes the input and copies it to two outputs.
* A word counter that counts the words (i.e., non-whitespace content) of a sentence.
* A letter counter that counts the letters (a-z and A-Z) of a sentence.
* A printer that prints its input.

None of these nodes knows about any of the other nodes, and does not need to.
*/

// Our two `counter` nodes (see below) send their results asynchronously to the `printer` node. To distinguish between the outputs of the two counters, we attach a tag to each count. (Yes, sending just a string including the count would be easier but also more boring. The `splitter` already sends strings, so let's try something different here.)
type count struct {
	tag   string
	count int
}

// The `splitter` receives strings and copies each one to its two output ports.
type splitter struct {
	flow.Component

	In         <-chan string
	Out1, Out2 chan<- string
}

// `OnIn` dispatches the input string to the two output ports.
func (t *splitter) OnIn(s string) {
	t.Out1 <- s
	t.Out2 <- s
}

// `WordCounter` is a `goflow` component that counts the words in a string.
type wordCounter struct {
	// Embed flow functionality.
	flow.Component
	// The input port receives strings that (should) contain words.
	Sentence <-chan string
	// The output port sends the word count as integers.
	Count chan<- *count
}

// `OnSentence` triggers on new input from the `Sentence` port.
// It counts the number of words in the sentence.
func (wc *wordCounter) OnSentence(sentence string) {
	wc.Count <- &count{"Words", len(strings.Split(sentence, " "))}
}

// `letterCounter` is a `goflow` component that counts the letters in a string.
type letterCounter struct {
	flow.Component
	Sentence <-chan string
	// The output port sends the letter count as integers.
	Count chan<- *count
	// To identify letters, we use a simple regular expression.
	re *regexp.Regexp
}

// `OnSentence` triggers on new input from the `Sentence` port.
// It counts the number of words in the sentence.
func (lc *letterCounter) OnSentence(sentence string) {
	lc.Count <- &count{"Letters", len(lc.re.FindAllString(sentence, -1))}
}

// An `Init` method allows to initialize a component. Here we use it to run
// the expensive `MustCompile` method once, rather than every time `OnSentence` is called.
func (lc *letterCounter) Init() {
	lc.re = regexp.MustCompile("[a-zA-Z]")
}

// A `printer` is a "sink" with no output channel. It prints the input
// to the console.
type printer struct {
	flow.Component
	Line <-chan *count // inport
}

// `OnLine` prints a count.
func (p *printer) OnLine(c *count) {
	fmt.Println(c.tag+":", c.count)
}

// `CounterNet` represents the complete network of nodes and data pipelines.
type counterNet struct {
	flow.Graph
}

/*

### Assembling the network

With the nodes in place, we can go foward and create the complete network, adding and connecting all the nodes.

*/

// Construct the network graph.
func NewCounterNet() *counterNet {
	n := &counterNet{}
	// Initialize the net.
	n.InitGraphState()
	// Add nodes to the net. (I derived from the documentation by using `&{}`
	// instead of `new`.) Each node gets a name assigned that is used later
	// when connecting the nodes.
	n.Add(&splitter{}, "splitter")
	n.Add(&wordCounter{}, "wordCounter")
	n.Add(&letterCounter{}, "letterCounter")
	n.Add(&printer{}, "printer")
	// Connect the nodes. The parameters are: Sending node, sending port,
	// receiving node, and receiving port.
	n.Connect("splitter", "Out1", "wordCounter", "Sentence")
	n.Connect("splitter", "Out2", "letterCounter", "Sentence")
	n.Connect("wordCounter", "Count", "printer", "Line")
	n.Connect("letterCounter", "Count", "printer", "Line")
	// Our net has 1 input port mapped to `splitter.In`.
	n.MapInPort("In", "splitter", "In")
	return n
}

/*

### Launching the network

Finally, we only need to activate the network, create an input port, and start feeding it with selected bits of wisdom.

*/

//
func main() {
	// Create the network.
	net := NewCounterNet()
	// We create a channel as the input port of the network.
	in := make(chan string)
	net.SetInPort("In", in)
	// Start the net.
	flow.RunNet(net)
	// Now we can send some text and see what happens. This is as easy as sending
	// text to the input channel. (All aphorisms by Oscar Wilde.)
	in <- "I never put off till tomorrow what I can do the day after."
	in <- "Fashion is a form of ugliness so intolerable that we have to alter it every six months."
	in <- "Life is too important to be taken seriously."
	// Closing the input channel shuts the network down.
	close(in)
	// Wait until the network has shut down.
	<-net.Wait()
}
