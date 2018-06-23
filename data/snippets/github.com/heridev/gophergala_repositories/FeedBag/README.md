# Feedbag

Feedbag is a TV first, multiuser Github events dashboard for organizations and teams.

Designed to be flexable, simple, and highly configurable, new widget dashboards can be added by simply adding a template to the `templates/` directory. The config options at the top of the templates define which github event triggers that template to render on the client.

An optional `condition` field can be added to further specify when the template is rendered. ie. `"and .PushEvent (eq .Branch 'master')"` will render the template if there is a force push to "master"

### Version
0.0.1

### Tech

Feedbag uses a number of open source projects to work:

* [Go] - Golang. The best ServerLanguage
* [Gin] - The fastest full-featured web framework for Golang. Crystal clear.
* [AngularJS] - HTML enhanced for web apps!
* [Twitter Bootstrap] - great UI boilerplate for modern web apps
* [Gulp] - the streaming build system
* [jQuery] - duh

## Developer Guide

### Installation

You need Gulp and Bower installed globally:

```sh
$ npm i -g gulp
$ npm i -g bower
```

Optionally, install gin for running the go server

```sh
$ go get github.com/codegangsta/gin
```

```sh
$ git clone [git-repo-url] feedbag
$ cd feedbag
$ godep restore
```

### Startup

Run the go server

```sh
//Point to the location of the gulp index file
$ export INDEX_FILE=web/.tmp/serve/index.html
$ go build
$ ./feedbag
//Or
$ gin
```

Run the gulp task

```sh
$ cd web
$ npm i
$ bower i
$ gulp serve
```

### Todo's

 - Write tests
 - Clean build process
 - Add more events and variables
 - Add more template examples
 - Add styleguide and classes
 - Add user control over repo events shown

License
----

MIT


**Free Software, BooYa!**

[Go]:http://golang.org
[Gin]:http://gin-gonic.github.io/gin/
[Twitter Bootstrap]:http://twitter.github.com/bootstrap/
[jQuery]:http://jquery.com
[AngularJS]:http://angularjs.org
[Gulp]:http://gulpjs.com
