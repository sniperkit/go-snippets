[![Datalayer](https://business.datalayer.io/img/logo-datalayer-horizontal.png)](http://datalayer.io)

[WIP] Collection of `golang` samples to learn it fast...

[![Go Report Card](https://goreportcard.com/badge/github.com/datalayer/go-snippets)](https://goreportcard.com/report/github.com/datalayer/go-snippets)

## Dep

### dep init

Initialize the project at filepath root by parsing its dependencies, writing
manifest and lock files, and vendoring the dependencies. If root isn't
specified, use the current directory.

When configuration for another dependency management tool is detected, it is
imported into the initial manifest and lock. Use the -skip-tools flag to
disable this behavior. The following external tools are supported:
glide, godep, vndr, govend, gb, gvt, govendor, glock.

Any dependencies that are not constrained by external configuration use the
GOPATH analysis below.

By default, the dependencies are resolved over the network. A version will be
selected from the versions available from the upstream source per the following
algorithm:

 - Tags conforming to semver (sorted by semver rules)
 - Default branch(es) (sorted lexicographically)
 - Non-semver tags (sorted lexicographically)

An alternate mode can be activated by passing -gopath. In this mode, the version
of each dependency will reflect the current state of the GOPATH. If a dependency
doesn't exist in the GOPATH, a version will be selected based on the above
network version selection algorithm.

A Gopkg.toml file will be written with inferred version constraints for all
direct dependencies. Gopkg.lock will be written with precise versions, and
vendor/ will be populated with the precise versions written to Gopkg.lock.

Flags:

  -gopath       search in GOPATH for dependencies (default: false)
  -no-examples  don't include example in Gopkg.toml (default: false)
  -skip-tools   skip importing configuration from other dependency managers (default: false)
  -v            enable verbose logging (default: false)

### Dep ensure

`dep ensure -update`

## Gorilla WebSocket

```
go get github.com/gorilla/websocket
cd `go list -f '{{.Dir}}' github.com/gorilla/websocket/examples/chat`
go run *.go
open http://localhost:8080
```

```
go get github.com/gorilla/websocket
cd `go list -f '{{.Dir}}' github.com/gorilla/websocket/examples/filewatch`
go run main.go <name of file to watch>
open http://localhost:8080
# Modify the file to see it update in the browser.
```
