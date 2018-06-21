/*
<!--
Copyright (c) 2016 Christoph Berger. Some rights reserved.
Use of this text is governed by a Creative Commons Attribution Non-Commercial
Share-Alike License that can be found in the LICENSE.txt file.

The source code contained in this file may import third-party source code
whose licenses are provided in the respective license files.
-->

<!--
NOTE: The comments in this file are NOT godoc compliant. This is not an oversight.

Comments and code in this file are used for describing and explaining a particular topic to the reader. While this file is a syntactically valid Go source file, its main purpose is to get converted into a blog article. The comments were created for learning and not for code documentation.
-->

+++
title = "Plugins in Go"
description = "Go is a statically compiled language. No dynamic libraries can be loaded at runtime, nor does the runtime support compiling Go on the fly. Still, there is a number of ways of creating and using plugins in Go."
author = "Christoph Berger"
email = "chris@appliedgo.net"
date = "2016-06-30"
publishdate = "2016-06-30"
domains = ["Architecture"]
tags = ["Plugin", "RPC", "in-process", "out-of-process"]
categories = ["Survey"]
+++

Go is a statically compiled language. The Go runtime cannot load dynamic libraries, nor does it support compiling Go on the fly. Still, there is a number of ways of creating and using plugins in Go.

<!--more-->

![Low Tech Plugins](bits.jpg)


## The case for plugins in Go

Plugins are useful for extending an application's feature list - but in Go, compiling a whole app from source is fast and easy, so why should anyone bother with plugins in Go?

* First, loading a plugin at runtime may be a requirement in the app's technical specification.
* Second, fast compilation and plugins are no contradiction. Go plugins can be created to be compiled into the binary - we'll look at an example of this later.

This article is a quick survey of plugin architectures and techniques in Go.


## Plugin criteria

Here is a wishlist for the ideal plugin architecture:

* **Speed:** Calling a plugin's methods must be fast. The slower the method call is, the more is the plugin restricted to implementing only big, long-running, rarely called methods.
* **Reliability:** Plugins should not fail or crash, and if they do, recovery must be possible, fast, easy, and complete.
* **Security:** Plugins should be secured against tampering, for example, through code signing.
* **Ease of use:** The plugin programmer should not be burdened with a complicated, error-prone plugin API.

The ideal plugin architecture should meet all of the above criteria, but in real life there is usually one or another concession to make. This becomes immediately clear when we look at the question that should be the first one when deciding upon a plugin architecture:

*Shall the plugins run inside the main process, or rather be separate processes?*


## In-process vs separate processes

Both approaches have advantages and disadvantages, and as we'll see, one approach's disadvantage may be the other's advantage.


### Advantages of in-process plugins

* **Speed:** Method calls are as fast as can be.
* **Reliability:** The plugins are available as long as the main process is available. An in-process plugin cannot suddenly become unavailable at runtime.
* **Easy deployment:** The plugin gets deployed along with the binary, either baked right in, or (only in non-Go languages for now) as a dynamic shared library that can be loaded either at process start or during runtime.
* **Easy runtime management:** No need for discovering, starting, or stopping a plugin process. No need for health checks. (Does the plugin process still live? Does it hang? Does it need a restart?)


### Advantages of plugins as separate processes

* **Resilience:** A crashing plugin does not crash main process.
* **Security:** A plugin in a separate process cannot mess with internals of the main process.
* **Flexibility (part 1):** Plugins can be written in any language, as long as there is a library available for the plugin protocol.
* **Flexibility (part 2):** Plugins can be activated and deactivated during runtime. It is even possible to deploy and activate new plugins without having to restart the main process.

With these feature lists in mind, let's look at a couple of different plugin solutions for the Go language.


## Plugin approaches in Go

As mentioned before, Go lacks an option for loading shared libraries at runtime, and so a variety of alternate approaches have been created. Here are the ones I could find through two quick searches on GitHub and on Google, in no particular order:


### External process using RPC via stdin/stdout

#### Description

This is perhaps the most straightforward approach:

* Main process starts plugin process
* Main process and plugin process are connected via stdin and stdout
* Main process uses RPC ([Remote Procedure Call](https://en.wikipedia.org/wiki/Remote_procedure_call)) via stdin/stdout connection

#### Example

The blog post [Go Plugins are as Easy as Pie](http://npf.io/2015/05/pie/) introduced this concept to Go in May 2015. The accompanying `pie` package is [here](https://github.com/natefinch/pie), and if you ask me, this could be my favorite plugin approach just for the yummy pumpkin pie picture in the readme! (Spoiler picture below.)

<a title="By Evan-Amos (Own work) [CC BY-SA 3.0 (http://creativecommons.org/licenses/by-sa/3.0)], via Wikimedia Commons" href="https://commons.wikimedia.org/wiki/File%3APumpkin-Pie-Slice.jpg"><img width="10%" alt="Pumpkin-Pie-Slice" src="https://upload.wikimedia.org/wikipedia/commons/thumb/8/84/Pumpkin-Pie-Slice.jpg/512px-Pumpkin-Pie-Slice.jpg"/></a>

And this is basically how Pie starts a plugin and communicates with it:

![Pie plugin diagram](pie.svg)

In Pie, a plugin can take one of two roles.

* As a Provider, it responds to requests from the main program.
* As a Consumer, it can actively call into the main program and receive the results.




### External process using RPC via network

#### Description

The main difference to the previous approach is the way how the RPC calls are implemented. Rather than using the stdin/stdout connection, the RPC calls can also be done via the (local) network.

#### Example

The package [`go-plugin`](https://github.com/hashicorp/go-plugin) by HashiCorp utilizes `net/rpc` for connecting to the plugin processes. `go-plugin` is a rather heavyweight plugin system with lots of features, clearly able to attract developers of enterprise software who look for a complete and industry tested solution.

![go-plugin diagram](go-plugin.svg)


### External process via message queue

#### Description

Message queue systems, especially the brokerless ones, provide a solid groundwork for creating plugin systems. My quick research did not return any MQ-based plugin solution, but this may well be due to the fact that not much is needed to turn a message queue system into a plugin architecture.

#### Example

I did not find any message queue based plugin systems, but maybe you remember the [first post](https://appliedgo.net/messaging) of this blog, where I introduced the [nanomsg](http://nanomsg.org/) system and its Go implementation [Mangos](https://github.com/go-mangos/mangos). The nanomsg specification includes a set of predefined communication topologies (called "scalability protocols" in nanomsg jargon) covering many different scenarios: Pair, PubSub, Bus, Survey, Pipeline, and ReqRep. Two of them come in quite handy for communicating with plugins.

* The ReqRep (or Request-Reply) protocol can be used for mimicking RPC calls to a particular plugin. It is not the real RPC thing, however, as the sockets handle plain `[]byte` data only. So the main process and the plugins must take care of serializing and de-serializing the request and reply data.
* The Survey protocol helps monitoring the status of all plugins at once. The main process sends a survey to all plugins, and the plugins respond if they can. If a plugin does not respond within the deadline, the main process can take measures to restart the plugin.

![MQ based plugin](mq-plugin.svg)

### In-process plugins, included at compile time

#### Description

Calling a package a plugin might seem debatable when it is compiled into the main application just like any other package. But as long as there is a plugin API defined that the plugin packages implement, and as long as the build process is able to pick up any plugin that has been added, there is nothing wrong with that.

The advantages of in-process plugins--speed, reliability, ease of use--have been outlined above. As a downside, adding, removing, or updating a plugin requires compiling and deploying the whole main application.

#### Example

Technically, any go library package can be a plugin package provided that it adheres to the plugin API that you have defined for your project.

Maybe the most common type of compile-time plugin is HTTP middleware. Go's [`net/http`](https://golang.org/pkg/net/http/) makes it super easy to plug in new handlers to an HTTP server:

* Write a package containing either one or more functions that implement the `Handler` interface, or functions with the signature `func(w http.ResponseWriter, r *http.Request)`.
* Import the package into your application.
* Call `http.Handle(<pattern>, <yourpkg.yourhandler>)` or `http.HandleFunc(<pattern>, <yourpkg.yourhandlefunc>)`, respectively, to register a handler.

![in-process plugin](in-process.svg)

Needless to say that this pattern can be used for any kind of plugin; the concept is not specific to HTTP handlers.


### Script plugins: In-process but not compiled

#### Description

Script plugin mechanisms provide an interesting middle ground between in-process and out-of-process plugin approaches. The plugin is written in a scripting language whose interpreter is compiled into your application. With this technique it is possible to load an in-process plugin at runtime--with the small caveat that the plugin is not native code but must be interpreted. Expect most of these approaches to have a rather low performance.

#### Example

The page "Awesome-go.com" [lists a couple of embeddable scripting languages for Go](http://awesome-go.com/#embeddable-scripting-languages). Be aware that some of them include an interpreter while others only accept pre-compiled byte code.

Just to list a few here:

* [Agora](https://github.com/PuerkitoBio/agora) is a scripting language with Go-like syntax.
* [GopherLua](https://github.com/yuin/gopher-lua) is an interpreter for the [Lua](https://www.lua.org/) scripting language.
* [Otto](https://github.com/robertkrimen/otto) is a JavaScript interpreter.

![Script Plugin](script-plugin.svg)


## Conclusion

The lack of shared, run-time loadable libraries did not stop the Go community from creating and using plugins. There are a number of different approaches to choose from, each one serving particular requirements.

Until Go supports creating and using shared libraries (and rumors about this appear to have been around since Go 1.4


## A simple(-minded) plugin concept

All of the examples listed above have very good documentation and/or examples available. I refrain from repeating any code here and take a bare-bone approach instead, based upon `net/rpc` (and a bit of `os/exec`).

If you are not familiar with RPC in Go, you will want to keep the [documentation](https://golang.org/pkg/net/rpc/) of the `net/rpc` package at hand while reading through the code.

*/

// ### Imports
package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"time"
)

// ### The plugin

// `Plugin` is the type that we register with `net/rpc`. The RPC client can
// then call Plugin's methods.
type Plugin struct {
	listener net.Listener
}

// `Revert` returns the input string reverted.
//
// All methods callable via RPC must satisfy these conditions:
//
// * exported method of exported type
// * two arguments, both of exported or built-in type
// * the second argument is a pointer
// * one return value, of type error
func (p Plugin) Revert(arg string, ret *string) error {
	fmt.Println("Plugin: revert")
	l := len(arg)
	r := make([]byte, l)
	for i := 0; i < l; i++ {
		r[i] = arg[l-1-i]
	}
	*ret = string(r)
	return nil
}

// `Exit` terminates the plugin process.
// `arg` and `ret` are only dummy parameters, to satisfy the
// function signature.
func (p Plugin) Exit(arg int, ret *int) error {
	fmt.Println("Plugin: done.")
	os.Exit(0) // using os.Exit here is suitable for demo code only.
	return nil
}

// `start` starts the plugin server.
func startPlugin() {
	fmt.Println("Plugin start")

	// register the Plugin type with RPC.
	p := &Plugin{}
	err := rpc.Register(p)
	if err != nil {
		log.Fatal("Cannot register plugin: ", err)
	}

	fmt.Println("Plugin: starting listener")
	// Start the listener.
	p.listener, err = net.Listen("tcp", "127.0.0.1:55555")
	if err != nil {
		log.Fatal("Cannot listen: ", err)
	}
	// Start serving requests to the default server.
	fmt.Println("Plugin: accepting requests")
	rpc.Accept(p.listener)
}

// ### The main application

// `app` is our main application. It starts the plugin process
// and calls `Shout()` via RPC.
func app() {
	fmt.Println("App start")
	// Get the plugin.
	p := exec.Command("./plugins", "true")
	p.Stdout = os.Stdout
	p.Stderr = os.Stderr

	// Start the plugin process.
	err := p.Start()
	if err != nil {
		log.Fatal("Cannot start ", p.Path, ": ", err)
	}
	// Ensure the plugin process is up and running before attempting to connect.
	// `net/rpc` has no `DialTimeout` method (unlike `net`).
	time.Sleep(1 * time.Second)

	// Create the RPC client.
	fmt.Println("App: registering RPC client")

	client, err := rpc.Dial("tcp", "127.0.0.1:55555")
	if err != nil {
		log.Fatal("Cannot create RPC client: ", err)
	}

	// Call the Revert method.
	fmt.Println("App: calling Revert")

	var reverse string
	err = client.Call("Plugin.Revert", "Live on time, emit no evil", &reverse)
	if err != nil {
		log.Fatal("Error calling Revert: ", err)
	}
	fmt.Println("App: revert result:", reverse)

	// Stop the plugin and terminate the app.
	fmt.Println("App: stopping the plugin")
	var n int
	client.Call("Plugin.Exit", 0, &n)
	p.Wait()
	fmt.Println("App: done.")
}

// ## main()

func main() {

	// Start either as plugin or as main app.
	if len(os.Args) > 1 && os.Args[1] == "true" {
		// Create the plugin struct and register it with RPC.
		startPlugin()
		// Exit after 10 seconds, in case the main app stopped
		// without terminating the plugin.
		time.AfterFunc(10*time.Second, func() {
			fmt.Println("Plugin: idle timeout - exiting")
			return
		})
	} else {
		// Start the main application.
		app()
	}
}

/*
As always, get the code with `go get -d` to avoid that the binary to show up in your $GOPATH/bin.

	go get -d github.com/appliedgo/plugins
	cd $GOPATH/src/github.com/appliedgo/plugins
	go build
	./plugins

Enjoy!
*/
