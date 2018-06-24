#!/bin/sh

docker build --build-arg SBT_VERSION="0.13.11" -t innoq/alpine-sbt -t innoq/alpine-sbt:0.13.11 .
