Working Demo
------------

http://where-you-at.com

<img src="https://s3.amazonaws.com/f.cl.ly/items/2i191V0S242c2o352Q1j/Image%202015-01-25%20at%203.26.06%20PM.png" width="300" height="532" title="iOS ScreenShot">

Installing
----------

You need to have Go installed.

```bash
go get -u github.com/gophergala/correct-horse-battery-staple
go get -u -d -tags=js github.com/gophergala/correct-horse-battery-staple/...
```

Running
-------

In the root project folder:

```bash
go build -o main && ./main
```

Deploy
------

```bash
go test ./... && GOOS=linux go build -o main && ./deploy.sh
```

Notes
-----

Features whiteboard: http://whiteboardfox.com/29046-5977-1449

Mockup
------

![](https://s3.amazonaws.com/f.cl.ly/items/0j2B2K2Y3T2u1O1m0h2y/Image%202015-01-23%20at%206.19.12%20PM.png)
