package main

import (
	"fmt"
)

var i int

func init() {
	fmt.Println("the init function gets started at first")
}

func main() {
	fmt.Println("the main function gets called immediately after")
	if i < 3 {
		i++
		fmt.Println("we can also call the main function ourself (but not init)")
		main()
	}
}

//or http://stackoverflow.com/questions/24790175/when-is-the-init-function-in-go-golang-run

/*
AnswerToLife() is guaranteed to run before init() is called, and init() is guaranteed to run before main() is called.

Keep in mind that init() is always called, regardless if there's main or not, so if you import a package that has an init function, it will be executed.

Also, keep in mind that you can have multiple init() functions per package, they will be executed in the order they show up in the code (after all variables are initialized of course).

A lot of the internal Go packages use init() to initialize tables and such, for example https://github.com/golang/go/blob/883bc6/src/compress/bzip2/bzip2.go#L480

*/

var WhatIsThe = AnswerToLife()

func AnswerToLife() int {
	return 42
}

func init() {
	WhatIsThe = 0
}

func main() {
	if WhatIsThe == 0 {
		fmt.Println("It's all a lie.")
	}
}
