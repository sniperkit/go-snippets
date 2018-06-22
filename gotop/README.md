#GoTop

A simple tool to get top like statistics about a running go process.

It uses the output from [expvars](http://golang.org/pkg/expvar/) so will work with any go process that has them exposed.

##Installation
```
go install github.com/gophergala/gotop
```

##Running
```
gotop -url http://golang.org/debug/vars
```

##Future plans
* Add a web interface

