#!/bin/sh

pushd .
for d in `ls -d  $PWD/*/`; do
  cd $d
  ./build.sh
done
popd
