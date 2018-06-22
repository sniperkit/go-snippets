#!/bin/bash

# from
# http://stackoverflow.com/a/24789597/2995824

watch() {

  echo watching folder $1/ every $2 secs.

  while [[ true ]]
  do
    files=`find $1 -type f -mtime -$2s`
    if [[ $files == "" ]] ; then
      true #echo "nothing changed"
    else
      echo changed, $files
      rake
    fi
    sleep $2
  done
}

watch assets 1
