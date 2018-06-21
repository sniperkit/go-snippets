
`master` mirrors v2 branch, to install:

	go get -u github.com/gizak/termui

## Introduction
This is a project for kernel source code analysis:

For single host emulation you just need to `make boot`.

For multi-host emulation you just need to `make mboot`, this will create two host and a router. And the topoloy will be:
__Absolute layout__
````c
    ---------             --------            -------
    | host1 |------------>|router|<-----------|host2|
    ---------             --------            -------
````

## Installation
    git clone https://github.com/lamproae/kemu.git 

## Use '. bootstrap.sh' to initialize the project and set related environment variables.
## Use '. script/setup-env.sh' to initialize the basic env.

## How to use
After downloading the source code:
  - [1]. cd kemu
  - [2]. . bootstrap.sh
  - [3]. cd kdebug
  - [2]. . script/setup-env.sh
  - [3]. make
  - [4]. make boot (make mboot)

`You can telnet to the mulated host from you local host. This will make the terminal looks better.`
