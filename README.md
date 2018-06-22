# promise [![Build Status](https://travis-ci.org/matthewdale/promise.svg?branch=master)](https://travis-ci.org/matthewdale/promise) [![codecov](https://codecov.io/gh/matthewdale/promise/branch/master/graph/badge.svg)](https://codecov.io/gh/matthewdale/promise) [![Go Report Card](https://goreportcard.com/badge/github.com/matthewdale/promise)](https://goreportcard.com/report/github.com/matthewdale/promise) [![GoDoc](https://godoc.org/github.com/matthewdale/promise?status.svg)](https://godoc.org/github.com/matthewdale/promise)

A simple, fast implementation of a promise in Go.

Example
=======
```go
// do something that takes 10 seconds to return
func do() (string, error) {
    time.Sleep(10 * time.Second)
    return "result", nil
}

func doAsync() *promise.Promise {
    p := promise.NewPromise()

    // Start a goroutine that completes the promise
    // once the result is returned.
    go func() {
        v, err := do()
        if err != nil {
            p.CompleteWithError(err)
            return
        }
        p.Complete(v)
    }()
    return p
}

func main() {
    p := doAsync()

    // Calling 'Get' blocks until the result is ready.
    v, err := p.Get()
    if err != nil {
        panic(err)
    }
    if result, ok := v.(string); ok {
        fmt.Printf("Promise result: %s", result)
    }
}
```
