/*
<!--
Copyright (c) 2017 Christoph Berger. Some rights reserved.

Use of the text in this file is governed by a Creative Commons Attribution Non-Commercial
Share-Alike License that can be found in the LICENSE.txt file.

Use of the code in this file is governed by a BSD 3-clause license that can be found
in the LICENSE.txt file.

The source code contained in this file may import third-party source code
whose licenses are provided in the respective license files.
-->

<!--
NOTE: The comments in this file are NOT godoc compliant. This is not an oversight.

Comments and code in this file are used for describing and explaining a particular topic to the reader. While this file is a syntactically valid Go source file, its main purpose is to get converted into a blog article. The comments were created for learning and not for code documentation.
-->

+++
title = "goman - the missing man pages for Go binaries"
description = "A Go binary without a man page? goman displays the project's README file instead."
author = "Christoph Berger"
email = "chris@appliedgo.net"
date = "2017-06-26"
draft = "false"
domains = ["DevOps"]
tags = ["man", "man page"]
categories = ["Tools and Libraries"]
+++


Most Go binaries come without any man page. The tool `goman` fills this gap. If the corresponding project includes a decent README file (and most projects do), `goman` find this README file and displays it on the terminal.

<!--more-->

It kept happening to me: I type `man <blah>` to get the man page of `<blah>`, only to find out that `<blah>` is a Go binary and hence has no man page (except for rare cases where the author took the time to write one and distribute it along with the binary via some installation manager like [Homebrew](https://brew.sh).)

Well, not anymore.

I wrote a tool named [`goman`](https://github.com/christophberger/goman) to get me some info about a Go binary when `man` can't.

![goman logo](goman.png)

From now on, when `man` finds no man page, `goman` will take a second attempt and try displaying the README file from the binary's project sources. This usually succeeds if the binary is a Go binary, and if this binary has a corresponding project either locally or on GitHub or GitLab.

This demo shows

* `goman`,
* `goman` with `less -R`, and
* `goman` launched by a bash script if `man` finds no man page.

[![goman demo](goman.gif)](/media/goman/goman.gif)
*(Click to enlarge)*

Here is how it works.


## Step 1: Locate the binary

A Go binary must reside within one of the directories contained in `$PATH`, or else I would not be able to call it from the command line. And when I run `man <binary>`, then the binary must be the first one of that name. So to find the binary, we simply need to do the same as the command `which` does.

Luckily, this is as easy as go-getting the `which` package from [bfontaine/which](https://github.com/bfontaine/which) and calling `which.One(name)`. Additionally, `goman` also calls `which.OneWithPath(name, path)` in order to locate the binary also in `$GOPATH/bin` or in the current path, in case any of the two is not included in `$PATH`.


## Step 2: Get the binary to reveal its source

Here is a maybe little-known fact: Go binaries contain the path of their source code project. With the help of the packages `debug/elf`, `debug/gosym`, and `debug/macho`, `goman` dives into the symbol table of the binary and locates the symbol "main.main". The text of this symbol is the path to the source code, relative to GOPATH (and sometimes it is an absolute path, but in this case, `goman` takes a guess and cuts off the path prefix up to the first occurrence of `/src/`, which is supposed to be the `/src/` directory right beneath GOPATH).

The code for reading the symbol table is gratefully taken from the [`gorebuild` tool](https://github.com/FiloSottile/gorebuild). As `gorebuild` is no library, I copied over `dwarf.go` that contains all the code I need for that behind a simple function call: `getMainPath(file)`.


## Step 3: Locate the README file

This part turned out to be a bit more complex. `goman` needs to find a README file that is either a Markdown file or a plain text file, which means we need to look for several possible file extensions: `.md`, `.markdown`, `.txt`, and no extension at all. (GitHub supports other markup langauges besides Markdown, but I find they are so rare that they can be safely ignored.)

In addition to that, a command can either reside in the project root, or in a `/cmd/abc` subproject. Sometimes the `/cmd/abc` subprojects contain their own README file but in most cases they don't. So `goman` has to check both the subbroject and the root project for a README file.

And if this is not enough, the source code may or may not exist on the local machine. The Go binary might have been installed via Homebrew or some other package manager, or the source code might have been removed after installing the binary.
In this case, `goman` must also look into the public repository of the project (which nowadays is mostly hosted on GitHub, with a few exceptions).

To make things even more complicated, the raw text version of the README file resides at a different path than the local README, so `goman` takes an additional step and builds the raw file URL from what it knows about the source path and about GitHub's and GitLab's raw file URL structure. (I had no joy with Bitbucket's file URL's, as they contain unpredictable hash values. I did not want to go as far as writing a web scraper to analyze the repository page in order to get the README URL.)

<!--
Last not least, the URL to the remote repository might be a canonical URL which gets redirected to a "real" URL.
-->

This diagram shows where `goman` searches for READMEs - if you seek some diversion for a few seconds, click the play button.

HYPE[how goman searches READMEs](goman.html)


## Step 4: Render Markdown as colored ANSI text

Markdown renderers usually produce HTML, in some cases PDF, and, if they value traditional typesetting, also LaTeX. Plain text with ANSI color codes is far less widespread.

Again, luck was on my side. [This fork of blackfriday](https://github.com/ec1oud/blackfriday) implements an ANSI renderer for Markdown. So getting color-coded ANSI output from a README file required little more than copying a couple of lines from [mdcat](https://github.com/ec1oud/mdcat) that uses the ANSI renderer under the hood.

Done.

## Get goman

Get `goman` via `go get`:

    go get -u github.com/christophberger/goman

## Usage

It's all in the [README](https://github.com/christophberger/goman/blob/master/README.md), so simply type `goman goman` ;-)

TL;DR:

```sh
goman <binary>  # find and display the README of <binary>

goman <binary> | less -R  # same but with paging
```

And a shell script can make `goman` blend in with the standard `man` command, to auto-reveal the README if the binary has no man page. (See the README for instructions.)

**Happy coding!**

*/
