package main

import (
	"log"
)

// A Handler responds to an HTTP request.
//
// ServeHTTP should write reply headers and data to the ResponseWriter
// and then return. Returning signals that the request is finished; it
// is not valid to use the ResponseWriter or read from the
// Request.Body after or concurrently with the completion of the
// ServeHTTP call.
//
// Depending on the HTTP client software, HTTP protocol version, and
// any intermediaries between the client and the Go server, it may not
// be possible to read from the Request.Body after writing to the
// ResponseWriter. Cautious handlers should read the Request.Body
// first, and then reply.
//
// Except for reading the body, handlers should not modify the
// provided Request.
//
// If ServeHTTP panics, the server (the caller of ServeHTTP) assumes
// that the effect of the panic was isolated to the active request.
// It recovers the panic, logs a stack trace to the server error log,
// and hangs up the connection. To abort a handler so the client sees
// an interrupted response but the server doesn't log an error, panic
// with the value ErrAbortHandler.
// type Handler interface {
//	ServeHTTP(ResponseWriter, *Request)
//}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
//type HandlerFunc func(ResponseWriter, *Request)

// ServeHTTP calls f(w, r).
//func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
//	f(w, r)
//}

//注意理解这里的设计
// 这里这么做的结果是： 可以将一个函数当做一个接口来使用
// 注意这里的表象是： 让一个函数实现了一个接口，而这个接口的方法和该函数有相同的签名。该函数对该接口的实现只是简单的调用该函数
// 先来看一下没有这种实现的时候我们是怎么做的：
// 1. 定义一个Handler接口来handle http request.
// 2. 我们需要显式的实现该接口（即我们需要一个结构体，或者其他类型）
// 3. 在需要使用该接口的地方使用我们自定义的类型。
// 我们这个实现的问题是什么？-----> 我们定义的结构体可能根本不需要任何成员，而只是为了实现该接口。
// 在采用了标准库的实现以后：我们只需要定义一个和该接口具有相同签名的函数就可以在需要该接口的地方直接使用该函数了。（不需要定义额外的结构体作为载体）。

type Greeter interface {
	Greet(string)
}

type GreeterFunc func(string)

func (g GreeterFunc) Greet(name string) {
	g(name)
}

func SayHello(name string) {
	log.Println("Hello, ", name)
}

func SayYes(name string) {
	log.Println("Yes, ", name)
}

func SayByeBye(name string) {
	log.Println("ByeBye, ", name)
}

func Greeting(g Greeter) {
	g.Greet("foo")
}

func main() {
	Greeting(GreeterFunc(SayHello))
	Greeting(GreeterFunc(SayYes))
	Greeting(GreeterFunc(SayByeBye))
}
