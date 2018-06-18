# Errors Package

Errors package is a copy of [Upspin errors package](https://github.com/upspin/upspin/blob/master/errors/errors.go) with some modification.

## How it is different from standard error package?

Errors package contains `Fields` which is rather helpful to add more context into `error`

```go
err := errors.New("This is new error", errors.Fields{"field1": "value1"})
fields := err.GetFields()
```

The implementation is pretty much like `logrus.Fields`, but is effective to add more context.

## Match function

Error with same string but from different `interface{}` implementation will not matched, so `Match` function is needed.

To check if error is actually match:

```go
import (
    stderr "errors"
    "github.com/albert-widi/go_common/errros"
)

func main() {
    err := stderr.New("This is an error")
    errs := errors.New("This is an error")
    // this will not work
    if err == errs {
        // do something
    }
    // do this instead
    if errors.Match(err, errs) {
        // do something
    }
}
```

## Runtime output

Errors provide a function called `SetRuntimeOutput`, when this on errors will automatically record the file and line where `error` is happened. This is enabled by calling runtime function.

Please note that this is experimental and thousands of runtime call might be very-very expensive

```go
errors.SetRuntimeOutput(true)
```

## Will this help you in the long run?

Yes and no, depends on your mental model. It depends on what you're gonna build, if you're building a service then yes maybe this is gonna help. But for a library, this kind of things will be an `overkill`, use standard `error` pakcage instead.