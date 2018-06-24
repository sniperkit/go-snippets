#!/bin/bash

# sudo apt-get update
# sudo apt-get -y upgrade

# sudo curl -O https://storage.googleapis.com/golang/go1.9.1.linux-amd64.tar.gz

# sudo tar -xvf go1.9.1.linux-amd64.tar.gz

# sudo mv go /usr/local

# echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile

go get github.com/pebbe/zmq4

go get github.com/chrislusf/gleam/flow

go get github.com/chrislusf/gleam/distributed/gleam

pip install bs4 ujson zmq warcio requests gevent


(go run page.go; python crawler.py) &