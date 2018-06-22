# CoditorX

## Introduction

CoditorX is an online text editor in browser for multiple players collaborating.

## Features

* allow user to register, everyone has its own workspace
* user can share documents to others with edit or view permission
* real-time synchronize content to the users with the same shared document
* Markdown preview

## Architecture

### Overview

![architecture](architecture.png)

### Key points

* efficient protocol: websocket
* version control: version number, incremental changes/merges
* permission control: private, shared edit/view

## Dependencies

### Frontend

* [Quill](https://github.com/quilljs/quill) (web editor)
* [jQuery](http://jquery.com) (utilities)

### Backend

* [Wide](https://github.com/b3log/wide) (Go utilities, such as logging/filesys)
* [go-diff](https://github.com/sergi/go-diff) (Google's diff-match-patch library to Go)
* [Gorilla](https://github.com/gorilla) (Go web lib, such as websocket/session)

