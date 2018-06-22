# LetsEncrypt

[![Docker Stars](https://img.shields.io/docker/stars/dockette/letsencrypt.svg?style=flat)](https://hub.docker.com/r/dockette/letsencrypt/)
[![Docker Pulls](https://img.shields.io/docker/pulls/dockette/letsencrypt.svg?style=flat)](https://hub.docker.com/r/dockette/letsencrypt/)

Create 90 days SSL certificates for given domains.

## How it works

Container creates simple Nginx server listening on port 80 and waiting for letsencrypt validation.

It handles only requests to `mydomain.com/.well-known`, all other requests are forbidden. 

```
server {
    listen 80;
    server_name $DOMAINS;

    location ^~ /.well-known/ {
        root /var/www/acme-certs;
    }

    location / {
        return 403;
    }
}
```

## Usage

```sh
docker run \
    -p 80:80 \ 
    -v /srv/certs/mydomain.com:/var/www/certs \
    --name le \
    -e DOMAINS='mydomain.com www.mydomain.com' \
    -e EMAIL='my@email.tld' \
    dockette/letsencrypt:latest
```

You can add `-it` for interactive shell.

After that you will have copies of certificates in your `/srv/certs/mydomain.com/` folder.

