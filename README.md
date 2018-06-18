# Motivation

This is a toy project to test the go 1.8 plugins feature (which only works on linux at the moment). Work in progress!

The main motivation of the project is to try out [go plugins](https://golang.org/pkg/plugin/) building something that might, some day, become useful (or not, its ok). 

# General Idea

`go-plugserver` is an http server built that registers handlers dynamically, during the runtime, using go plugins. The server keeps polling the `handlers` directory looking for files that look like compiled plugins (`.so` extension), loads those plugins and then run a function that returns two things: 

1. A `string`, the handler pattern url, e.g, `/users/`.
2. A `func(w http.ResponseWriter, r *http.Request))` that is the handler function to be added to the server and will handle requests for that pattern

If a file is removed from the `handlers` directory, the registered pattern for that file will also be removed from the server. **Without any downtime (needs to be confirmed by trusted sources :-))!**

# Using

To register new handlers, compile any go file that has a function:

```go
func Handler() (pattern string, handleFunc func(w http.ResponseWriter, r *http.Request))
```

using `go build -buildmode=plugin -o myplugin.so myplugin.go` and place the generated `.so` file
on the `handlers` directory (check the examples already there). 

Start the server with `make run`.

# TODO

1. Compile the plugins during runtime
2. Use fsnotify instead of polling every 5 seconds (maybe)
3. Investigate a way to enable de-registering patterns without having to implement a ServerMux
4. General support for muxers (depends on 3)

# Disclaimer

On the `plughttp` directory there is a slightly modified version of `golang/go`'s `http.ServerMux` featuring a `DeRegister(pattern string)` method. The LICENSE on that file is the same as `golang/go` and is on the `plughttp` directory.

