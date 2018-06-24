#!/bin/bash

if [ "$#" == "0" ]; then
    echo "This script creates a client runner connected with master server. The usage is as follows:"
    echo "crawl - crawl [number_of_pages] webpages from the internet and save to ScrapeStore, starting from [root_url]
    runclient.sh crawl [root_url] [number_of_pages] 
    runclient.sh crawl news.google.com 10
    "
    echo "getlink - retrieve all links contained in [target_url] crawled previously
    runclient.sh getlink [target_url]
    runclient.sh getlink news.google.com
    "
    echo "pagerank - run pagerank algorithm with the webpages crawled
    runclient.sh pagerank
    "
    echo "getrank - run pagerank algorithm with the webpages crawled
    runclient.sh getrank [target_url]
    runclient.sh getrank news.google.com
    "
    exit 1
fi 

if [ -z $GOPATH ]; then
    echo "FAIL: GOPATH environment variable is not set"
    exit 1
fi

if [ -n "$(go version | grep 'darwin/amd64')" ]; then    
    GOOS="darwin_amd64"
elif [ -n "$(go version | grep 'linux/amd64')" ]; then
    GOOS="linux_amd64"
else
    echo "FAIL: only 64-bit Mac OS X and Linux operating systems are supported"
    exit 1
fi

# Build the student's client node implementation.
# Exit immediately if there was a compile-time error.
go install github.com/cmu440-F15/paxosapp/runners/crunner
if [ $? -ne 0 ]; then
   echo "FAIL: code does not compile"
   exit $?
fi

# Pick random ports between [10000, 20000).
NODE_PORT0=$(((RANDOM % 10000) + 10000))
MASTER_PORT=`cat paxosnodesport.txt`
NUM=$3
if [ "$#" != "3" ]; then
  NUM=10
fi

CLIENT_NODE=$GOPATH/bin/crunner

##################################################

# Start client node.
if [ $1 == "crawl" ]; then
  ${CLIENT_NODE} -port=$NODE_PORT0 -masterPort=$MASTER_PORT -opt=$1 -url=$2 -num=$NUM  &
elif [ $1 == "getlink" ]; then
  ${CLIENT_NODE} -port=$NODE_PORT0 -masterPort=$MASTER_PORT -opt=$1 -url=$2  &
elif [ $1 == "pagerank" ]; then
  ${CLIENT_NODE} -port=$NODE_PORT0 -masterPort=$MASTER_PORT -opt=$1  &
elif [ $1 == "getrank" ]; then
  ${CLIENT_NODE} -port=$NODE_PORT0 -masterPort=$MASTER_PORT -opt=$1 -url=$2  &
else 
  echo "unknown command"  
fi

CLIENT_NODE_PID0=$!
sleep 10

# Kill client node.
# kill -9 ${CLIENT_NODE_PID0}

wait ${CLIENT_NODE_PID0} 2> /dev/null
