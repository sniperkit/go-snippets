#!/usr/bin/env bash

source ~/.bashrc

unset GOPATH
export GOPATH=/golang/create-branch



#(cd /golang/create-branch/src/github.com/barakb/create-branch; git fetch; git rebase)

#cd /golang/create-branch

#go build -o bin/create-branch github.com/barakb/create-branch/main

#cd /golang/create-branch/src/github.com/barakb/create-branch

#nvm use v5.7.1
#npm install
#webpack
source /create-branch-conf/env.sh
/golang/create-branch/bin/create-branch -repos /create-branch-conf/repos.txt "$*"
