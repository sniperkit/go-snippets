#!/bin/sh

docker build --rm -t wallneradam/docker-openalpr-alpine:latest .
docker push wallneradam/docker-openalpr-alpine:latest
