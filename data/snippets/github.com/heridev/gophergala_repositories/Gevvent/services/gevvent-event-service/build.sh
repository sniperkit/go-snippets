#!/bin/bash
CGO_ENABLED=0 GOOS=linux go build -o service -a -tags netgo -ldflags '-w' . && \
docker build -t gevvent/gevvent-event-service . && \
docker tag -f gevvent/gevvent-event-service gcr.io/gevvent/gevvent-event-service && \
gcloud preview docker push gcr.io/gevvent/gevvent-event-service
