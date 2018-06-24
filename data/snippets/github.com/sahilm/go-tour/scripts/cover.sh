#!/usr/bin/env bash
# The coverage tool cannot handle multiple packages see (https://github.com/golang/go/issues/6909)
# This script is a workaround that runs coverage for each package and then unifies it

set -e

WORKDIR=coverage
PROFILE="$WORKDIR/cover.out"
MODE=count

test_and_cover() {
    rm -rf "$WORKDIR"
    mkdir -p "$WORKDIR"

    for pkg in "$@"; do
        f="$WORKDIR/$(echo ${pkg} | tr / -).cover"
        go test -covermode="$MODE" -coverprofile="$f" "$pkg"
    done

    echo "mode: $MODE" >"$PROFILE"
    grep -h -v "^mode:" "$WORKDIR"/*.cover >>"$PROFILE"
}

generate_html_cover_report() {
    go tool cover -html="$PROFILE" -o "$WORKDIR"/index.html
}

test_and_cover $(go list ./...)
generate_html_cover_report
