#!/bin/bash

go install github.com/gophergala/gomp

pushd .
cd $GOPATH/src/github.com/gophergala/gomp/tests

rm -rf *_modified* *_result tmp

for entry in *.go
do
    echo -n "Processing $entry... "
    OLD_SOURCE=$entry
    OLD_PROG=${OLD_SOURCE%'.go'}
    OLD_RESULT=$OLD_PROG'_result'
    NEW_SOURCE=${OLD_SOURCE%'.go'}'_modified.go'
    NEW_PROG=${NEW_SOURCE%'.go'}
    NEW_RESULT=$NEW_PROG'_result'
    go build $OLD_SOURCE
    gomp < $OLD_SOURCE > $NEW_SOURCE
    go build $NEW_SOURCE
    ./$OLD_PROG | sort > $OLD_RESULT
    ./$NEW_PROG | sort > $NEW_RESULT

    if diff $OLD_RESULT $NEW_RESULT >/dev/null 2>&1
    then
        echo Passed
    else
        echo Failed
        exit 1
    fi

    rm -rf $NEW_SOURCE $NEW_PROG $NEW_RESULT $OLD_PROG $OLD_RESULT tmp
done

popd 2>/dev/null 1>&2
