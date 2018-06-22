# alpine-bigchaindb

[![Build Status](https://travis-ci.org/artus/alpine-bigchaindb.svg?branch=master)](https://travis-ci.org/artus/alpine-bigchaindb)

## A small docker image for bigchaindb running on alpine.

> This project is still in alpha, please test extensively before using it 'for real'.

## Why alpine linux?

![filesize difference between running on ubuntu and alpine](https://artus.github.io/alpine-bigchaindb/img/filesize.png)

## How to run

You can run the image by using the same steps listed on [the documentation of bigchaindb.](https://docs.bigchaindb.com/projects/server/en/latest/appendices/run-with-docker.html)

```shell
docker run -d  \
           --name=bigchaindb \
           -p 59984:9984 \
           --restart=always \ 
           -v $HOME/bigchaindb_docker:/data \
           artusvranken/alpine-bigchaindb \
           start
```
