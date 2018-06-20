# genericop

[![GoDoc](https://godoc.org/github.com/motemen/go-genericop?status.svg)](https://godoc.org/github.com/motemen/go-genericop)

Package genericop provides functions that does operators (e.g. "+" or ">")
on interface{} args. Each function returns the result along with an error,
which indicates if the argument types were incompatible on the operation.
You can use MustT(x, err) family to assert the resulting type.

## Examples

### Add

```go
fmt.Println(Add(1, 2))
fmt.Println(Add("a", "b"))
fmt.Println(Add(1+1i, 2+2i))
fmt.Println(Add(1.1, "x"))
```

Output:
```
3 <nil>
ab <nil>
(3+3i) <nil>
<nil> incompatible types: float64 + string

```

### MustInt

```go
MustInt(Add(1, 2))
MustInt(Add(1.1, 2.2))
MustInt(Add("a", 2))
```

## Author

motemen <https://github.com/motemen>
