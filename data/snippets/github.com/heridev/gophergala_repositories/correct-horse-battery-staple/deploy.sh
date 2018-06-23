#!/bin/bash

echo "Stopping service"
ssh root@where-you-at.com "service whereyouat stop"

echo "Uploading..."
scp ./main root@where-you-at.com:/home/whereyouat/
scp -r ./assets/ root@where-you-at.com:/home/whereyouat/
scp -r ./js/ root@where-you-at.com:/usr/local/gopath/src/github.com/gophergala/correct-horse-battery-staple/
scp -r ./common/ root@where-you-at.com:/usr/local/gopath/src/github.com/gophergala/correct-horse-battery-staple/
scp -r ./mapview/ root@where-you-at.com:/usr/local/gopath/src/github.com/gophergala/correct-horse-battery-staple/
ssh root@where-you-at.com "chown -R whereyouat:whereyouat /home/whereyouat/"

echo "Starting service"
ssh root@where-you-at.com "service whereyouat start"
