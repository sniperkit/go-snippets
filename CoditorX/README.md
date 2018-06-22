# CoditorX [![Build Status](https://img.shields.io/travis/gophergala/CoditorX.svg?style=flat)](https://travis-ci.org/gophergala/CoditorX)

## Introduction

CoditorX is an in-browser markdown editor for multiple players collaborating.

_Have a try first: [http://coditorx.b3log.org](http://coditorx.b3log.org)_

## Features

* Allow user to register, everyone has its own workspace
* User can share documents to others with edit or view permission
* Real-time synchronize content to the users with the same shared document
* Markdown preview

## Architecture

### Overview

![architecture](https://github.com/gophergala/CoditorX/raw/master/architecture.png)

### Key points

* Efficient protocol: websocket
* Version control: version number, incremental changes/merges
* Permission control: private, shared edit/view

## Setup

### Build CoditorX for yourself

1. [Download](https://github.com/gophergala/CoditorX/master.zip) source or by `git clone`
2. Get dependencies with `go get`
3. Compile CoditorX with `go build`
4. then run ./CoditorX or CoditorX.exe

### Docker

1. Get image: `sudo docker pull 88250/coditorx:latest`
2. Run: `sudo docker run -p 127.0.0.1:9090:9090 88250/coditorx:latest ./CoditorX -channel=ws://127.0.0.1:9090`
3. Open browser: http://127.0.0.1:9090

## License

Copyright (c) 2015, B3log Team (http://b3log.org)

Licensed under the [Apache License 2.0](https://github.com/b3log/wide/blob/master/LICENSE).

## Credits

* [Wide](https://github.com/b3log/wide)
* [CodeMirror](https://github.com/codemirror/codemirror)
* [go-diff](https://github.com/sergi/go-diff)
* [Gorilla](https://github.com/gorilla)
* [Docker](https://docker.com)
