# Hugo

[![Docker Stars](https://img.shields.io/docker/stars/dockette/hugo.svg?style=flat)](https://hub.docker.com/r/dockette/hugo/)
[![Docker Pulls](https://img.shields.io/docker/pulls/dockette/hugo.svg?style=flat)](https://hub.docker.com/r/dockette/hugo/)

Take a look at [official documentation](https://gohugo.io/overview/introduction/).

## Discussion / Help

[![Join the chat](https://img.shields.io/gitter/room/dockette/dockette.svg?style=flat-square)](https://gitter.im/dockette/dockette?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

## Usage

Show help:

```
docker run -it --rm dockette/hugo help new
```

Start embedded server: 

```
docker run -it --rm -v $(pwd):/srv dockette/hugo server -D -v
```

Build site:

```
docker run -it --rm -v $(pwd):/srv dockette/hugo
```
