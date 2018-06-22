#!/bin/bash

set -u

run() {
    tests_do $BIN -c $CONFIG -s $SOURCES ${@}
    output
}

init() {
    CONFIG="$TEST_DIR/guntalina.conf"
    SOURCES="$TEST_DIR/sources"
    OUTPUT="$TEST_DIR/output"

    touch $CONFIG $SOURCES
    mkdir $TEST_DIR/conf.d/
}

config() {
    local content=$(cat)

    tests_put $CONFIG "$content"
}

confd() {
    local content=$(cat)
    local package="$1"
    local file="$2"

    mkdir $TEST_DIR/conf.d/$package
    tests_put $TEST_DIR/conf.d/$package/$file "$content"
}

sources() {
    local content=$(cat)

    tests_put $SOURCES "$content"
}

output() {
    # remove prefix '2015/10/22/ 10:17:58'
    cat `tests_stderr` | sed 's/^[^:]*:[^:]*:.. //' > "$OUTPUT"
}

assert_diff() {
    local expected="$(cat)"
    tests_diff "$expected" "$OUTPUT"
}
