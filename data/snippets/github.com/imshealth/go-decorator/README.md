# Go-Decorator

`go-decorator` creates an implementation of an interface suitable for
the [decorator pattern](https://en.wikipedia.org/wiki/Decorator_pattern).

# Motivation
Rather than adding logging to every method of an api, write one logging method
and add it to every method of an interface.  This could also be used for
* [circuit-breaking](https://godoc.org/github.com/afex/hystrix-go/hystrix)
* [backoff](https://github.com/cenkalti/backoff)
* [metrics](https://godoc.org/github.com/afex/hystrix-go/hystrix/metric_collector)

# Using go-decorator

## Generate the decorator
    $ go-decorator -type MyApi myapi.go > myapi_decorator.go

Go-decorator will include the necessary imports, and place the resulting
structure in the same package.

If the generated import list is missing a library, use the `-import` flag to add it.

    $ go-decorator -type MyApi -import github.com/... myapi.go

## In code
```
    logMethod := func(name string, call func() error) {
        err := call()
        log.Infof("Called Method %s, error? %v", name, err)
        return err
    }

    api := &MyApi{}
    api = MyApiDecorator{Inner: api, Decorator: logMethod}

    // use API as before, it implements the same interface!
    ...
```

# Contributing

Missing a feature? [Open an issue](issues) or fork and send over a pull request.

# License
This work is published under the MIT license.

Please see the `LICENSE` file for details.
