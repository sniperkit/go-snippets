# Go by Example: Constants

Go supports _constants_ of character, string, boolean, and numeric values.

`const` declares a constant value.

A `const` statement can appear anywhere a `var` statement can.

Constant expressions perform arithmetic with arbitrary precision.

A numeric constant has no type until it's given one, such as by an explicit cast.

A number can be given a tpye by using it in a context that requires one, such as a variable
assignment or function call. For example, here `math.Sin` expects a `float64`.

```sh
$ go run constant.go
constant
6e+11
600000000000
-0.28470407323754404
```
