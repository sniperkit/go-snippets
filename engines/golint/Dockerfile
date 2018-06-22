FROM linterhub/image-go
LABEL maintainer="linterhub@repometric.com"
RUN apk update && apk --no-cache add git && \
    go get -u github.com/golang/lint/golint && \
    apk del git
