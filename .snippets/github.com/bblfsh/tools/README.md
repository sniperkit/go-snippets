# Babelfish Tools

[![Build Status](https://travis-ci.org/bblfsh/tools.svg?branch=master)](https://travis-ci.org/bblfsh/tools)
[![codecov](https://codecov.io/gh/bblfsh/tools/branch/master/graph/badge.svg)](https://codecov.io/gh/bblfsh/tools)

Language analysis tools on top of Babelfish

## Build

```sh
go get -u -v -d github.com/bblfsh/tools/...
```

If you want to build a docker image containing Babelfish Tools, run:

`make build`

If you clone the repo manually, remember that it should be in your 
`$GOPATH/src/github.com/bblfsh` directory to successfully build.

## Usage

Babelfish Tools provides a set of tools built on top of Babelfish, to
see which tools are supported, run:

`bblfsh-tools --help`

To make use of any of these tools you need to have the Babelfish
server up and running. The easiest way to do so is through docker:

`docker run --privileged -p 9432:9432 --name bblfsh bblfsh/server`

Look at [server site](https://github.com/bblfsh/server/) for more
information.

Once you have a server running, you can use the dummy tool, which
should let you know if the connection with the server succeeded:

`bblfsh-tools dummy path/to/source/code`

If the server is in a different location, use the `address` parameter:

`bblfsh-tools dummy --address location:port path/to/source/code`

Once the connection with the server is working fine, you can use any other
available tool in a similar way.

### Available tools

Apart from the dummy tool, the following tools are currently provided:

* cyclomatic: Parses a code file and prints its
  [cyclomatic complexity](https://en.wikipedia.org/wiki/Cyclomatic_complexity)
* npath: Parses a code file and prints the
  [npath complexity](https://pmd.github.io/pmd-5.7.0/pmd-java/xref/net/sourceforge/pmd/lang/java/rule/codesize/NPathComplexityRule.html)
  of its functions
* tokenizer: Parses a code file and extracts and prints its tokens

## How to add a new tool to Babelfish Tools

Adding a new tool to Babelfish Tools involves two steps: implementing
the Tool interface and adding it as a command to the CLI interface.

### Implementing the Tooler interface

The `Tooler` interface has a single method `Exec(*uast.Node) error`,
that is, a tool must implement a method called `Exec` that receives a
pointer to an UAST node and returns an optional `error`.

It's also convenient to create a new type for the new tool, to be used
in the CLI interface command. In the simplest case, an empty struct
will do: `type Dummy struct{}`

### Adding the new tool as a command to the CLI interface

Create a new file for the tool command in the `cmd/bblfsh-tools`
directory.

Create a new type there for your tool. It should at least include
Common struct. If your tool will support additional parameters, add
them there. In the simplest case, it'd just be:

```go
type Dummy struct {
	Common
}
```

Implement the
[Commander interface](https://godoc.org/github.com/jessevdk/go-flags#Commander). There's
a helper common method `execute` that will provides a good default, in
most cases calling this method should be enough:

```go
func (c *Dummy) Execute(args []string) error {
	return c.execute(args, tools.Dummy{})
}
```

Note that `tools.Dummy{}` is the instance of the type that implements
the `Tooler` interface that we described in the previous section.

At this point, only adding the command to the parser is left. This is
done at `cmd/bblfsh-tools/main.go`:

```go
parser.AddCommand("dummy", "", "Run dummy tool", &Dummy{})
```

And that's it, rebuild and your new tool should be ready to use.

## License

GPLv3, see [LICENSE](LICENSE)
