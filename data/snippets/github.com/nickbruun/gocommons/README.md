# Go commons

[![GoDoc](https://godoc.org/github.com/nickbruun/gocommons?status.svg)](https://godoc.org/github.com/nickbruun/gocommons) [![Build status](https://travis-ci.org/nickbruun/gocommons.svg?branch=master)](https://travis-ci.org/nickbruun/gocommons)

Common libraries for the type of Go applications I toy with and build.


## Contents

The libraries are a mixture of foundation, utility and system libraries for various purposes.


### Foundation

The foundation libraries are libraries on which all other libraries in the collection depends.

**logging**

A thin wrapper library around [Simon Eskildsen's](http://sirupsen.com/) excellent [logrus](https://github.com/Sirupsen/logrus) that provides caller annotations to make it easier to figure out where log messages are coming from:

```Go
package main

import (
	logrus "github.com/Sirupsen/logrus"
	log "github.com/nickbruun/gocommons/logging"
	"sync"
)

var wg sync.WaitGroup

func MyFunc() {
	log.Debug("Debug time!")

	wg.Add(1)
	go func() {
		log.Info("Goroutines are presented understandably as well")
		wg.Done()
	}()
	wg.Wait()
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	log.Info("Hola, mundo")
	MyFunc()
}
```

Results in:

    INFO[0000] Hola, mundo                                       file=/Users/nick/Code/gocommonsdev/src/github.com/nickbruun/gocommons/play/main.go func=main.main line=25
    DEBU[0000] Debug time!                                       file=/Users/nick/Code/gocommonsdev/src/github.com/nickbruun/gocommons/play/main.go func=main.MyFunc line=12
    INFO[0000] Go-routines are presented understandable as well  file=/Users/nick/Code/gocommonsdev/src/github.com/nickbruun/gocommons/play/main.go func=main.func·001 line=16


### Utilities

**zkutils**

Utilities for interacting with and testing integrations with ZooKeeper, built on top of [go-zookeeper](https://github.com/samuel/go-zookeeper). Among other things, the library provides a connection manager wrapper class, that allows multiplexed access to all events received by a connection. Derived from that is a session watcher, which can be used to identify when a connection loss to a ZooKeeper cluster brings the validity of an ephemeral node into question.


### Systems libraries

#### Distributed systems libraries

**distributed/locking**

Distributed mutually exclusive locking facilities. Currently only implemented for ZooKeeper, with a cautious approach to the validity of ephemeral nodes when clusters are unstable.

**distributed/leadership**

Leadership election and observation facilities. Currently only implemented for ZooKeeper as an extended version of the [leader election recipe](http://zookeeper.apache.org/doc/trunk/recipes.html#sc_leaderElection), which takes a cautious approach to the validity of ephemeral nodes when clusters are unstable, preferring giving up leadership than risking races.
