#!/bin/sh
set -x
docker build --no-cache=true -t jkremser/mini-jre:8 .
