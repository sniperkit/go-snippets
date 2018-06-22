# Nginx

Nginx with SSL / HSTS and HTTP2 module. 

-----

[![Docker Stars](https://img.shields.io/docker/stars/dockette/nginx.svg?style=flat)](https://hub.docker.com/r/dockette/nginx/)
[![Docker Pulls](https://img.shields.io/docker/pulls/dockette/nginx.svg?style=flat)](https://hub.docker.com/r/dockette/nginx/)

## Discussion / Help

[![Join the chat](https://img.shields.io/gitter/room/dockette/dockette.svg?style=flat-square)](https://gitter.im/dockette/dockette?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

## Default

Without any volumes Nginx listen IPv4 and IPv6 on both ports 80/443 and serve 444 No content.

## Usage

```sh
docker run \
    -p 80:80 \ 
    -p 443:433 \
    -v /path/to/site:/etc/nginx/sites.d/site \
    --name nginx \
    dockette/nginx:latest
```

## Configuration

There are 2 main files:

- [nginx.conf](https://github.com/dockette/nginx/blob/master/nginx.conf)
- [mime.types](https://github.com/dockette/nginx/blob/master/mime.types)

Diffie-Hellman files:

- [dhparam2048.pem](https://github.com/dockette/nginx/blob/master/dhparam2048.pem)
- [dhparam4096](https://github.com/dockette/nginx/blob/master/dhparam4096.pem)

Modules:

- nginx-module-xslt
- nginx-module-geoip
- nginx-module-image-filter
- nginx-module-njs

### Http config

There are many files in `/etc/nginx/conf.d`.

- [basic.conf](https://github.com/dockette/nginx/blob/master/conf.d/basic.conf)
- [charset.conf](https://github.com/dockette/nginx/blob/master/conf.d/charset.conf)
- [gzip.conf](https://github.com/dockette/nginx/blob/master/conf.d/gzip.conf)
- [logging.conf](https://github.com/dockette/nginx/blob/master/conf.d/logging.conf)
- [ssl.conf](https://github.com/dockette/nginx/blob/master/conf.d/ssl.conf)

These files are loaded manully in the main `http` section. You can extend only part of these files which you need.

Please do not place your custom configurations to this folder. There is a folder `/etc/nginx/user.conf.d` for this reason.

Take a look at [nginx.conf](https://github.com/dockette/nginx/blob/master/nginx.conf).

### Sites

Sites are loaded from folder `/etc/nginx/sites.d`.

There are 2 sites predefined:

The [default](https://github.com/dockette/nginx/blob/master/sites.d/default) returns http code 444 for every requests, it's marked as `default_server`.

And the [example](https://github.com/dockette/nginx/blob/master/.examples) site.

### Sites config

For little fragments shared accross sites you can use folder `/etc/nginx/site.conf.d`.

### User config

Your custom global http section configuration, `/etc/nginx/user.conf.d`. These files are loaded automatic from `http` section.

Take a look at [nginx.conf](https://github.com/dockette/nginx/blob/master/nginx.conf).

### SSL

By default nginx uses 2048 bits dhparam cypher. There is also 4096 bits pregenerated file. 

This approach is for developing, for your stagging / production, please generate your own dhparam.

## Tip

You should have A+ on [SSL Server Test](https://www.ssllabs.com/ssltest/).

![](https://raw.githubusercontent.com/dockette/nginx/master/.docs/ssllabs.png "SSL Server Test")
