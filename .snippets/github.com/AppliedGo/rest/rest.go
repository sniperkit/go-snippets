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
title = "Take a REST!"
description = "The basics of a RESTful Web API, with a tiny REST server in Go."
author = "Christoph Berger"
email = "chris@appliedgo.net"
date = "2016-09-14"
publishdate = "2016-09-14"
domains = ["Internet and Web"]
tags = ["REST", "Web", "API", "Video"]
categories = ["Tutorial"]
+++

RESTful Web API's are ubiquitous. Time for a minimalistic, five-minutes video tutorial about REST, RESTful API's, and buidling a REST server in Go.

<!--more-->

{{< youtube iVXaPD_Jbu0 >}}

- - -
*This is the transcript of the video.*
- - -

Hello and welcome to a little crash course about RESTful Web API Basics in Go. I am Christoph Berger, and in this video we will look into REST basics and how to build a tiny REST server in Go.

This course consists of two parts. First, we’ll have a brief look at the basics of RESTful Web API’s. After that, we’ll build a tiny REST server together.


## Rest in a Nutshell

Before touching any code, let’s have a closer look at the concepts behind REST.


### What is Rest?

REST is an architectural paradigm that can be seen as an abstraction of the structure and behavior of the World Wide Web. Just like a human clicks a hyperlink and gets a new page of data back, a REST client sends a request to a server and gets a response back.

![REST Client And Server](RESTClientAndServer.png)


### Basic Operations

The REST protocol is based on four basic operations: Create, Read, Update, and Delete. These operations are often described by the acronym CRUD.

![CRUD](CRUD.png)


### REST and HTTP

REST is often but not always based on the Hypertext Transfer Protocol. Methods are built from HTTP verbs like GET, POST, PUT, DELETE, et cetera.

![HTTP Verbs](HTTPVerbs.png)

URL’s just point to resources that the methods can operate on. A URL never includes a method name.

![REST URL's](RESTURLs.png)

This is how the CRUD operations can be mapped onto HTTP verbs:

* Create is mapped to HTTP POST.
* Read is mapped to GET,
* Update to PUT, and
* Delete to DELETE.

This assignment is not absolute. For example, both PUT and POST could be used to create or update data.

In RESTful Web API’s, a URL can describe a collection of elements, or an individual element.
An HTTP verb may therefore represent two different operations, depending on whether the resource is a collection or a single element.


### How To Pass Data

For sending data to the server, there are two options.

* First, you can send small amounts of data within the URL itself.
* Second, data can reside in the body of the HTTP request.

The server always returns data via the body of the HTTP response.


### REST and SQL

Coincidentally, database operations also follow the CRUD scheme:

* Create maps to SQL INSERT,
* Read maps to SELECT,
* Update to UPDATE, and
* Delete to DELETE.

Well, ok, this is perhaps not much of a coincidence at all. But anyway, this is a straightforward way of making database operations accessible to Web clients.


## A Tiny REST Server In Go

Our code consists of standard Go, except for the HTTP router. The standard ServerMux router is very simplistic and provides no path variables nor any complex pattern matching. A better choice is `httprouter`.  `httprouter` provides path variables, designated by a leading colon, as well as a simple way of mapping HTTP verbs to CRUD operations.

![httprouter features](httprouter.png)

- - -
**UPDATE:** The code has been kept simple for clarity. The original version as seen in the video does not even check for concurrent access to the data store. This is no problem when testing the code by sending `curl` calls one-by-one, but in real-world applications this can mess up your data. Hence the below code uses [sync.Mutex](https://golang.org/pkg/sync/#Mutex) to guard access to the global data store.
- - -

*/

// ## Imports and globals
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	// This is `httprouter`. Ensure to install it first via `go get`.
	"github.com/julienschmidt/httprouter"
)

// We need a data store. For our purposes, a simple map
// from string to string is completely sufficient.
type store struct {
	data map[string]string

	// Handlers run concurrently, and maps are not thread-safe.
	// This mutex is used to ensure that only one goroutine can update `data`.
	m sync.RWMutex
}

var (
	// We need a flag for setting the listening address.
	// We set the default to port 8080, which is a common HTTP port
	// for servers with local-only access.
	addr = flag.String("addr", ":8080", "http service address")

	// Now we create the data store.
	s = store{
		data: map[string]string{},
		m:    sync.RWMutex{},
	}
)

// ## main
func main() {
	// The main function starts by parsing the commandline.
	flag.Parse()

	// Now we can create a new `httprouter` instance...
	r := httprouter.New()

	// ...and add some routes.
	// `httprouter` provides functions named after HTTP verbs.
	// So to create a route for HTTP GET, we simply need to call the `GET` function
	// and pass a route and a handler function.
	// The first route is `/entry` followed by a key variable denoted by a leading colon.
	// The handler function is set to `show`.
	r.GET("/entry/:key", show)

	// We do the same for `/list`. Note that we use the same handler function here;
	// we'll switch functionality within the `show` function based on the existence
	// of a key variable.
	r.GET("/list", show)

	// For updating, we need a PUT operation. We want to pass a key and a value to the URL,
	// so we add two variables to the path. The handler function for this PUT operation
	// is `update`.
	r.PUT("/entry/:key/:value", update)

	// Finally, we just have to start the http Server. We pass the listening address
	// as well as our router instance.
	err := http.ListenAndServe(*addr, r)

	// For this demo, let's keep error handling simple.
	// `log.Fatal` prints out an error message and exits the process.
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// ## The handler functions

// Let's implement the show function now. Typically, handler functions receive two parameters:
//
// * A Response Writer, and
// * a Request object.
//
// `httprouter` handlers receive a third parameter of type `Params`.
// This way, the handler function can access the key and value variables
// that have been extracted from the incoming URL.
func show(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// To access these parameters, we call the `ByName` method, passing the variable name that we chose when defining the route in `main`.
	k := p.ByName("key")

	// The show function serves two purposes.
	// If there is no key in the URL, it lists all entries of the data map.
	if k == "" {
		// Lock the store for reading.
		s.m.RLock()
		fmt.Fprintf(w, "Read list: %v", s.data)
		s.m.RUnlock()
		return
	}

	// If a key is given, the show function returns the corresponding value.
	// It does so by simply printing to the ResponseWriter parameter, which
	// is sufficient for our purposes.
	s.m.RLock()
	fmt.Fprintf(w, "Read entry: s.data[%s] = %s", k, s.data[k])
	s.m.RUnlock()
}

// The update function has the same signature as the show function.
func update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// Fetch key and value from the URL parameters.
	k := p.ByName("key")
	v := p.ByName("value")

	// We just need to either add or update the entry in the data map.
	s.m.Lock()
	s.data[k] = v
	s.m.Unlock()

	// Finally, we print the result to the ResponseWriter.
	fmt.Fprintf(w, "Updated: s.data[%s] = %s", k, v)
}

/*
After saving, we can run the code locally by calling

```
go run rest.go
```

Now we can call our server. For this, let's use curl. Curl is an HTTP client for the command line. By default, it sends GET requests, but the X parameter lets us create a PUT request instead.

First, let's fill the map with some entries. We do that by sending PUT requests with a key and a value.

Then we can request a list of all entries, as well as individual entries by name.

```
curl -X PUT localhost:8080/entry/first/hello
curl -X PUT localhost:8080/entry/second/hi
curl localhost:8080/list
curl localhost:8080/entry/first
curl localhost:8080/entry/second
```

As always, the code (with all comments) is available on GitHub: https://github.com/appliedgo/rest

(No Playground link this time, as the Go Playground does not allow running Web servers.)


That’s it for today, thanks for reading, and happy coding!

---

Errata: 2017-03-05 Fixed: mutex logic for the `data` map.
*/
