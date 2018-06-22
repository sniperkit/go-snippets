# Welcome

I have written some code snippets while learning golang. Look around and learn!

## main function in golang

Every Go program starts in a package main by calling identically named function. When this function returns the program ends its execution.

```go
  package main
  import "fmt"

  func main() {
      fmt.Printf("Hello from the GO programming language\n")
      fmt.Printf("This is lot like C programming language\n")
  }
```

### golang commands on bash

```bash
go run hello-world.go
```

```bash
go build hello-world.go
./hello-world
```

### project structure for golang code

refer to the link below to read on how to write your first library in golang

[project structure] (https://golang.org/doc/code.html)

### GOPATH setup

```bash
# Edit ~/.bash_profile

export GOPATH=/Users/navkar/go
export PATH=$GOPATH/bin:$PATH

# Reload profile : source ~/.bash_profile
```

### init function in golang and more on packages

* init functions are defined in package block.
* init function doesn’t take arguments neither returns any value.
* identifier init is not declared so cannot be referenced.
* Many init functions can be defined in the same package.
* init identifier doesn’t introduce binding, just as blank identifier represented by underscore character.
* Package initialization is done only once even if package is imported many times.
* To use a imported package it needs to be initialized first.
* Package in Go can contain many files. To ensure reproducible initialization behavior, build systems are encouraged to present multiple files belonging to the same package in lexical file name order to a compiler.
* init is called after all the global variables in the package (all files) are initialized.
* The most common use case of init function is to assign a value which cannot be calculated as a part of initialization expression
* for instance, it’s not possible to use for loop as an expression (it’s a statement in Go) so putting it into init function solves the problem.

```go
  var precomputed = [20]float64{}
  func init() {
      var current float64 = 1
      precomputed[0] = current
      for i := 1; i < len(precomputed); i++ {
          precomputed[i] = precomputed[i-1] * 1.2
      }
  }
```

### solving unused imports with a blank identifier

* Go is very strict when it comes to unused variables/imports.

```go
  import _ "ThisIsGoingToBuild"
```

### what are blocks in golang

* Block is a sequence of statements (empty sequence is also valid). Blocks can be nested and are denoted by curly brackets

```go
  package main
  import “fmt”
  func main() {
      { // start outer block
          a := 1
          fmt.Println(a)
          { // start inner block
              b := 2
              fmt.Println(b)
          } // end inner block
      } // end outer block
  }
```

* universe block contains all source code
* package block contains all package’s source code (package can be spread over several files in a single directory)
* file block contains file’s source code (importing packages names have scope set to file block)
* for statement is in its own implicit block

```go
  for i := 0; i < 5; i++ {
      fmt.Println(i)
  }
```

```go
  if i := 0; i >= 0 {
      fmt.Println(i)
  }
```

```go
  switch i := 2; i * 4 {
  case 8:
      fmt.Println(i)
  default:
      fmt.Println(“default”)
  }
```

```go
  v := 1
  {
      v := 2  // short variable declaration
      fmt.Println(v)
  }
  fmt.Println(v)
```

### Find the difference between the following code snippets

#### snippet #1

```go
  package main
  import "fmt"
  var (
      a int = b + 1
      b int = 1
  )
  func main() {
      fmt.Println(a)
      fmt.Println(b)
  }
```

#### snippet #2

```go
  package main
  import "fmt"
  func main() {
      var (
          a int = b + 1
          b int = 1
      )
      fmt.Println(a)
      fmt.Println(b)
  }
```

### Golang method receivers

* you define a method receiver to specify which struct to attach a certain function to in order to make it invoke-able as a method.

For instance, `func (d Dog)` is part which defines the method receiver in the following program:

```go
  package main

  import "fmt"

  type Dog struct {
  }

  func (d Dog) Say() {
      fmt.Println("Woof!")
  }

  func main() {
      d := &Dog{}
      d.Say()
  }
```

#### what’s the difference between pointer and non-pointer method receivers

You can treat the receiver as if it was an argument being passed to the method. All the same reasons why you might want to pass by value or pass by reference apply.

```go
  func (s *MyStruct) pointerMethod() { } // method on pointer
  func (s MyStruct)  valueMethod()   { } // method on value
```

First, and most important, does the method need to modify the receiver? If it does, the receiver must be a pointer. (Slices and maps act as references, so their story is a little more subtle, but for instance to change the length of a slice in a method the receiver must still be a pointer.) In the examples above, if pointerMethod modifies the fields of s, the caller will see those changes, but valueMethod is called with a copy of the caller's argument (that's the definition of passing a value), so changes it makes will be invisible to the caller.

Reasons why you would want to pass by reference as opposed to by value:

* You want to actually modify the receiver (“read/write” as opposed to just “read”)
* The struct is very large and a deep copy is expensive
* Consistency: if some of the methods on the struct have pointer receivers, the rest should too. This allows predictability of behavior

```go
  package main

  import "fmt"

  type Mutable struct {
      a int
      b int
  }

  func (m Mutable) StayTheSame() {
      m.a = 5
      m.b = 7
  }

  func (m *Mutable) Mutate() {
      m.a = 5
      m.b = 7
  }

  func main() {
    // With method receivers that take pointers
    // Go conveniently allows both pointers and non-pointers to be passed and it automatically does the conversion
      m := &Mutable{0, 0}
    //  m := new (Mutable) - this will work as well
      fmt.Println(m)
      m.StayTheSame()
      fmt.Println(m)
      m.Mutate()
      fmt.Println(m)
  }
```

### Struct in golang - Employee-Manager example

```go
package main
import (
  "fmt"
)

type Employee struct {
  FirstName string
  LastName string
  Role string
  Manager *Employee
}

func (e *Employee) SetManager(mgr *Employee) {
  e.Manager = mgr
}

//the type *Employee is the receiver of the ManagerRole method.
func (e *Employee) ManagerRole() {
  e.Role = "manager"
}

func (e *Employee) EmployeeRole() {
  e.Role = "employee"
}

// factory function to create a new instance
func NewEmployee(firstName string, lastName string) *Employee {
  return &Employee {
    FirstName: firstName,
    LastName: lastName,
    Role: "",
    Manager: nil }

    // emp := new (Employee)
    // emp.FirstName = firstName
    // emp.LastName = lastName
    // return emp
}

//emp := new (Employee) // same as
//emp := &Employee{}

func main() {
  mgr := NewEmployee ("Bill", "Gates")
  mgr.ManagerRole()

  emp := new (Employee)
  emp.FirstName = "Satya"
  emp.LastName = "Nadella"
  emp.EmployeeRole()
  emp.SetManager(mgr)

  fmt.Println("FirstName | LastName | Role | Manager")
  fmt.Println(mgr)
  fmt.Println(emp)
  fmt.Println(emp.Manager.FirstName + " is the manager of " + emp.FirstName)
}
```

#### Output

```bash
Naveens-MacBook-Pro:structs navkar$ go run employee.go
FirstName | LastName | Role | Manager
&{Bill Gates manager <nil>}
&{Satya Nadella employee 0xc420014200}
Bill is the manager of Satya
Naveens-MacBook-Pro:structs navkar$
```

### golang interfaces

* an interface is a set of methods that we use to define a set of actions
* an interface defines a set of methods, so if a type implements all the methods we say that it implements the interface.
* an interface can be implemented by any type, and one type can implement many interfaces simultaneously.
* any type implements the empty interface `interface{}` because it doesn't have any methods and all types have zero methods by default.

```go
  // define a as empty interface
  var a interface{}
  var i int = 5
  s := "Hello world"
  // a can store value of any type
  a = i
  a = s
```

* an interface is a set of abstract methods, and can be implemented by non-interface types.
* an empty interface is an interface that doesn't contain any methods, so all types implement an empty interface. This fact is very useful when we want to store all types at some point, and is similar to `void*` in C.
* If a function uses an empty interface as its argument type, it can accept any type; if a function uses empty interface as its return value type, it can return any type.

### type of a variable in an interface

If a variable is the type that implements an interface, we know that any other type that implements the same interface can be assigned to this variable. The question is how can we know the specific type stored in the interface.

#### Assertion of Comma-ok pattern

Go has the syntax `value, ok := element.(T)`. This checks to see if the variable is the type that we expect, where "value" is the value of the variable, "ok" is a variable of boolean type, "element" is the interface variable and the T is the type of assertion.
If the element is the type that we expect, ok will be true, false otherwise.

```go
  package main

  import (
      "fmt"
      "strconv"
  )

  type Element interface{}
  type List []Element

  type Person struct {
      name string
      age  int
  }

  func (p Person) String() string {
      return "(name: " + p.name + " - age:     " + strconv.Itoa(p.age) + " years)"
  }

  func main() {
      list := make(List, 3)
      list[0] = 1       // an int
      list[1] = "Hello" // a string
      list[2] = Person{"Dennis", 70}

      for index, element := range list {
          if value, ok := element.(int); ok {
              fmt.Printf("list[%d] is an int and its value is %d\n", index, value)
          } else if value, ok := element.(string); ok {
              fmt.Printf("list[%d] is a string and its value is %s\n", index, value)
          } else if value, ok := element.(Person); ok {
              fmt.Printf("list[%d] is a Person and its value is %s\n", index, value)
          } else {
              fmt.Printf("list[%d] is of a different type\n", index)
          }
      }
  }
```

#### Assertion of switch pattern

* `element.(type)` cannot be used outside of the `switch` body, which means in that case you have to use the comma-ok pattern .

```go
  func main() {
      list := make(List, 3)
      list[0] = 1       //an int
      list[1] = "Hello" //a string
      list[2] = Person{"Dennis", 70}

      for index, element := range list {
          switch value := element.(type) {
          case int:
              fmt.Printf("list[%d] is an int and its value is %d\n", index, value)
          case string:
              fmt.Printf("list[%d] is a string and its value is %s\n", index, value)
          case Person:
              fmt.Printf("list[%d] is a Person and its value is %s\n", index, value)
          default:
              fmt.Println("list[%d] is of a different type", index)
          }
      }
  }
```

### Embedded Interfaces in golang

* If an interface has another interface embedded within it, it will behave as if it has all the methods that the embedded interface has.

We can see that the source file in `container/heap` has the following definition:

```go
  type Interface interface {
    sort.Interface // embedded sort.Interface
    Push(x interface{}) //a Push method to push elements into the heap
    Pop() interface{} //a Pop method that pops elements from the heap
  }
```

We see that sort.Interface is an embedded interface, so the above Interface has the three methods contained within the `sort.Interface` implicitly.

```go
type Interface interface {
  // Len is the number of elements in the collection.
  Len() int
  // Less returns whether the element with index i should sort
  // before the element with index j.
  Less(i, j int) bool
  // Swap swaps the elements with indexes i and j.
  Swap(i, j int)
}
```

### The interface{} type

The interface{} type is the interface that has no methods. Since there is no implements keyword, all types implement at least zero methods, and satisfying an interface is done automatically, all types satisfy the empty interface. That means that if you write a function that takes an interface{} value as a parameter, you can supply that function with any value.

```go
func PrintAnyInterface(any interface{}) {

// REF : https://research.swtch.com/interfaces
// if v, ok := any.(Stringer); ok {
//     return v.String()
// }

    switch v := any.(type) {
    case IDisposable:
        fmt.Println(v.dispose())
        return
    default:
        fmt.Println(v)
    }
}
```

An interface value is constructed of two words of data; one word is used to point to a method table for the value’s underlying type, and the other word is used to point to the actual data being held by that value. Understand that an interface value is two words wide and it contains a pointer to the underlying data.

When you are given an interface value, there’s no guarantee whether the underlying type is or isn’t a pointer.

### channels in go

Channels are typed values that allow goroutines to synchronize and exchange information.

```go
    c := make(chan int)
    go func() { c <- 3} ()
    n := <-c
```

#### fan-in example in go

A `select` blocks until one of its `case` can execute.

```go
func main() {
    c := selectFanIn(sleepCounter("SFO", 2), sleepCounter("AUS", 4))
    for i := 1; i < 15; i++ {
        fmt.Println(<-c)
    }

    fmt.Println("Excellent example of an fan-in function")
}

// Combines the output of both the channels
func selectFanIn(input1 <-chan string, input2 <-chan string) <-chan string {
    c := make(chan string)
    go func() {
        for {
			// Send to channel c, recv from input1
			select {
			case s := <-input1:
				c <- s
			case s := <-input2:
				c <- s
			}
		}
	}()

	return c
}

// Returns receive-only channel of strings
func sleepCounter(name string, wait int) <-chan string {
	c := make(chan string)

	// goroutine inside a function
	go func() {
		for i := 1; ; i++ {
			c <- fmt.Sprintf("%s %d", name, i)
			time.Sleep(time.Duration(wait) * time.Millisecond)
		}
	}()
	return c // return the channel to the caller
}
```

### Reflection in go

* Reflection in Go is used for determining information at runtime. We use the `reflect` package.

1. There are three steps involved when using reflect. First, we need to convert an interface to reflect types (`reflect.Type` or `reflect.Value`, this depends on the situation).

```go
  t := reflect.TypeOf(i)    // get meta-data in type i, and use t to get all elements
  v := reflect.ValueOf(i)   // get actual value in type i, and use v to change its value
```

1. Convert the reflected types to get the values that we need

```go
  var x float64 = 3.4
  v := reflect.ValueOf(x)
  fmt.Println("type:", v.Type())
  fmt.Println("kind is float64:", v.Kind() == reflect.Float64)
  fmt.Println("value:", v.Float())
```

1. Change the values of the reflected types, we need to make it modifiable

```go
  var x float64 = 3.4
  p := reflect.ValueOf(&x)
  v := p.Elem()
  v.SetFloat(7.1)
```

### go naming conventions

* In Golang, any variable (or a function) with an identifier starting with an upper-case letter (example, CamelCase) is made public (accessible) to all other packages in your program, whereas those starting with a lower-case letter (example, camelCase) is not accessible to any package except the one it is being declared in.

* Names are as important in Go as in any other language. They even have semantic effect: the visibility of a name outside a package is determined by whether its first character is upper case. It's therefore worth spending a little time talking about naming conventions in Go programs.

#### Package

* When a package is imported, the package name becomes an accessor for the contents.

* After import "bytes" the importing package can talk about `bytes.Buffer`. It's helpful if everyone using the package can use the same name to refer to its contents, which implies that the package name should be good: short, concise, evocative. By convention, packages are given lower case, single-word names; there should be no need for underscores or mixedCaps. Err on the side of brevity, since everyone using your package will be typing that name. And don't worry about collisions a priori. The package name is only the default name for imports; it need not be unique across all source code, and in the rare case of a collision the importing package can choose a different name to use locally. In any case, confusion is rare because the file name in the import determines just which package is being used.

* Another convention is that the package name is the base name of its source directory; the package in `src/encoding/base64` is imported as "encoding/base64" but has name `base64`, not encoding_base64 and not encodingBase64.

* The importer of a package will use the name to refer to its contents, so exported names in the package can use that fact to avoid stutter. (Don't use the import . notation, which can simplify tests that must run outside the package they are testing, but should otherwise be avoided.) For instance, the buffered reader type in the `bufio` package is called `Reader`, not BufReader, because users see it as `bufio.Reader`, which is a clear, concise name. Moreover, because imported entities are always addressed with their package name, `bufio.Reader` does not conflict with `io.Reader`.

* Similarly, the function to make new instances of `ring.Ring` which is the definition of a constructor in Go would normally be called `NewRing`, but since `Ring` is the only type exported by the package, and since the package is called `ring`, it's called just `New`, which clients of the package see as `ring.New`. Use the package structure to help you choose good names.

* Another short example is `once.Do`; `once.Do(setup)` reads well and would not be improved by writing `once.DoOrWaitUntilDone(setup)`. Long names don't automatically make things more readable. A helpful doc comment can often be more valuable than an extra long name.

#### Getters

Go doesn't provide automatic support for getters and setters. There's nothing wrong with providing getters and setters yourself, and it's often appropriate to do so, but it's neither idiomatic nor necessary to put Get into the getter's name. If you have a field called owner (lower case, unexported), the getter method should be called `Owner` (upper case, exported), not `GetOwner`. The use of upper-case names for export provides the hook to discriminate the field from the method. A setter function, if needed, will likely be called `SetOwner`. Both names read well in practice:

```go
owner := obj.Owner()
if owner != user {
    obj.SetOwner(user)
}
```

#### Interfaces

* By convention, one-method interfaces are named by the method name plus an `-er` suffix or similar modification to construct an agent noun: Reader, Writer, Formatter, Stringer, CloseNotifier etc.

* There are a number of such names and it's productive to honor them and the function names they capture. Read, Write, Close, Flush, String and so on have canonical signatures and meanings. To avoid confusion, don't give your method one of those names unless it has the same signature and meaning. Conversely, if your type implements a method with the same meaning as a method on a well-known type, give it the same name and signature; call your string-converter method `String` not `ToString`.

#### MixedCaps

The convention in Go is to use `MixedCaps` or `mixedCaps` rather than underscores to write multiword names.

### GoLang Tools

#### Command stringer

Stringer is a tool to automate the creation of methods that satisfy the fmt.Stringer interface. Given the name of a (signed or unsigned) integer type T that has constants defined, stringer will create a new self-contained Go source file implementing

```go
func (t T) String() string

//go:generate stringer -type=Pill
```
### golang references

* [multiple-return-values] (https://gobyexample.com/multiple-return-values)
* [golang] (https://golang.org/)
* [overview-of-Go] (https://www.codeproject.com/Articles/1160635/Quick-Overview-of-Go)
* [golang-book] (https://www.golang-book.com/books/intro/10)
* [go-lang-cheat-sheet] (https://github.com/a8m/go-lang-cheat-sheet/)
* [godoc.org] (https://godoc.org/)

### golang useful packages

* [Package goth provides a simple, clean, and idiomatic way to write authentication packages for Go web applications] (https://github.com/markbates/goth)
* [Go Secure Example Application] (https://github.com/komand/gosea)

### Rapid Web development in golang using Buffalo

* [gobuffalo.io] (http://gobuffalo.io/docs/getting-started)

#### Creating a new application using Buffalo

```bash

navkar$ go get -u github.com/gobuffalo/buffalo/...

navkar$ pwd
/Users/navkar/go/src/github.com/navkar

navkar$ buffalo new trips
...
You can find your new application at:
/Users/navkar/go/src/github.com/navkar/trips

navkar$
```

#### Creating databases

Once the `database.yml` has been configured with the appropriate settings, and the database server is running, Buffalo can create all of the databases in the `database.yml` file with a simple command:

```bash
navkar$ pwd
/Users/navkar/go/src/github.com/navkar/trips

trips navkar$ buffalo db create -e development
v3.17.0

CREATE DATABASE trips_development;
created database trips_development
```

#### Running the GO Web App

```bash
navkar$ buffalo dev
```

#### Running the GO Web App

![Buffalo](/img/golang_buffalo.png?raw=true "Buffalo Dev")

#### Create Models

Buffalo uses Fizz to generate migrations that are both easy to work with and work across multiple types of databases.

```bash
navkar$ buffalo db generate model trip
v3.17.0

--> models/trip.go
--> models/trip_test.go
--> goimports -w .
> migrations/20170413183709_create_trips.up.fizz
> migrations/20170413183709_create_trips.down.fizz
```

Create a `trip` model.

```golang
create_table("trips", func(t) {
    t.Column("id", "uuid", {"primary": true})
    t.Column("name", "string", {"size": 50})
    t.Column("location", "string", {"size": 50})
    t.Column("duration", "integer", {"default": 0})
    t.Column("is_active", "boolean", {"default": true})
    t.Column("is_deleted", "boolean", {"default": false})
    t.Column("started_at", "timestamp", {})
    t.Column("ended_at", "timestamp", {})
})
```

#### Running Migrations

Once migrations have been created they can be run with either of the following commands:

```bash
navkar$ buffalo db migrate
v3.17.0

> 20170413183709_create_trips.up.fizz

0.0605 seconds
```



