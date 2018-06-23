# Go by Example: Variables

In Go, variables are explicitly declared and used by the compiler to
e.g. check type-correctness of function calls.

`var` declares 1 or more variables.

You can declare multiple variables at once.

Go will infer the type of initialized variables.

Variables declared without a corresponding initialization are _zero-valued_.
For example, the zero value for an `int` is 0.

The `:=` syntax is shorthand for declaring and initializing a variable,
e.g. `for var f string = "short"` in this case.

```sh
$ go run variables.go
initial
1 2
true
0
short
```
