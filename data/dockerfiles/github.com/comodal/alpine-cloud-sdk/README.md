# [Google Cloud CLI Alpine Docker Image](https://hub.docker.com/r/comodal/alpine-cloud-sdk/) [![](https://images.microbadger.com/badges/image/comodal/alpine-cloud-sdk.svg)](https://microbadger.com/images/comodal/alpine-cloud-sdk "microbadger.com")  [![](https://images.microbadger.com/badges/commit/comodal/alpine-cloud-sdk.svg)](https://microbadger.com/images/comodal/alpine-cloud-sdk "microbadger.com")

## Supported Tags

* [`latest`](https://github.com/comodal/alpine-cloud-sdk/blob/master/Dockerfile) FROM openjdk:7-jre-alpine

## Docker Run gcloud

```sh
docker run -i -t --rm \
 -v $HOME/.config/gcloud:/.config/gcloud\
 -v $PWD:/data\
 -w /data\
  comodal/alpine-cloud-sdk:latest
```

## Docker Run gsutil

```sh
docker run -i -t --rm\
 -v $HOME/.config/gcloud:/.config/gcloud\
 -v $PWD:/data\
 -w /data\
 --entrypoint=bash\
  comodal/alpine-cloud-sdk:latest gsutil help
```

## Shell Alias

```sh
alias gcloud="docker run -i -t --rm\
 -v $HOME/.config/gcloud:/.config/gcloud\
 -v $PWD:/data\
 -w /data\
  comodal/alpine-cloud-sdk:latest "
```
