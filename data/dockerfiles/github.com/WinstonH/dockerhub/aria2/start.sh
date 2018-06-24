#!/bin/sh
nginx
touch /tmp/aria2.session
touch /tmp/aria2.log
aria2c --enable-rpc --rpc-listen-all --conf-path=/root/.aria2/aria2.conf
