# Nginx

## Examples

Base Container, HTTP

```
sudo docker run --rm -p 80:80 kafebob/rpi-nginx
```

Serve your content, HTTP

```
sudo docker run --rm \
  -p 80:80 \
  -v $YOUR_CONTENT_FOLDER:/var/www \
  kafebob/rpi-nginx
```

Serve your content, HTTP(s). Will generate an insecure HTTPS certificate and enable Diffie-Hellman for TLS.

```
sudo docker run --rm \
  -p 80:80 \
  -p 443:443 \
  -e NGINX_ENABLE_SSL=1 \
  -v $YOUR_CONTENT_FOLDER:/var/www \
  kafebob/rpi-nginx
```
