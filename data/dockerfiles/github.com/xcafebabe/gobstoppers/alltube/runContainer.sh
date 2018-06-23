#!/bin/bash

ALLTUBE_VERSION=0.10.0
ALLTUBE_IMAGE=kafebob/rpi-alltube:0.10.0
ALLTUBE_CONTAINER=myTube

wget https://github.com/Rudloff/alltube/archive/$ALLTUBE_VERSION.zip
unzip $ALLTUBE_VERSION.zip
mv alltube-$ALLTUBE_VERSION source
sudo docker build -t $ALLTUBE_IMAGE .
sudo docker run -p 80:80 --name $ALLTUBE_CONTAINER $ALLTUBE_IMAGE

