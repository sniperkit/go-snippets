#!/bin/sh

# List of servers
servers=('127.0.0.1' 'localhost')

# Look thru the list
for i in "${servers[@]}"
do
  # SSH and run grep to check for last error
  { echo 'grep "ERROR" /var/log/mylog.log | tail -1'; } | ssh $i -t -T
done
