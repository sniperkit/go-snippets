examples
========

This is a repository of interesting examples or tutorials.

## Installation

Installation should be performed via `go-get`:

```sh
$ go get github.com/go-hep/examples/...
```

(yes, the 3 dots are on purpose)

## Content

- `croot/go-croot-ex-tree-0X` where `X=[0-3]`: show how to read/write
  `ROOT::TTree`, involving builtins or structs.
  
- `croot/go-croot-ex-datareader` shows how an analysis code would look
  like when using the `go-croot-gen-datareader` code generator
  facility
  
- `croot/go-croot-ex-hist` shows how to fill a `ROOT::TH1F` histogram

- `hplot/hplot-ex-hist` shows how to fill and plot a `hist.Hist1D`
  histogram

