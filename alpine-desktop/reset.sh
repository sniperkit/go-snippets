#!/bin/sh
while [ 1 ]
do
PID=`pgrep firefox`
kill -9 $PID
echo "It's been done."
date
sleep 3600
done
