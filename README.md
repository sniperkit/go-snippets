knmr
===

knmr - Command Line Utility written in Go

# DESCRIPTION

knmr is command-line interface tool in Go. 

# INSTALLATION

## Build from source

```bash
$ go get github.com/Ken2mer/knmr
$ go install github.com/Ken2mer/knmr
```

## EXAMPLES

```
knmr agent
```
* every 3s tick until signal receive

The `knmr agent` command has an '-d' option for debug output.

```
knmr -d agent
```
* log output in DEBUG level

# CONTRIBUTION

1. Fork ([https://github.com/ken2mer/knmr/fork](https://github.com/ken2mer/knmr/fork))
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Run test suite with the `go test ./...` command and confirm that it passes
6. Run `gofmt -s`
7. Create new Pull Request

# TODO

## Implementation of loadConfig

* use TOML parser

cf. TOML parser  
[BurntSushi/toml: TOML parser for Golang with reflection.](https://github.com/BurntSushi/toml)

## As a library

* go-libs
* Go client libraries for using third-party APIs

cf. awesome-go#third-party-apis  
[avelino/awesome-go: A curated list of awesome Go frameworks, libraries and software](https://github.com/avelino/awesome-go#third-party-apis)

## Additional functions

* auth, login
* upload, postFile
* session

cf. 
[LearnServerProgramming · golang/go Wiki](https://github.com/golang/go/wiki/LearnServerProgramming)

cf. 
[Custom Handlers and Avoiding Globals in Go Web Applications · request / response](https://blog.questionable.services/article/custom-handlers-avoiding-globals/)

* How to build your own handler type
* How to explicitly pass a “context” containing our database pool, template map, a custom logger, ...
