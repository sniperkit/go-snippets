#!/bin/bash
CGO_ENABLED=0 GOOS=linux go build -o service -a -tags netgo -ldflags '-w' . && \
docker build -t gevvent/gevvent-user-service . && \
docker tag -f gevvent/gevvent-user-service gcr.io/gevvent/gevvent-user-service && \
gcloud preview docker push gcr.io/gevvent/gevvent-user-service
