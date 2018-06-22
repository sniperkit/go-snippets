#!/bin/bash

if [ "$1" = git ];then
  git fetch origin master
  exec git reset --hard origin/master
fi

rm -rf 5* 7*

if [ -z $WSL_HOME ];then
    cp -a ~/lnmp/dockerfile/php-fpm/. .
fi

cp -a $WSL_HOME/lnmp/dockerfile/php-fpm/. .
