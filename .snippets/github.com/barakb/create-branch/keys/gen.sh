#!/usr/bin/env bash

openssl genrsa -out server.key 2048
openssl req -config openssl.cnf -new -x509 -key server.key -out server.pem -days 3650
#openssl x509 -req -days 1024 -in server.csr -signkey server.key -out server.crt

#openssl req -config openssl.cnf -new -key server.key -out server.csr