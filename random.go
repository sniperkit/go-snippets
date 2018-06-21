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
title = "A Random Blog Post"
description = "Random generators and their usage in Go"
author = "Christoph Berger"
email = "chris@appliedgo.net"
date = "2016-10-11"
draft = "false"
publishdate = "2016-10-11"
domains = ["Algorithms and Data Structures"]
tags = ["random", "mathematics", "cryptography"]
categories = ["Tutorial"]
+++

How to generate random numbers, and the difference between math/rand and crypto/rand.

<!--more-->

## Generating randomness on deterministic machines

The ideal computer is completely deterministic. (Well, except perhaps for the oldy moldy i486 Tower PC in the back of your garage running Linux 1.0 and sometimes acting kinda capricious for no particular reason.) For every input, the output is foreseeable. Trying to generate random data on a machine where everything is determined seems odd at first. But there are two ways to overcome the dichotomy between determinism and randomness.


### Exploiting sources of real randomness

Real-world computers are not quite the ideal machines that computer sciences would like them to be. There are a lot of sources in every computing device that produce more or less random data. Mouse movement, time between two keystrokes, the wall clock, activity counters (CPU load, disk activity, network activity, number of processes, etc), GPS receivers, movement sensors, and more can be used to generate continuous streams of random bits.

{{< figure src="/media/random/noise.png" class="imageLeft" alt="" >}} Simple electronic circuits can also produce true random data. For example, resistors can generate thermal noise, and Zener diodes can generate Zener breakdown noise. Transistors can be wired to produce static noise. If you turn up an amplifier that has no input device attached, you hear amplified noise from the circuits inside. Or tune an AM or FM radio between two stations, and you get atmospheric noise.

And even nature itself provides sources of random data. Photons that arrive at a semi-transparent mirror are either reflected or can pass through, in a random way. Nuclear decay creates events in a Geiger counter in random intervals. Vacuum energy fluctuates randomly.

All this noise is truly random information. An analog/digital converter can turn the noise into never-ending sequences of random bits.

Still, these "natural" sources of random data suffer from asymmetries and systematic biases caused by various physical phenomena that are inherent to the given source. Simply put, the produced bit stream may contain much more 1s than 0s on average, or vice versa. Thus the generated random numbers are not uniformly distributed.  Luckily, so-called "randomness extractors" can fix this, at the cost of a lower output rate. (I won't go into the details here, see [here](https://en.wikipedia.org/wiki/Hardware_random_number_generator) for more on this topic.)


### Generating pseudo-random numbers

The second way is to simulate a source of random data. But how, if "randomness" is not part of the concept of a deterministic machine? The trick is to produce long sequences of bits and bytes that *appear* to be random. After a while, the sequence repeats, but for many consumers of random data this is perfectly fine.

Today, a range of pseudo-random number generators (usually abbreviated as "PRNG") exist. A very simple PRNG is a bit shift register with a feedback loop.

* At each clock cycle, all bits in the register are shifted to the right.
* The rightmost bit is added to the outgoing bit stream.
* At two (arbitrary) positions of the register, the bits are extracted and fed into an Exclusive-OR (XOR) gate. The result of the XOR operation is fed back to the first bit of the register.

(To recap: The XOR operation returns "true" if both input values are different, and "false" otherwise.)

If this sounds a bit too abstract, watch the animation below. The values "true" and "false" are represented by "1" and "0", respectively.

!HYPE[Bit Shift Register](bitshift.html)

At some point in time, however, the register contains a value that it contained earlier, and at this point, the cycle repeats.

**In a generalized sense, this is how all pseudo-random number generators work. A deterministic algorithm produces a long series of seemingly random bits and bytes. Eventually the series will repeat, but depending on the algorithm, the output may still successfully pass various tests for randomness.**

There is one more point to consider: Being a deterministic algorithm, a new PRNG instance always starts at the same point of the cycle. So each time a PRNG is reset, it would deliver exactly the same sequence of numbers again. To avoid this, the algorithm can be set to start at an arbitrary value called **seed value**. This seed can be taken from a source that is known to change over time. In the simplest case, this can be the system time (`time.Now().UnixNano()` comes to mind), but sources that are more random (as described above) deliver better results.


## Go's rand packages

Go has two packages that generate random numbers: `math/rand` and `crypto/rand`. Now with the above in mind you surely already have an idea why there are two of them: One is a pseudo-random number generator, the other makes use of a source of truly random data (provided by the operating system).

But yet - why do we need both? Can't we just use `crypto/rand` for everything and enjoy truly random numbers for all purposes?

A brief look into each of the two packages may help answering this question.


### math/rand

One aspect that sets `math/rand` apart from `crypto/rand` is the rich API that includes:

* Methods that return uniformly distributed random values in different numeric formats (float33, float64, int32, int64,...).
* Methods that return `float64` values according to a non-uniform distribution - either Normal distribution and exponential distribution.
* A type named Zipf that generates Zipf-distributed values.
* And finally, a method for generating a slice of permuted (i.e., shuffled) integers.

Another one is speed. Not because `math/rand` is such a darn fast, micro-optimized package, but rather because `crypto/rand` is slow. It has to be - more about this later.

The internal pseudo-random number generator is quite simple, in the sense of "does not require complex math calculations". You can find it in [`src/math/rand/rng.go`](https://golang.org/src/math/rand/rng.go) implemented by the function `Int63()`:

```go
// from rng.go - (c) the Go team

type rngSource struct {
	tap  int         // index into vec
	feed int         // index into vec
	vec  [_LEN]int64 // current feedback register
}

func (rng *rngSource) Int63() int64 {
	rng.tap--
	if rng.tap < 0 {
		rng.tap += _LEN
	}

	rng.feed--
	if rng.feed < 0 {
		rng.feed += _LEN
	}

	x := (rng.vec[rng.feed] + rng.vec[rng.tap]) & _MASK
	rng.vec[rng.feed] = x
	return x
}
```

Remarks on the variables and constants used in this code snippet:

* `_LEN` is a constant value set to 607.
* `rng.vec` is an array of length `_LEN` that gets initialized with seed values through the `Seed()` function.
* `tap` and `feed` are initialized to 0 and 334, respectively.
* `_MASK` is a 64-bit value that has all bits except the highest one set to 1.

As you can see, the algorithm consist of four simple steps:

1. Step backwards through the array at two indexes (`tap` and `feed`). If an index reaches zero, it is set to the end of the array.
2. Add the values found at the two indexes and set the highest bit to zero (by ANDing it with `_MASK`), to ensure a positive value.
3. Save the result to the array at index `feed`.
4. Return the result.

Although it might not be easy to see at a first glance, this algorithm is a variant of the bit shift register model discussed above. See the `vec` array as a very large bit register, and `tap` and `feed` as the two positions where the values are extracted from the register, to be XOR'ed and re-inserted into the register. However, rather than shifting bits through a register (which would be fine if done in hardware but expensive if done in software), the code just cycles two indexes through the array.  Also, instead of XOR'ing tap and feed, it adds the two and adjusts the result to fit into the range of `[0..2^63)`.

This animation should make the similarities (and the differences) more apparent:

!HYPE[math/rand algorithm](mathrand.html)

The downside of `math/rand` is that the quality of the generated "randomness" is not high enough for being used in cryptographic algorithms. The data it generates might contain unforeseen repetitions or other patterns. Cryptanalysts can reveal these patterns using statistical methods. So for cryptographic purposes, we need something different, which is why `crypto/rand` exists.


### crypto/rand

`crypto/rand` does not implement an RNG algorithm; rather, it relies on the operating system to deliver cryptographically secure random numbers. On Unix-like systems, this is usually a virtual device named like [`/dev/random`](https://en.wikipedia.org/wiki//dev/random).

True sources of randomness can produce only so many bits at a time. (Side note: Crypto experts tell you the same by saying things like, "cryptographic random sources have a *limited pool of entropy*".) And the aforementioned randomness extractor reduces the throughput even more.

For this reason, Unix systems offer another device, `/dev/urandom`, that does not have this rate limitation. Usually, `/dev/urandom` is a [cryptographically secure pseudo-random number generator (CSPRNG)](https://en.wikipedia.org/wiki/Cryptographically_secure_pseudorandom_number_generator). Some Unixes even use a CSPRNG for both `/dev/urandom` and `/dev/random`. This CSPRNG must still be seeded with a truly random value though.

Even though CSPRNG's are faster than a source of true randomness, they usually are considerably slower than a typical "standard" PRNG. First, the algorithm must be more complex as it has to ensure that the generated value stream is cryptographically secure. Second, a CSPRNG instance must be guarded against simultaneous access from multiple processes, since each request for a new random value also changes the CSPRNG's internal state. [This blog post](http://blog.sgmansfield.com/2016/01/locking-in-crypto-rand/) explains the details and also compares the speed of `crypto/rand` directly against that of `math/rand`.


### Which rand package to choose

At this point, you surely already have an idea which rand package you need to use for your purpose. Nevertheless, let's summarize the criteria for picking the correct package:

* Unless you really need cryptographically secure random numbers, use `math/rand`. It will be sufficient for most of your needs. Plus, it offers a rich API (compared to `crypto/rand`) that offers different result types as well as a couple of non-uniform distributions (normal, exponential, and Zipf distribution).

* On the other hand, if the generated random value is to be used anywhere in a security context, `crypto/rand` is the only choice. Don't even think of using `math/rand` for any security-related code, only because it is faster or has more features. What you need here is cryptographically strong random number generation, period.


## Some code - just for fun

If you have time to kill, inspect the following code and try to find out how it shuffles the bytes around to generate its output. Hint: The code is an (incomplete) implementation of an algorithm is called "xoroshiro128+". [This PRNG shootout](http://xoroshiro.di.unimi.it/) includes this and a couple of other PRNG algorithms. I ported the code straight from the C implementation available on that site. (Although I must admit that there is [prior art](https://github.com/dgryski/go-xoroshiro/blob/master/xoro.go) available.)
*/

// (No explanations this time.)
package main

import (
	"fmt"
	"time"
)

var (
	s [2]uint64
)

func rotl(x uint64, k uint) uint64 {
	return (x << k) | (x >> (64 - k))
}

func next() uint64 {
	s0, s1 := s[0], s[1]
	result := s0 + s1
	s1 ^= s0
	s[0] = rotl(s0, 55) ^ s1 ^ (s1 << 14)
	s[1] = rotl(s1, 36)
	return result
}

func main() {
	s[0], s[1] = uint64(time.Now().UnixNano()^0x3bfa8764f685bd1c), uint64(time.Now().UnixNano()^0x5a2fdc2bf68cedb3) // silly seed
	for i := 0; i < 10; i++ {
		fmt.Println(next())
	}
}

/*

As usual, you can get this code via `go get github.com/appliedgo/random`, but this time you might be faster copying & pasting the code right into your editor.

(Also avaialble on the [playground](https://play.golang.org/p/bzQjF5_9g7).)

## Odds and Ends


### Third party packages for fun and profit

Below are some packages that I came across while doing research for this blog post. The list is not complete, and neither the selection nor the sort order were driven by any particular criteria other than, "hmm, that looks interesting." Here we go:

**[`random:`](https://godoc.org/github.com/DexterLB/traytor/random)** A package that extends the `math/rand` API by new result types (bool, sign, unit vector) and result ranges (between 0 and 2*pi, between a and b,...). It is part of a raytracer package.

**[`distuv:`](https://godoc.org/github.com/gonum/stat/distuv#pkg-variables)** A rather heavyweight package, featuring a large range of distribution types: Bernoulli, Beta, Categorial, Exponential, Gamma, Laplace, LogNormal, Normal, Uniform, and Weibull. It is part of the [`stat` package](https://github.com/gonum/stat) from the [`gonum` project](https://github.com/gonum).

**[`golang-petname:`](https://github.com/dustinkirkland/golang-petname)** Delivers random combinations of words to be used as a readable "ID". Similar to, for example, auto-generated container names in Docker, so that you can refer to a Docker image  as "awesome_swartz" instead of "5fe15f7e7876".

**[`go-randomdata:`](https://github.com/Pallinder/go-randomdata)** Generates random first names, last names, city names, email addresses, paragraphs, dates, and more. Good for creating mock-up data.


### An apology

Last not least, an apology is in place. My long-time readers are used to get an article every one or two weeks, but this time I failed badly delivering in time. I plan to overhaul my publishing strategy. Maybe I'll post shorter articles but then I have to struggle keeping the posts interesting for you. I am still not decided but for now rest assured that I have no intention to abandon this blog; quite the opposite is true.

So the next post *will* arrive, and until then, happy coding!

- - -

**Changelog**

2016-10-20: Section math/rand: Pointed out that `math/rand` is not cryptographically secure.
*/
