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
title = "Regular Expressions demystified"
description = "A short video tutorial on Regular Expression basics."
author = "Christoph Berger"
email = "chris@appliedgo.net"
date = "2016-08-29"
publishdate = "2016-08-29"
domains = ["Text Processing"]
tags = ["regular expression", "regexp", "search", "pattern", "video"]
categories = ["Tutorial"]
+++

Regular Expressions are slow, ugly, error-prone, incomprehensible,... Or are they? Find out by learning regexp basics.

<!--more-->

Regular Expressions seem to divide software developers. Some love them and use them without thinking twice, some frown upon any regexp they spot in someone else's code. To others, regular expressions are all Greek. Who is right? Is the truth somewhere in the middle, as is so often the case when people take extreme standpoints on a topic?

I'd say the best way is to find out by yourself. For this I made a short video about the basic building blocks of regular expressions. Here we go:

{{< youtube nDiFzmGFjS8 >}}

The video does not cover Go regexp methods (or only very briefly near the end), so let's examine some useful methods from the `regexp` library here.

*/

// ## Importing the regexp package
package main

// The `regexp` package is the only one we need. `regexp/syntax` contains some low-level functions that usually are not used directly. `regexp` uses these methods internally.
import (
	"fmt"
	"regexp"
)

// prettyMatches formats Matches nicely.
func prettyMatches(m []string) string {
	s := "["
	for i, e := range m {
		s += e
		if i < len(m)-1 {
			s += "|"
		}
	}
	s += "]"
	return s
}

// PrettySubmatches formats Submatches nicely.
func prettySubmatches(m [][]string) string {
	s := "[\n"
	for _, e := range m {
		s += "    " + prettyMatches(e) + "\n"
	}
	s += "]"
	return s
}

// Let's define some regular expressions and text to search in. Note the backticks used for the regexp that contains a backslash. If we used double quotes, we would need to double the backslash to avoid an "unknown escape sequence" error.
var (
	exps = []string{"b.*tter", "b(i|u)tter", `batter (\w+)`}

	text = `Betty Botter bought some butter 
But she said the butter’s bitter 
If I put it in my batter, it will make my batter bitter 
But a bit of better butter will make my batter better 
So ‘twas better Betty Botter bought a bit of better butter`
)

// PrintSlice prints a slice in a more readable way. Standard Println or Printf separate the elements by a space, but our text also contains spaces, so we need something
// Now try some of the various Find functions.
func main() {
	for _, e := range exps {
		re := regexp.MustCompile(e)
		fmt.Println(e + ":")
		fmt.Println("1. FindString: ", re.FindString(text))
		fmt.Println("2. FindStringIndex: ", re.FindStringIndex(text))
		fmt.Println("3. FindStringSubmatch: ", re.FindStringSubmatch(text))
		fmt.Printf("4. FindAllString: %v\n", prettyMatches(re.FindAllString(text, -1)))
		fmt.Printf("5. FindAllStringIndex: %v\n", re.FindAllStringIndex(text, -1))
		fmt.Printf("6. FindAllStringSubmatch: %v\n\n", prettySubmatches(re.FindAllStringSubmatch(text, -1)))
	}
}

/* ## Closing remarks

I hope you enjoyed the video. As always, the code is available on GitHub:

```
go get -d github.com/appliedgo/regexp
cd $GOPATH/src/github.com/appliedgo/regexp
go run regexp.go
```

Also available [on the Go Playground](https://play.golang.org/p/yNusC1jDQ3).

Feel free to experiment with the expressions and see if the outcome is what you expected!


### Links from the video

[Go Regexp Syntax Reference](https://golang.org/pkg/regexp/syntax/)

[RegexPlanet](http://www.regexplanet.com/advanced/golang/index.html)

Happy coding!
*/
