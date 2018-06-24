#!/bin/bash

app_name=nginx-rtmp
app_version=latest

docker build -t ncc0706/${app_name}:${app_version} .