# Golang Advanced Snippets
Go Advanced Snippet records during learning Go ...

## Concurrency
* atomic
* channel
* goroutine
* mutex
* select
* pool
* runner
* worker

## go2cache
Concurrency-safe golang caching library with expiration capabilities.

### Build
Make sure you have a working Go environment. See the [install instructions](http://golang.org/doc/install.html).

To install go2cache, simply run:
```sh
    go get github.com/Jiangew/hancock
```

To compile it from source:
```sh
    cd $GOPATH/src/github.com/Jiangew/hancock/go2cache
    go get -u -v
    go build && go test -v
```

### Test
Also see our test-cases in cache_test.go for further working examples.

## Crawler
Golang implement crawler framework.

## Snowflake
Golang implement snowflake algorithm.

## Consul
* Server
* Client
* Distributed Locks
    * Distributed locks with consul and golang.

## NSQ
* producer
* consumer

## Rank List
Golang implement rank list based on redis and consul.

## Short Url
Golang implement short url service based on beego web framework.

## Elasticsearch Client
* elasticsearch client based on v5.x
* bulk insert
* bulk processor
* sliced scroll

## Gin
Gin is a web framework written in Go.
It features a martini-like API with much better performance, up to 40 times faster thanks to httprouter. 


## Contents
- [Go 高并发 concurrency 使用总结](http://www.grdtechs.com/2016/02/17/go-concurrency-summarize)
- [Go 缓存库 cache2go 介绍](http://time-track.cn/cache2go-introduction.html)
- [Go 实现 snowflake 算法：分布式唯一 id 生成器](http://studygolang.com/articles/9753)
- [Go 使用 consul 做服务发现](http://changjixiong.com/use-consul-in-golang)
- [An Example of Using NSQ From Go](http://tleyden.github.io/blog/2014/11/12/an-example-of-using-nsq-from-go)
- [Go 使用实时消息中间件 NSQ](http://changjixiong.com/golang%E5%AE%9E%E6%97%B6%E6%B6%88%E6%81%AF%E5%B9%B3%E5%8F%B0nsq%E7%9A%84%E4%BD%BF%E7%94%A8)
- [Go 实现基于 Redis and Consul 的排行榜服务](http://changjixiong.com/a-scalable-rank-server-example-base-on-redis-and-consul-in-golang)
- [Distributed Locks with Consul and Golang](https://medium.com/m/global-identity?redirectUrl=https://distributedbydefault.com/distributed-locks-with-consul-and-golang-c4eccc217dd5)
- [Elasticsearch Client for Go](https://olivere.github.io/elastic/)
- [Elasticsearch Client for Go wiki](https://github.com/olivere/elastic/wiki)
- [Gin - A high performance HTTP web framework](https://github.com/gin-gonic/gin)
- [Testing Gin with JSON response](https://medium.com/@craigchilds94/testing-gin-json-responses-1f258ce3b0b1)
- [A high performance HTTP request router](https://github.com/julienschmidt/httprouter)
- [Go Channel 应用模式](http://colobu.com/2018/03/26/channel-patterns/#eapache)
- [Go sync:mutex.TryLock vs Channel](https://github.com/golang/go/issues/6123)
- [Concurrency in Go](https://github.com/kat-co/concurrency-in-go-src)
