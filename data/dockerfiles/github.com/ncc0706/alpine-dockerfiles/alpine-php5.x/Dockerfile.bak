
FROM alpine:latest

MAINTAINER niuyuxian <ncc0706@gmail.com>

ENV TIMEZONE Asia/Shanghai

#root用户密码
ENV ROOT_PASSWORD=123456

#nginx版本
ENV NGINX_VER=1.12.0
#redis版本
ENV REDIS_VER=3.2.9
#redis密码
ENV REDIS_PASS=123456

ENV PHP_VERSION 5.6.31
ENV PHP_URL="https://secure.php.net/get/php-5.6.31.tar.xz/from/this/mirror"
ENV PHP_ASC_URL="https://secure.php.net/get/php-5.6.31.tar.xz.asc/from/this/mirror"
ENV PHP_SHA256="c464af61240a9b7729fabe0314cdbdd5a000a4f0c9bd201f89f8628732fe4ae4"
ENV PHP_MD5=""

## 国内时区
RUN apk add --update tzdata \
    && cp /usr/share/zoneinfo/${TIMEZONE} /etc/localtime \
    && echo "${TIMEZONE}" > /etc/timezone

ENV PHP_INI_DIR /usr/local/etc/php

RUN

RUN mkdir -p /usr/src; \
    cd /usr/src; \
    wget -O php.tar.xz "$PHP_URL"; \
    mkdir -p $PHP_INI_DIR/conf.d



CMD ["php", "-a"]

