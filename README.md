# new
new is a command to create projects off of github repos, it simply downloads a github repository and extracts it to a directory with your specified project name

[![Build Status](https://travis-ci.org/divanvisagie/new.svg?branch=master)](https://travis-ci.org/divanvisagie/new)


## Install 

#### Requirements 
 - The git command must be accessable from the terminal you execute `new` in

### Windows 

First install [scoop](http://scoop.sh/)

```
scoop bucket add divanvisagie https://github.com/divanvisagie/scoop-bucket
scoop install new
```

### macOS

#### Homebrew

First install [homebrew](https://brew.sh/)

```
brew install divanvisagie/homebrew-tap/new
```

#### Manual Install 

Download the tar.gz and run `install.sh`


### Linux
#### Linuxbrew

First install [linuxbrew](http://linuxbrew.sh/)

```
brew install divanvisagie/linuxbrew-tap/new
```

#### Manual Install

Download the tar.gz and run `install.sh`



## Run

To run , simply type the new command, the name of the project you want to generate and the github repo name that you want to seed from

```sh
new myProjectName divanvisagie/kotlin-tested-seed
```
