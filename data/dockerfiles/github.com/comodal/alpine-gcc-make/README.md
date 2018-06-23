# [GCC & Make Alpine:edge Docker Image](https://hub.docker.com/r/comodal/alpine-gcc-make/) [![](https://images.microbadger.com/badges/image/comodal/alpine-gcc-make.svg)](https://microbadger.com/images/comodal/alpine-gcc-make "microbadger.com")  [![](https://images.microbadger.com/badges/commit/comodal/alpine-gcc-make.svg)](https://microbadger.com/images/comodal/alpine-gcc-make "microbadger.com")

This image is intended for compiling binaries to be used by Alpine containers.

## Docker Run

```sh
docker run --rm -it\
 -v $(pwd):/myproj\
 -w /myproj\
  comodal/alpine-stunnel:latest sh

root /myproj > make ...
```
