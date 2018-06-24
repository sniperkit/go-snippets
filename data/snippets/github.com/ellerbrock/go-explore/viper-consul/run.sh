#!/usr/bin/env bash

# start consul in dev mode
consul agent -ui -dev

# set default values from local json file encrypted in consul
crypt set -backend="consul" -endpoint="127.0.0.1:8500" -plaintext /crypt/config.json config.json

# get values
crypt get -backend="consul" -endpoint="127.0.0.1:8500" -plaintext /crypt/config.json
