#!/bin/bash

# Clone source code from Gliderlabs
if [ ! -d "docker-alpine" ]; then
  git clone https://github.com/gliderlabs/docker-alpine.git
fi
cd docker-alpine/

# the builder image has to be FROM a rpi based alpine
sed -i "/FROM/ s:.*:FROM hypriot/rpi-alpine-scratch:" builder/Dockerfile

# build the builder image
sudo docker build -t alpine-builder builder

if [ -z "$ALPINE_VERSION" ]; then
  export ALPINE_VERSION="3.6"
fi

if [ -z "$ALPINE_CUSTOM_IMAGE" ]; then
  export ALPINE_CUSTOM_IMAGE="kafebob/rpi-alpine"
fi

# Generate Alpine root filesystem
sudo docker run --name alpine-build-v${ALPINE_VERSION} alpine-builder -s -t UTC -r v${ALPINE_VERSION} -m http://dl-4.alpinelinux.org/alpine > ./versions/gliderlabs-${ALPINE_VERSION}/rootfs.tar.xz

# Package filesystem in your Alpine image
sudo docker build -t ${ALPINE_CUSTOM_IMAGE}:${ALPINE_VERSION} versions/gliderlabs-${ALPINE_VERSION}

# Builded ok
echo "Successfully built. Try:  sudo docker run --rm -it --entrypoint /bin/sh --name alpine2 ${ALPINE_CUSTOM_IMAGE}"
