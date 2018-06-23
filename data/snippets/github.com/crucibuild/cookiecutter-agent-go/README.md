cookiecutter-agent-go
=====================
[![Cookiecutter](https://img.shields.io/badge/generator-cookiecutter-5aadbb.svg?style=flat)](http://yeoman.io)
[![Build Status](https://travis-ci.org/crucibuild/cookiecutter-agent-go.svg?branch=master)](https://travis-ci.org/crucibuild/cookiecutter-agent-go)

> A [cookiecutter] template to create new [Crucibuild] Agents in [Go] following best practices.

## Features

The template scaffolds a _project_ with a minimal implementation of a Crucibuild agent, including:

 * the agent manifest
 * a Makefile managing the project build:
    * installation of the all necessary go dependencies
    * code generation of the static resources using [go-resources](https://github.com/omeid/go-resources/cmd/resources)
    * check of Go sources for errors and warnings using [gometalinter](https://github.com/alecthomas/gometalinter)
    * built of the agent
 * TravisCI integration.

## Installation

##### Prerequisites

The following must be installed and available:

- [Python]
- [Go]

##### Project
 
Install [cookiecutter] using the following command line:

```sh
$ pip install cookiecutter
```

alternatively you can install cookiecutter with `homebrew`:

```sh
$ brew install cookiecutter
```

## Usage

To create a new Crucibuild agent, just run `cookiecutter` with this template by typing:

```sh
$ cookiecutter https://github.com/crucibuild/cookiecutter-agent-go
```

You'll be prompted for various configuration options (see `cookiecutter.json` for the full list).

Finally, enter the project and take a look around:

```sh
$ cd <agent-name>
$ ls
```

You can build the agent using the make command:
```sh
$ make
```

[Python]: https://www.python.org/
[Go]: https://golang.org/
[cookiecutter]: https://github.com/audreyr/cookiecutter
[Crucibuild]: https://github.com/crucibuild
