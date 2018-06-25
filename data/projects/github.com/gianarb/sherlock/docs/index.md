`sherlock` is a tool designed to visualize traces. This project starts as PoC to
build a tracing visualization on top of InfluxDB.

Currently [Telegraf a collector with multiple input/output
plugins](https://github.com/influxdata/telegraf) support
[Zipkin](https://github.com/openzipkin/zipkin)
as [input
plugin](https://github.com/influxdata/telegraf/tree/master/plugins/inputs/zipkin). It means that we can instrument our applications using:

* [Zipkin's sdk for different lanagues](https://github.com/openzipkin?utf8=%E2%9C%93&q=&type=&language=)
* [Opentracing sdk](https://github.com/opentracing-contrib)

To send traces to Telegraf. It can store them in InfluxDB or all the other
supported [output
plugin](https://github.com/influxdata/telegraf/tree/master/plugins/outputs). But
at the moment we will stay focused on InfluxDB.

`sherlock` is an application capable of reading and visualize traces.

The [`docs`](/docs) folder contains all the information you need to know.

We are building this project because this is the missing steps to have a fully
integrated experience of tracing using InfluxDB and Telegraf. This is the
missing part for now. We know how to store the traces and we would like to have
a way to read them out!

## Table of content
1. [Tracing, zipkin, opentracing...](/docs/tracing.md): what is it?
2. [Telegraf, InfluxDB](/docs/tick_stack.md): links and docs about these projects
and why they are good for tracing
3. [Architecture](/docs/architecture.md): all about how this project is made and
why.

Have fun!

[** Next > What tracing is**](/docs/tracing.md)
