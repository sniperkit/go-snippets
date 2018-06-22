#!/bin/bash
CGO_ENABLED=0 GOOS=linux go build -o service -a -tags netgo -ldflags '-w' . && \
docker build -t gevvent/gevvent-web . && \
docker tag -f gevvent/gevvent-web gcr.io/gevvent/gevvent-web && \
gcloud preview docker push gcr.io/gevvent/gevvent-web
