# h5ai

[h5ai](http://larsjung.de/h5ai/) is a modern web server index.
This [docker](https://www.docker.io/) image makes it trivially easy to
spin up a webserver and start sharing your files through the web.

## Usage

This docker image is available as a [trusted build on the docker index](https://hub.docker.com/r/bixidock/h5ai/),
so there's no setup required.
Using this image for the first time will start a download automatically.
Further runs will be immediate, as the image will be cached locally.

The recommended way to run this container looks like this:

```bash
$ sudo docker run -it --rm -p 8080:80 -v /your/data:/var/www bixidock/h5ai
```

You can also use docker-compose.yml to run this container


```yaml
version: '2'

services:

  h5ai:
    image: bixidock/h5ai
    container_name: h5ai
    restart: always
    volumes:
      - "/your/data:/var/www"
```

You can now point your webbrowser to this URL:

```
http://localhost/
```
