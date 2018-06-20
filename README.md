[![GoDoc](https://godoc.org/github.com/mickep76/encdec?status.svg)](https://godoc.org/github.com/mickep76/encdec)
[![codecov](https://codecov.io/gh/mickep76/encdec/branch/master/graph/badge.svg)](https://codecov.io/gh/mickep76/encdec)
[![Build Status](https://travis-ci.org/mickep76/encdec.svg?branch=master)](https://travis-ci.org/mickep76/encdec)
[![Go Report Card](https://goreportcard.com/badge/github.com/mickep76/encdec)](https://goreportcard.com/report/github.com/mickep76/encdec)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/mickep76/mlfmt/blob/master/LICENSE)

# encdec

Package provides a generic interface to encoders and decoders

## Example

More examples can be found in [examples](https://github.com/mickep76/encdec/tree/master/examples).

```go
package main
  
import (
        "fmt"
        "log"

        "github.com/mickep76/encdec"
        _ "github.com/mickep76/encdec/yaml"
)

type Message struct {
        Name, Text string
}

type Messages []*Message

func main() {
        in := Messages{
                &Message{Name: "Ed", Text: "Knock knock."},
                &Message{Name: "Sam", Text: "Who's there?"},
                &Message{Name: "Ed", Text: "Go fmt."},
                &Message{Name: "Sam", Text: "Go fmt who?"},
                &Message{Name: "Ed", Text: "Go fmt yourself!"},
        }

        b, err := encdec.ToBytes("yaml", in)
        if err != nil {
                log.Fatal(err)
        }

        out := Messages{}
        if err := encdec.FromBytes("yaml", b, &out); err != nil {
                log.Fatal(err)
        }

        for _, m := range out {
                fmt.Printf("%s: %s\n", m.Name, m.Text)
        }
}
```
