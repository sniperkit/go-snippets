Useful [pre-commit](http://pre-commit.com) hooks written in Go language.

Currently it only includes [mrtazz Makefile linter](https://github.com/mrtazz/checkmake).

Requires that you have `go` installed.

## Usage

```
-   repo: https://github.com/Lucas-C/pre-commit-hooks-go
    sha: v1.0.0
    hooks:
    -   id: checkmake
```

## Development

`checkmake.go` is simply a copy of `cmd/checkmake/main.go` in https://github.com/mrtazz/checkmake.

To update it:

    curl https://raw.githubusercontent.com/mrtazz/checkmake/master/cmd/checkmake/main.go -o checkmake.go

