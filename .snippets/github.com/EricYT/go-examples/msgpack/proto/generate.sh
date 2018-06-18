#!/bin/bash

# use this script to generate message pack encoding/decoding 
# modules.

# debug
#set -x

function usage() {
  echo "usage: $(basename $0) proto-file [outputfile]"
  exit -2
}

[ $# -lt 1 ] && usage

if [ "$2" != "" ];then
  outputfile="-o $2"
else
  outputfile=""
fi

if [ -f "$1" ];then
  msgp -file $1 $outputfile
else
  echo "file not exists"
  exit -1
fi

exit 0
