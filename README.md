# Serverless frameworks
A curated list of open source serverless frameworks by popularity with a short description from the repositories.

1. [Serverless](https://github.com/kaxap/serverless_list/#1-serverless)
2. [OpenFaaS](https://github.com/kaxap/serverless_list/#2-openfaas)
3. [Fission](https://github.com/kaxap/serverless_list/#3-fission)
4. [Apache OpenWhisk](https://github.com/kaxap/serverless_list/#4-apache-openwhisk)
5. [Gordon](https://github.com/kaxap/serverless_list/#5-gordon)
6. [Iron.io](https://github.com/kaxap/serverless_list/#6-ironfunctions-ironio)
7. [fnproject.io](https://github.com/kaxap/serverless_list/#7-fn-fnprojectio)
8. [Kubeless](https://github.com/kaxap/serverless_list/#8-kubeless)
9. [nuclio](https://github.com/kaxap/serverless_list/#9-nuclio)
10. [Microcule](https://github.com/kaxap/serverless_list/#10-microcule)
11. [Rainbound](https://github.com/kaxap/serverless_list/#11-rainbond)
12. [.architect](https://github.com/kaxap/serverless_list/#12-architect)



## 1. [Serverless](#serverless)
[The Serverless Framework](https://github.com/serverless/serverless) – Build applications comprised of microservices that run in response to events, auto-scale for you, and only charge you when they run. This lowers the total cost of maintaining your apps, enabling you to build more logic, faster.

The Framework uses new event-driven compute services, like AWS Lambda, Google CloudFunctions, and more. It's a command-line tool, providing scaffolding, workflow automation and best practices for developing and deploying your serverless architecture. It's also completely extensible via plugins.

Serverless is an MIT open-source project, actively maintained by a full-time, venture-backed team.

#### Languages
* Node.js
* Python
```python
mport logging
import requests

log = logging.getLogger()
log.setLevel(logging.DEBUG)

API_HOST = 'https://jsonplaceholder.typicode.com'


def list_posts(event, context):
    """ Retrieve posts list """
    url = API_HOST + '/posts'

    log.debug('calling ' + url)
    res = requests.get(url)

    response = {
        "statusCode": res.status_code,
        "body": res.text
    }

    return response


def get_post(event, context):
    """ Retrieve post by id """
    url = API_HOST + '/posts/' + str(event['pathParameters']['id'])

    log.debug('calling ' + url)
    res = requests.get(url)

    response = {
        "statusCode": res.status_code,
        "body": res.text
    }

    return response
```
* Java
* Scala
* C#
* F#
* Groovy
* Kotlin
* PHP 
* Swift

## 2. [OpenFaaS](#openfaas)
[OpenFaaS](https://github.com/openfaas/faas) (Functions as a Service) is a framework for building serverless functions with Docker which has first class support for metrics. Any process can be packaged as a function enabling you to consume a range of web events without repetitive boiler-plate coding.

#### Supported languages/environments

 * Javascript
```javascript
"use strict"

module.exports = (callback, context) => {
    callback(null, {"message": "You said: " + context})
}
```
 * Python
```python
import requests

def handle(req):
        r =  requests.get(req, timeout = 1)
        print(req +" => " + str(r.status_code))
```

## 3. [Fission](#fission)
[Fission](https://github.com/fission/fission) is a fast serverless framework for Kubernetes with a focus on developer productivity and high performance.

Fission operates on just the code: Docker and Kubernetes are abstracted away under normal operation, though you can use both to extend Fission if you want to.

Fission is extensible to any language; the core is written in Go, and language-specific parts are isolated in something called environments (more below). Fission currently supports NodeJS, Python, Ruby, Go, PHP, Bash, and any Linux executable, with more languages coming soon.

#### Supported languages/environments
* Go	
```go
package main

import (
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	msg := "Hello, World!"
	w.Write([]byte(msg))
}
```
* .NET
* NodeJS
```javascript
module.exports = async function(context) {
    return {
        status: 200,
        body: "Hello, world!\n"
    };
}
```
* Perl
* PHP 7
* Python 3
```python
def main():
    return "Hello, world!\n"
```
* Ruby
* Running binaries/scripts

## 4. [Apache OpenWhisk](#openwhisk)
[OpenWhisk](https://github.com/apache/incubator-openwhisk) is a cloud-first distributed event-based programming service. It provides a programming model to upload event handlers to a cloud service, and register the handlers to respond to various events.

## 5. [Gordon](#gordon)
[Gordon](https://github.com/jorgebastida/gordon) is a tool to create, wire and deploy AWS Lambdas using CloudFormation.
#### Supported languages/environments

* Python
```python
print('Loading function')


def handler(event, context):
    # print("Received event: " + json.dumps(event, indent=2))
    print("value1 = " + event['key1'])
    print("value2 = " + event['key2'])
    print("value3 = " + event['key3'])
    return event['key1']  # Echo back the first key value
    # raise Exception('Something went wrong')
```
* Javascript
* Java
* Golang
```go
package main

import "fmt"

func main() {
	fmt.Println("Hello, 世界")
}
```
* Scala

## 6. [IronFunctions (Iron.io)](#ironfunctions)
[IronFunctions](https://github.com/iron-io/functions) is an open source serverless platform, or as we like to refer to it, Functions as a Service (FaaS) platform that you can run anywhere.

#### Supported languages/environments

* Go
```go
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Person struct {
	Name string
}

func main() {
	p := &Person{Name: "World"}
	json.NewDecoder(os.Stdin).Decode(p)
	fmt.Printf("Hello %v!", p.Name)
}
```
* Anything in a container

## 7. [Fn (fnproject.io)](#fn)
[Fn](https://github.com/fnproject/fn) is an event-driven, open source, functions-as-a-service compute platform that you can run anywhere.

#### Supported languages/environments

* Go
```go
package main

import (
  "fmt"
)

func main() {
  fmt.Println("Hello from Fn!")
}
```
* Anything in a container

## 8. [Kubeless](#kubeless)
[Kubeless](https://github.com/kubeless/kubeless) is a Kubernetes-native serverless framework that lets you deploy small bits of code without having to worry about the underlying infrastructure plumbing. It leverages Kubernetes resources to provide auto-scaling, API routing, monitoring, troubleshooting and more.

Kubeless stands out as we use a Custom Resource Definition to be able to create functions as custom kubernetes resources. We then run an in-cluster controller that watches these custom resources and launches runtimes on-demand. The controller dynamically injects the functions code into the runtimes and make them available over HTTP or via a PubSub mechanism.

#### Supported languages/environments
* Python
```python
def foo():
    return "hello world"
```
* dotnetcore
```csharp
using System;
using Microsoft.AspNetCore.Http;

public class helloget
{
    public string foo(HttpRequest request)
    {
        return "hello world";
    }
}
```
* Ruby
```ruby
def foo(context)
  "hello world"
end
```

## 9. [nuclio](#nuclio)
[nuclio](https://github.com/nuclio/nuclio) is a new serverless project, derived from iguazio's elastic data life-cycle management service for high-performance events and data processing. nuclio is being extended to support a large variety of event and data sources. You can use nuclio as a standalone binary (for example, for IoT devices), package it within a Docker container, or integrate it with a container orchestrator like Kubernetes.

nuclio is extremely fast. A single function instance can process hundreds of thousands of HTTP requests or data records per second. This is 10–100 times faster than some other frameworks.
#### Supported languages/environments

* Python
```python
def handler(context, event):
    response_body = f'Got {event.method} to {event.path} with "{event.body}"'

    # log with debug severity
    context.logger.debug('This is a debug level message')

    # just return a response instance
    return context.Response(body=response_body,
                            headers=None,
                            content_type='text/plain',
                            status_code=201)
```
* Go
```go
package handler

import (
    "github.com/nuclio/nuclio-sdk"
)

func Handler(context *nuclio.Context, event nuclio.Event) (interface{}, error) {
    context.Logger.Info("Request received: %s", event.GetPath())

    return nuclio.Response{
        StatusCode:  200,
        ContentType: "application/text",
        Body: []byte("Response from handler"),
    }, nil
}
```

## 10. [Microcule](#micromule)
[Microcule](https://github.com/Stackvana/microcule) is a software Development Kit and Command Line Interface for spawning streaming stateless HTTP microservices for any programming language or arbitrary binary.

Think of it as serverless functions meets [Unix Philosophy](https://en.wikipedia.org/wiki/Unix_philosophy).

#### Supported languages/environments

   - C ( with `gcc` )
   ```c
#include<stdio.h>

int main(void)
{ // This is a comments
    printf("Hello world!\n");
    return 0;
}
```
   - Java
```java
class hook
{
        public static void main(String args[])
        {
           System.out.println("Hello JavaWorld!!!!");
        }
}
```
   - Javascript ( first-class support )
```javascript
module.exports = function (req, res) {
  res.write('Hello, this is a JavaScript function.\n');
  res.write('hook.params is populated with request parameters\n');
  res.write(JSON.stringify(req.params, true, 2));
  res.end('');
};
```
   - Babel ( ES6 / ES7 / etc ... )
   - Coffee-script
   - Common Lisp
   - Bash
   - Lua
   - Golang
```go
package main
import "fmt"
func main() {
    fmt.Println("hello world")
}
```
   - Ocaml
   - Perl
   - PHP
   - Python
```python
import pprint
print "hello world"
```
   - Python3
   - Ruby
   - Rust
```rust
// This is the main function
fn main() {
    // The statements here will be executed when the compiled binary is called

    // Print text to the console
    println!("Hello RustWorld!");
}
```
   - R
   - Scheme
   - Smalltalk
   - tcl

## 11. [Rainbond](https://github.com/goodrain/rainbond)
Rainbond is a production grade serverless PaaS. It integrates Kubernetes based containers management, CI/CD and multiple data center resource management, to provide full management of cloud native applications, build connections between application and infrastructure, application and application, infrastructure and infrastructure.

#### Supported languages/environments
- Java
- PHP 
- Static (Nginx, Apache)
- Python 2,3
- NodeJS
- Golang
- Ruby

## 12. [.architect](https://arc.codes/)
 - Version control your architecture and create cloud infra in minutes from an .arc manifest
 - Deploy in seconds with first class support for staging and production
 - Work locally while completely offline with a speedy in-memory database
 - Primitives not Frameworks; define app architecture agnostic of vendor arcana
 
#### Supported languages/environments
 - NodeJS
