#!/bin/sh
while [ 1 ]
do
  ps -fe|grep firefox |grep -v grep
  if [ $? -ne 0 ]
  then
  echo "start firefox ..."
  firefox &> /dev/null &
  else
  echo "firefox is running"
  fi
  sleep 30
 done
#####
