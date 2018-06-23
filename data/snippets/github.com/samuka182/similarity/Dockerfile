FROM golang:alpine

ENV APP_ENV ${APP_ENV}
 
ADD . /go/src/similarity

RUN go install similarity

ENTRYPOINT  /go/bin/similarity
 
EXPOSE 8080