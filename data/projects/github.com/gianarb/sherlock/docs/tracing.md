`What is tracing?`  this is gonna be a quick intro to the tracing topic,
probably with links and so on.

As a developer, you probably discover tracing in different occasion. When you
profile your applications, you are building a trace for all the methods that you
are calling inside your application.

On Safari, Chrome, Firefox you can inspect `Network` to discover how many
requests and how many files a particular page downloads. That's tracing.

This concept applies well to a distributed system because just as it is
for multi-threads applications when a call jumps between threads in a
    microservices environment a particular request jump over the network between
    applications.

    You need to answer questions like:
    * Where my requests failed and why?
    * Who is the slower service for a particular request?
    * Why latency between these two services is growing? Is it?
    * and a lot more

    Traces helps...

## [Opentracing](http://opentracing.io/)
As you can imagine you need to instrument your applications to build and send
traces. It means that if you have a single page application in Javascript, two
mobile apps in Java and C++ plus a backend in Go you need to write sdks for all
of them. And probably it already sounds like a waste of time. Plus there are a
lot of tracers (applications to store and manage traces). So what about
vendor-lock-in?
Opentracing is a standard developed and promoted by the Cloud Native Foundation
to solve these problems and many more.

## Tracers
There are a lot of tracers. [Zipkin](https://github.com/openzipkin/zipkin) and
[Jeager](https://www.jaegertracing.io/) are famous in the open source world; the
first one is in Java from Twitter the second in Go from Uber.

AWS ships an as a service tracer called
[X-Ray](https://aws.amazon.com/it/xray/), Google has its own called
[StackDriver](https://cloud.google.com/trace/).

All the quoted tracers and many more are in some way compatible with
Opentracing. It means that you can switch between them transparently. Or at
least that's the goal.

[** Next > TICK Stack?? **](/docs/tick_stack.md)
