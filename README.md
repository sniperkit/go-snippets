# qu has been superceded by the code at github.com/grolang

A command `qufmt` and associated packages to format Qu source code into Go code. The syntax of Qu code is a derivative of Go's syntax, using Kanji for keywords, special identifiers, and common package names, with the aim of eliminating the need for all whitespace in source code.

### License

Copyright © 2016 Gavin "Groovy" Grover

Distributed under the same BSD-style license as Go that can be found in the LICENSE file.

### Status

Version 0.3.

All the functionality described below is implemented.

### Installation

Run `go get github.com/gavingroovygrover/qu` to get the command and packages.

Run `go install github.com/gavingroovygrover/qu/cmd/qufmt` to compile and install the `qufmt` command from the downloaded source.

Run `go get github.com/gavingroovygrover/qutests` to get some sample qu code, i.e. most of the examples from the GoByExample website.

Run `qufmt -o src/github.com/gavingroovygrover/qutests/goByEg1.go src/github.com/gavingroovygrover/qutests/goByEg1.qu` to format one of the supplied qu code samples, which can then be run using the standard `go run src/github.com/gavingroovygrover/qutests/goByEg1.go`. Most of the examples from GoByExample have been translated into Qu.


## Intro: Spaceless Programming

By making semicolons optional at line ends, Go departs from other C-derivative languages by allowing newlines to determine semantics in some places in the code. Go doesn't actually require the code to have any newlines, though, even though its style guide and `gofmt` utility recommend their use. For the other whitespace, however, Go does require spaces (or tabs) in the syntax to determine semantics.

Qu, the Chinese word for "Go", extends the optionality of newlines to the other whitespace, enabling a Go program to be written without any whitespace. It does so by introducing one small prohibition to Go's syntax: prohibiting the use of Kanji in identifier names. We then use the Kanji as aliases for Go's keywords, special identifiers, and package names in a slightly modified syntax. Because all of the approx 80,000 Kanji available in Unicode have an implied space both before and after it in the Qu grammar, this allows Qu code to be written without any spaces or tabs, as well as no newlines.

Kanji-less Go code and Qu code can be mixed freely in the same source file. But when a program is written using Kanji in all the places it can be in the code, any of the 25 Go keywords can be used as identifier or label names, so a dedicated Qu programmer doesn't need to know any Go-specific naming exceptions to write Qu code. The Kanji, unlike other non-Ascii characters like `÷`, `≥`, or `←`, are easily enterable via the many IME's (input method editors) available for Chinese and Japanese that ship for free on OS's such as Linux and Windows, so Qu code can be typed in quickly.

Qulang is the format program that translates Qu code and Kanji-less Go code mixed together into standard Go code.


## How Qulang differs from Golang

Six basic rules completely summarize how Qulang 0.3 differs from Golang to enable spaceless programming. With those rules, it's possible to write any valid Qu program without any whitespace, e.g:

```go
包正;种A整;功正(){a:=123;形Println(整64(a));变b这A;变c这A;形Println(b,c)}//example of totally spaceless code (except this comment)
```


#### Rule 1. No Kanji in identifier names

The use of Kanji in identifier names is prohibited in the Go grammar. So `myName` is valid in both Go and Qu, but `my名` and `性名` are invalid in Qulang. This is the only way in which Go code is restricted in Qulang. Virtually no-one uses Kanji in identifiers anyway -- even Chinese and Japanese programmers really only use them inside strings and comments -- so in practise this shouldn't be a problem for anyone wanting to program in Qu.


#### Rule 2. Kanji are aliases for keywords, etc

Various single Kanji are used as aliases for Go's keywords, special identifiers, and certain package names in the Qu grammar. Each Kanji has an implied space both before and after it so the spaces needn't be written.

The 25 keywords of Go can be substituted by any of their respective Kanji below:

* `包` package, `入` import
* `变` var, `久` const, `种` type, `功` func
* `构` struct, `图` map, `面` interface, `通` chan
* `择` switch, `事` case, `别` default, `掉` fallthrough
* `如` if, `否` else, `为` for, `围` range
* `选` select, `去` go, `终` defer
* `回` return, `破` break, `继` continue, `跳` goto

So we can write Qu code using keyword aliases:

```go
包正 //包 is short for `package` with an implied space after it
     //正 is short for identifier `main` which is required here
入"fmt" //入 is short for `import`

//功 is short for `func`, and 正 for identifier `main` which is required here
功正(){
	fmt.Println("你好,世界")
}
```

The 39 special identifiers in Go can also be substituted by their associated Kanji:

* `真` true, `假` false, `空` nil, `毫` iota
* `双` bool, `节` byte, `字` rune, `串` string, `错` error, `镇` uintptr
* `整` int, `整8` int8, `整16` int16, `整32` int32, `整64` int64
* `绝` uint, `绝8` uint8, `绝16` uint16, `绝32` uint32, `绝64` uint64
* `漂32` float32, `漂64` float64, `复` complex, `复64` complex64, `复128` complex128
* `造` make, `新` new, `关` close, `删` delete, `能` cap, `度` len, `加` append, `副` copy
* `实` real, `虚` imag, `丢` panic, `抓` recover, `写` print, `线` println

The ones suffixed with a number have special support in the grammar, and are the only cases in Qu of Kanji having extra tokens associated with them.

As well as Go aliases `byte` (`节`) and `rune` (`字`), Qu adds alias `任` for `interface{}`, best verbalized as "any".

```go
package main
import("fmt")
//we use 整 for `int` and 漂64 for `float64`...
功plus(a整,b整)漂64{回漂64(a+b)} //回 for `return`
func main(){
	a:= 复(0, plus(4,6))
	fmt.Println("a is: ", a)
}
```

We also enable Kanji aliases for package names, and they aren't followed by a dot when used. Only 6 packages are implemented for now:

* `形` fmt, `网` net, `序` sort, `数` math, `大` math/big, `时` time

When a package name is aliased by a kanji, it needn't and mustn't be explicitly imported.

```go
包正;功正(){
	//形 is short for `fmt.` which is automatically imported when used
	形Println("你好,世界")
}
```

More packages will progressively be added over time.


#### Rule 3. Go and Qu code can be mixed

Code conforming to the Kanji-less Go grammar and that conforming to the Qu grammar can be mixed easily. Qu code is embedded in Go code simply by using the kanji alias of the keyword or special identifier at the head of the scope, or using special kanji `做`, best verbalized as "do", at the head of a block. Go code is embedded in Qu code by using special kanji `英`, best verbalized as "ascii", at the head of a block.

To understand the details, we must first understand the categories of identifier in Go and in Qu. Just as Go has 3 categories of identifier, i.e. global (the 25 keywords and 39 special identifiers), public (uppercase-initial identifiers), and private (all identifiers beginning with an underscore or lowercase letter, except the 25 keywords), so Qu also has 3 categories. The categories of identifier in Qu match more intuitively to their lexical class, however. They are:

* Kanji. All single-token Kanji which are aliases for keywords, special identifiers, and package names.

* public. Uppercase-initial identifiers are visible outside a package in Qu, just like in Go. They have the same format as in Go.

* protected. Identifiers that begin with an underscore followed by lowercase. They are accessible by both Qu and Go code within a single file. When defined or used within Qu code (i.e. within the static scope of a Kanji), the initial underscore is omitted.

Go's private identifiers are inaccessible within Qu code, which generally isn't a problem because they're usually used as parameters and local variables. If a private variable needs to be accessed by both Go and Qu code, put an underscore in front of it when declaring it in Go context.

The identifier `main` can't have an underscore in front of it in the generated Go code, so we provide the built-in Kanji `正` to use with `包`(package) and `功`(func).

We can see how Go code and Qu code is easily mixed:

```go
//Go by Example: Collection Functions

package正 //Go code headed by ascii `package`
入"strings" //Qu code headed by kanji 入: translated to `import _strings "strings"`
import"fmt" //Go code

功正(){ //Qu code headed by 功, so all identifiers within have underscore prefixed
  变strs=[]串{"peach","apple","pear","plum"} //`_strs` is actually generated
  形Println(Index(strs,"pear"))
  if 真{ 形Println(Include(strs,"grape")) }
  英{ //Go code embedded within Qu -- signified by 英
	形Println(Any(_strs,功(v串)双{ //referencing Qu-defined identifier from Go code, so prefix _
      回strings.HasPrefix(v,"p")
    }))

    var strs = []string{"peach", "apple", "pear", "plum"} //`strs` generated
    fmt.Println(Index(strs, "pear"))
    fmt.Println(Include(strs, "grape"))
    fmt.Println(Any(strs, func(v string) bool {
		//referencing Go-defined package identifier from Qu -- must prefix underscore...
        return _strings.HasPrefix(v, "p")
    }))
  }
}

//We can translate public functions to Qu...
功Include(vs[]串,t串)双{回Index(vs,t)>=0}

功Any(vs[]串,f功(串)双)双{
  为_,v:=围vs{如f(v){回真}}
  回假
}

//...or we can leave function defn unchanged and everything still just works
func Index(vs []string, t string) int {
    for i, v := range vs {
        if v == t {
            return i
        }
    }
    return -1
}
```

It is good style, though, to use either Go or Qu but not both as much as possible in a single file. Being able to mix both helps with gradual conversions of code from one language to the other.

The rules for Kanji acting as a header for Qu code are:

* `包`package heads everything after it in the package.
* `入`import, `变`var, `久`const, `种`type each head everything they define.
* `功`func heads everything in the receiver, parameters, results, and block after it.
* `构`struct, `图`map, `面`interface, `通`chan each head the types after them. They also head the literal data when used that way, as also do slices and arrays.
* `择`switch, `如`if, `否`else, `为`for, `选`select, `去`go, `终`defer each head everything up to the end of the following block.
* `事`case, `别`default each head the block of statements afterwards.
* `围`range, `回`return each head the expression immediately following it.
* `破`break, `继`continue, `跳`goto each head any labels that follow them.
* 36 special identifiers (all, plus `任`, except `真`true, `假`false, `空`nil, `毫`iota) head the expression they convert in a type conversion.


#### Rule 4. Keywords are identifiers within Qu

, and so all possible lowercase-initial identifiers, including Go keywords, are permitted there.

Qu code permits all possible identifiers beginning with a lowercase letter, including all 25 keywords. It does this by automatically putting an underscore at the front of all lowercase-initial identifiers. Such an identifier referenced by surrounding or embedded Go code must explicitly have the initial underscore. Special identifier `让` is available to optionally head declarations and assignments, and even required when an identifier that doubles as a Go keyword is used.

```go
包正;功正(){
	//added to demo 让:
    让range:="abc" //when used with 让, Go keywords like "range" can be used as identifiers...
    让range="abcdefg" //...and this style should be the prefered style for Qu programmers
	形Printf("range: %v\n",range)
}
```


#### Rule 5. Special identifier Kanji are protected

Go only has 25 reserved words which can't be used as variables, but all the special identifiers such as `false` can be. Kanji used as special identifiers in Qu can't be locally declared and assigned to as they can in Go:

```go
package main
func main() {
	a:= true
	b:= 真 //Kanji for `true` used on right-hand side
	nil:= true //Qu still allows special identifiers (here, `nil`) to be used on the left-hand side, ...
	iota:= 真
	//假:= true // ... but when the Kanji version is used, e.g. 假 for false,
	           // generates a parse error "expected non-kanji special identifier on left hand side"
	形Printf("a: %v, b: %v, nil: %v, iota: %v\n", a, b, nil, iota)

	abc:=图[双]整{}
	abc[真]=789 //of course, kanji can still be on the LHS when not being assigned to
}
```


#### Rule 6. Packages can be aliased with Kanji

A package not in the registry of Kanji-aliased packages can be given a temporary Kanji when imported. The Kanji is any of those with the `口` radical on the left hand side. There's about 2000 such Kanji to choose from out of the 80,000 in Unicode. Types defined in the current package can be prefixed with `这`, best verbalized as "this", if desired.

```go
package main

入"fmt"
import 吧"fmt" //we can use any Kanji with 口-radical on LHS...
import 哪_fg"fmt" //...with imports that don't have their own dedicated Kanji
入㕤hij"fmt"
入卟"unicode/utf8"
入吗嗎kl"unicode/utf8" //can even put in two aliases

type A int
变n = 50
功正(){
  var b A
  var c这A //can use `这` with locally-defined type to achieve spaceless program
  形Println(b, c)

  如假{
    fg.Printf("Len: %d\n", 度("hijk") + n)
  }否{
    hij.Printf("Len: %d\n", 度("hi") + n)
  }
  fr,_:= utf8.DecodeRune([]节("lmnop"))
  fmt.Printf("1st rune: %s; Len: %d\n", 串(fr), 度("lmnop") + n)
  让_,_=吗DecodeRune([]节("lmnop"))
  㕤Printf("Fifty: %d\n", n)
  哪Printf("Fifty: %d\n", n)
  吧Printf("Fifty: %d\n", n)
}
```


## Examples

Most of the Go examples from gobyexample.com (or its Chinese translation gobyexample.everyx.in ) are available as Qu source in github.com/gavingroovygrover/qutests , all except the last few. (Go By Example is copyrighted by Mark McGranaghan.)

For example, the code from the first 3 pages of "Go by Example" can be replaced by the terser:

```go
包正 //包 is shorthand Kanji for `package`; 正 is short for `main`
久s串="constant" //久 short for `const`
功正(){ //功 short for `func`

  //Go by Example: Hello world
  //形 is short for `fmt.` which is automatically imported when used
  形Println("你好,世界")

  //Go by Example: Values
  形Println("go"+"lang") //Canonical "spaceless" Qu style: no spaces within statements and expressions
  形Println("1+1 =",1+1)
  形Println("7.0/3.0 =",7.0/3.0)
  形Println(真&&假) //真 for `true`; 假 for `false`
  形Println(真||假)
  形Println(!真)

  //Go by Example: Variables
  做{ //optionally, begin standalone block with `做`
    变a串= "initial";形Println(a) //变 for `variable`; 串 for `string`
    变b,c整=1,2;形Println(b, c) //整 for `int`
    变d=真;形Println(d)
    变e整;形Println(e)
    f:="short";形Println(f)
    //canonical Qu style: join stmts on one line with semicolons if they're written
    //as separate lines between blank lines in corresponding Go style
  }
}
功plus(a整,b整)整{回a+b} //回 for `return`
功plusPlus(a,b,c整)整{回a+b+c}
```

If we shorten the local identifier names, and use semicolon (`;`) to join lines together, we can achieve much greater tersity.


## Rationale

Qu has the purpose of generating discussion on which Kanji should map to which keyword, special identifier, and package in Go, in both China and Japan. If one standard repertoire of Kanji is eventually adopted under the control of some other party/s, then the author encourages others to clone, modify, and promote their own editions and clones of Qu. My choice of Kanji are only interim because native Chinese and Japanese speakers will make the final choices.

The name "qu" is the pinyin of the Mandarin translation of "to go". The Qu syntax is a modification of Go's, parsed by a modified edition of the new recursive descent parser shipped in Go 1.6. The `parser`, `scanner`, and `cmd/gofmt` packages were copied and modified, and the `golang.org/x/tools/go/ast/astutil` package copied.

