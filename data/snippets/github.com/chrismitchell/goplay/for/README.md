# Go by Example: For

`for` is Go's only looing construct. Here are three basic types of `for` loops.

The most basic type, with a single condition.

A classic initial/condition/after `for` loop.

`for` without a condition will loop repeatedly until you `break` out for the loop or `return` from
the enclosing function.

You can also `continue` to the next iteration of the loop.

```sh
$ go run for.go
1
2
3
7
8
9
loop
1
2
3
```
