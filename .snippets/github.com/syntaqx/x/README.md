# X - experiments

[![GoDoc](https://godoc.org/github.com/syntaqx/x?status.svg)](https://godoc.org/github.com/syntaqx/x)
[![Build Status](https://travis-ci.org/syntaqx/x.svg?branch=master)](https://travis-ci.org/syntaqx/x)
[![codecov](https://codecov.io/gh/syntaqx/x/branch/master/graph/badge.svg)](https://codecov.io/gh/syntaqx/x)
[![Go Report Card](https://goreportcard.com/badge/github.com/syntaqx/x)](https://goreportcard.com/report/github.com/syntaqx/x)
[![Sourcegraph for Repo Reference Count](https://sourcegraph.com/github.com/syntaqx/x/-/badge.svg)](https://sourcegraph.com/github.com/syntaqx/x?badge)

[license]: https://syntaqx.com/license

`x` is a collection of experimental packages and prototypes. The purpose of this
repository is provide technical journaling for projects I'm building, intend to
build, or have abandoned while building.

Some packages may one day be promoted to their own repositories or merged with
existing software, or they be modified arbitrarily or even disappear altogether.
In short, code in this repository makes no compatibility promises, and is very
likely to do so. __All usage is at your own risk__. ¯\\\_(ツ)\_/¯

## Packages

[bytes]: https://github.com/syntaqx/x/tree/master/bytes
[render]: https://github.com/syntaqx/x/tree/master/render

| Package    | Status                        | Synopsis                                                                        |
|------------|-------------------------------|:--------------------------------------------------------------------------------|
| [bytes][]  | ![status][status-development] | provides functionality for measuring and displaying data units                  |
| [render][] | ![status][status-development] | wraps [unrolled/render][unrolled/render] with sugar                             |

[unrolled/render]: https://github.com/unrolled/render

### Under consideration

| Package      | Status                    | Synopsis                                                       |
|--------------|---------------------------|:---------------------------------------------------------------|
| Mux/Router++ | ![status][status-concept] | I kind of wanna just build one...                              |

## License

Released under the [MIT license][license]

[status-concept]: https://img.shields.io/badge/status-concept-lightgrey.svg
[status-research]: https://img.shields.io/badge/status-research-orange.svg
[status-scaffolding]: https://img.shields.io/badge/status-scaffolding-yellow.svg
[status-development]: https://img.shields.io/badge/status-development-yellowgreen.svg
[status-promotable]: https://img.shields.io/badge/status-promotable-green.svg
[status-abandoned]: https://img.shields.io/badge/status-abandoned-red.svg
