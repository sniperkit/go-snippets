#!/bin/sh
set -x
docker build --no-cache=true -t jkremser/mini-wildfly:10.1.0.Final .
