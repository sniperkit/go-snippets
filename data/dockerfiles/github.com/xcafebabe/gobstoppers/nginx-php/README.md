# Nginx + Php 7

A quick and easy way to setup your PHP application in your PI using Docker.

## Usage

### Default App

```
docker run --name nginx-php7 \
 --rm -p 8080:80 \
  kafebob/rpi-nginx-php
```

You can then browse to http://\<docker_host\>:8080 to view the default install files.

### Your app

Simply run the container with a volume mapped to your PHP app hosted in your PI.

```
docker run --name nginx-php7 \
 --rm -p 8080:80 \
 -v $YOUR_APP_FOLDER:/var/www/html
 kafebob/rpi-nginx-php
```

