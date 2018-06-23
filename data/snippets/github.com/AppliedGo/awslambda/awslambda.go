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
title = "In the news: Go on AWS Lambda"
description = "AWS Lambda announced support for Go. A summary."
author = "Christoph Berger"
email = "chris@appliedgo.net"
date = "2018-01-28"
draft = "false"
domains = ["DevOps"]
tags = ["Lambda", "FaaS", "AWS"]
categories = ["News"]
+++

Just recently, Amazon announced support for Go on AWS Lambda. Here is a summary of last week's news around this topic.

<!--more-->

On Jan 15th, Amazon [announced](https://aws.amazon.com/blogs/compute/announcing-go-support-for-aws-lambda/) Go support for AWS Lambda. This was exciting news for many, according to the number of blog posts that followed this announcement. To save you a bit of searching, here are some of the posts from the last week.

## But wait, what is AWS Lambda?

Fair point, not everyone is familiar with the idea of "serverless" in general and AWS Lambda in particular.

In a nutshell, with AWS Lambda you can upload *a single function* to Amazon's servers and have the Lambda environment invoke this function when particular events occur. No need for setting up and maintaining an HTTP server, or any server at all - all you need to do is writing your Lambda function. This is often called "Function as a Service" (or FaaS), and if you prefer an open-source, self-hostable[^1] alternative, have a look at [OpenFaaS](https://github.com/openfaas/faas) or the [fn project](http://fnproject.io/).
The Events that trigger a Lambda function can come from various sources, mostly from within the AWS ecosystem, but also from the outside via a REST API gateway. Your function then can process the event and respond by accessing other services (AWS services as well as non-AWS services).

Before starting to write your own Lambda function, you will want to have a look at the [AWS Lambda documentation](https://docs.aws.amazon.com/lambda/latest/dg/welcome.html).

[^1]: Does the word "hostable" even exist?


But now on to the news review.

## Speed and Stability: Why Go is a Great Fit for Lambda

Let me start with a [motivational article](https://brandur.org/go-lambda) by [@brandur](https://twitter.com/brandur).

In this article, Brandur observes that while AWS Lambda supports quite a few languages, it seems difficult for them to keep the runtime environments of these languages up to date. Most of them lag behind by quite a few releases.

Go, on the other hand, comes with two properties that seem a perfect fit for Lambda.

First, there is no runtime. A pure Go binary can be statically linked and needs no pre-installed runtime of any kind. AWS Lambda takes advantage of this by allowing to upload Go code in compiled form. No need for a compilcated dependency deployment system because all dependencies are already baked into the binary.

Second, Go 1.x has an extraordinarily stable language specification, paired with an equally stable standard library. (This is a direct result from the [Go 1 Compatibility Promise](https://golang.org/doc/go1compat).) AWS Lambda's Go support thus does not necessarily need an update when a new minor Go version is released.

Brandur concludes that,

> Along with the languages normal strengths – incredible runtime speed, an amazing concurrency story, a great batteries-included standard library, and the fastest edit-compile-debug loop in the business – Go’s stability and ease of deployment is going to make it a tremendous addition to the Lambda platform. I’d even go so far as to say that you might want to consider not writing another serverless function in anything else.


## Go Based AWS Lambda Tutorial

After reading Brandur's article, you might be thinking, "wow, that's fantastic! Why am I still sitting here reading blogs? Get me started quiiiiick!!"

Noting easier than that. In less than 10 minutes, Elliott from TutorialEdge walks you through the steps of creating, uploading, and running a minimal Lambda function, in [this video](https://www.youtube.com/watch?v=x_yCX4kSchY).


## Test Lambda functions locally (1): mtojek/aws-lambda-go-proxy

Once you start going beyond simple "Hello, world"-style Lambda functions, you might wish to speed up the develop-deploy-test cycles. In this case, [aws-lambda-go-proxy](https://github.com/mtojek/aws-lambda-go-proxy) is for you. This tool uses a cool trick to achieve the speedup: It removes the "deploy" part from "develop-deploy-test", by establishing a proxy function that routes all Lambda events to your local machine. A nice aspect is that the code needs no modifications in order to run behind the proxy function.


## Test Lambda functions locally (2): djhworld/go-lambda-invoke

While `aws-go-lambda-proxy` employs a proxy function, `go-lambda-invoke` goes a different route and directly redirects the TCP traffic between the AWS Lambda service and the Lambda function to the local machine. Another difference to the aforementioned package: The code needs to explicitly start the lambda function in order to have it listen to AWS events.


## Test Lambda functions locally (3): SAM Local (Beta)

Amazon provides a [Beta version of SAM Local](https://docs.aws.amazon.com/lambda/latest/dg/test-sam-local.html) for testing Lambda functions in an entirely local environment. Support for Go was added just four days ago, on Jan 24th.


## Turn HTTP handlers into Lambda functions: apex/gateway and akrylysov/algnhsa

So you have written that pretty cool HTTP handler, and now you want to run it as a Lambda function? With [apex/gateway](https://github.com/apex/gateway), you simply swap out http.ListenAndServe for gateway.ListenAndServe - it can't get any easier.

Package [`akrylysov/algnhsa`](https://github.com/akrylysov/algnhsa) basically does the same, and comes with a [blog post](http://artem.krylysov.com/blog/2018/01/18/porting-go-web-applications-to-aws-lambda/) that discusses the steps required to port an HTTP handler to Lambda without using a convenience package like `algnhsa` or `gateway`.


## Conclusion

AWS Lambda seems to be quite popular among Gophers, given the number of articles and packages that appeared only days after support for Go has been officially announced. I hope this quick summary helps you getting started quickly with writing your own Lambda functions in Go.

**Happy coding!**

*/
